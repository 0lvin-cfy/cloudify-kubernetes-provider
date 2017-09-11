/*
Copyright (c) 2017 GigaSpaces Technologies Ltd. All rights reserved

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cloudifyprovider

import (
	"encoding/json"
	"fmt"
	cloudify "github.com/0lvin-cfy/cloudify-rest-go-client/cloudify"
	"github.com/golang/glog"
	"io"
	"k8s.io/kubernetes/pkg/cloudprovider"
	"k8s.io/kubernetes/pkg/controller"
	"os"
)

const (
	providerName = "cloudify"
)

// CloudProvider implents Instances, Zones, and LoadBalancer
type CloudProvider struct {
	deployment string
	client     *cloudify.CloudifyClient
	instances  *CloudifyIntances
	balancers  *CloudifyBalancer
	zones      *CloudifyZones
}

// Initialize passes a Kubernetes clientBuilder interface to the cloud provider
func (r *CloudProvider) Initialize(clientBuilder controller.ControllerClientBuilder) {
	glog.Warning("Initialize")
}

// ProviderName returns the cloud provider ID.
func (r *CloudProvider) ProviderName() string {
	return providerName
}

// LoadBalancer returns a balancer interface. Also returns true if the interface is supported, false otherwise.
func (r *CloudProvider) LoadBalancer() (cloudprovider.LoadBalancer, bool) {
	glog.Warning("LoadBalancer")
	if r.client != nil {
		if r.balancers != nil {
			return r.balancers, true
		} else {
			r.balancers = NewCloudifyBalancer(r.client)
			return r.balancers, true
		}
	}
	return nil, false
}

// Zones returns a zones interface. Also returns true if the interface is supported, false otherwise.
func (r *CloudProvider) Zones() (cloudprovider.Zones, bool) {
	glog.Warning("Zones")
	if r.client != nil {
		if r.zones != nil {
			return r.zones, true
		} else {
			r.zones = NewCloudifyZones(r.client)
			return r.zones, true
		}
	}
	return nil, false
}

// Instances returns an instances interface. Also returns true if the interface is supported, false otherwise.
func (r *CloudProvider) Instances() (cloudprovider.Instances, bool) {
	glog.Warning("Instances")
	if r.client != nil {
		if r.instances != nil {
			return r.instances, true
		} else {
			r.instances = NewCloudifyIntances(r.client, r.deployment)
			return r.instances, true
		}
	}
	return nil, false
}

// Clusters returns a clusters interface.  Also returns true if the interface is supported, false otherwise.
func (r *CloudProvider) Clusters() (cloudprovider.Clusters, bool) {
	glog.Warning("Clusters")
	return nil, false
}

// Routes returns a routes interface along with whether the interface is supported.
func (r *CloudProvider) Routes() (cloudprovider.Routes, bool) {
	glog.Warning("Routers")
	return nil, false
}

// HasClusterID returns true if a ClusterID is required and set
func (r *CloudProvider) HasClusterID() bool {
	return false
}

// ScrubDNS provides an opportunity for cloud-provider-specific code to process DNS settings for pods.
func (r *CloudProvider) ScrubDNS(nameservers, searches []string) (nsOut, srchOut []string) {
	return nameservers, searches
}

type CloudifyProviderConfig struct {
	Host       string `json:"host,omitempty"`
	User       string `json:"user,omitempty"`
	Password   string `json:"password,omitempty"`
	Tenant     string `json:"tenant,omitempty"`
	Deployment string `json:"deployment,omitempty"`
}

func newCloudifyCloud(config io.Reader) (cloudprovider.Interface, error) {
	glog.Warning("New Cloudify client")

	var cloudConfig CloudifyProviderConfig
	cloudConfig.Host = os.Getenv("CFY_HOST")
	cloudConfig.User = os.Getenv("CFY_USER")
	cloudConfig.Password = os.Getenv("CFY_PASSWORD")
	cloudConfig.Tenant = os.Getenv("CFY_TENANT")
	if config != nil {
		err := json.NewDecoder(config).Decode(&cloudConfig)
		if err != nil {
			return nil, err
		}
	}

	if len(cloudConfig.Host) == 0 {
		return nil, fmt.Errorf("You have empty host")
	}

	if len(cloudConfig.User) == 0 {
		return nil, fmt.Errorf("You have empty user")
	}

	if len(cloudConfig.Password) == 0 {
		return nil, fmt.Errorf("You have empty password")
	}

	if len(cloudConfig.Tenant) == 0 {
		return nil, fmt.Errorf("You have empty tenant")
	}

	if len(cloudConfig.Deployment) == 0 {
		return nil, fmt.Errorf("You have empty deployment")
	}

	glog.Warning("Config %+v", cloudConfig)
	return &CloudProvider{
		deployment: cloudConfig.Deployment,
		client: cloudify.NewClient(
			cloudConfig.Host, cloudConfig.User,
			cloudConfig.Password, cloudConfig.Tenant),
	}, nil
}

func init() {
	glog.Warning("Cloudify init")
	cloudprovider.RegisterCloudProvider(providerName, func(config io.Reader) (cloudprovider.Interface, error) {
		return newCloudifyCloud(config)
	})
}
