.PHONY: image deploy

VERSION = v0.0.2

image:
	nerdctl build --namespace k8s.io -f docker/Dockerfile --tag spolab/petclinic-owner:$(VERSION) .
