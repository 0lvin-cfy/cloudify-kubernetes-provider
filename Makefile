all: bin/cfy-go bin/cfy-kubernetes

PACKAGEPATH := github.com/0lvin-cfy/cloudify-rest-go-client

reformat:
	rm -rfv pkg/*
	rm -rfv bin/*
	gofmt -w src/${PACKAGEPATH}/cloudifyrest/*.go
	gofmt -w src/${PACKAGEPATH}/cloudifyutils/*.go
	gofmt -w src/${PACKAGEPATH}/cloudify/*.go
	gofmt -w src/${PACKAGEPATH}/cfy-go/*.go
	gofmt -w src/cloudifyprovider/*.go
	gofmt -w src/*.go

define colorecho
	@tput setaf 2
	@echo -n $1
	@tput setaf 3
	@echo $2
	@tput sgr0
endef

# cloudify rest
CLOUDIFYREST := \
	src/${PACKAGEPATH}/cloudifyrest/rest.go \
	src/${PACKAGEPATH}/cloudifyrest/types.go

pkg/linux_amd64/${PACKAGEPATH}/cloudifyrest.a: ${CLOUDIFYREST}
	$(call colorecho,"Build: ", $@)
	go build -v -i -o pkg/linux_amd64/${PACKAGEPATH}/cloudifyrest.a ${CLOUDIFYREST}

# cloudify utils
CLOUDIFYUTILS := \
	src/${PACKAGEPATH}/cloudifyutils/utils.go

pkg/linux_amd64/${PACKAGEPATH}/cloudifyutils.a: ${CLOUDIFYUTILS}
	$(call colorecho,"Build: ", $@)
	go build -v -i -o pkg/linux_amd64/${PACKAGEPATH}/cloudifyutils.a ${CLOUDIFYUTILS}

# cloudify
CLOUDIFYCOMMON := \
	src/${PACKAGEPATH}/cloudify/client.go \
	src/${PACKAGEPATH}/cloudify/events.go \
	src/${PACKAGEPATH}/cloudify/blueprints.go \
	src/${PACKAGEPATH}/cloudify/status.go \
	src/${PACKAGEPATH}/cloudify/executions.go \
	src/${PACKAGEPATH}/cloudify/deployments.go

pkg/linux_amd64/cloudify.a: ${CLOUDIFYCOMMON} pkg/linux_amd64/${PACKAGEPATH}/cloudifyrest.a
	$(call colorecho,"Build: ",$@)
	go build -v -i -o pkg/linux_amd64/cloudify.a ${CLOUDIFYCOMMON}

bin/cfy-go: src/${PACKAGEPATH}/cfy-go/cfy-go.go pkg/linux_amd64/${PACKAGEPATH}/cloudifyutils.a pkg/linux_amd64/cloudify.a
	$(call colorecho,"Install: ", $@)
	go install -v -ldflags "-X main.versionString=`cd src/${PACKAGEPATH} && git rev-parse --short HEAD`" src/${PACKAGEPATH}/cfy-go/cfy-go.go

bin/cfy-kubernetes: pkg/linux_amd64/cloudifyprovider.a pkg/linux_amd64/cloudify.a src/cfy-kubernetes.go
	$(call colorecho,"Install: ", $@)
	go install -v src/cfy-kubernetes.go

# cloudify provider
CLOUDIFYPROVIDER := src/cloudifyprovider/init.go src/cloudifyprovider/instances.go

pkg/linux_amd64/cloudifyprovider.a: pkg/linux_amd64/cloudify.a ${CLOUDIFYPROVIDER}
	go build -v -i -o pkg/linux_amd64/cloudifyprovider.a ${CLOUDIFYPROVIDER}

test:
	go test ./src/${PACKAGEPATH}/...
