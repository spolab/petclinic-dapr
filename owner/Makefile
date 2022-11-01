.PHONY: image deploy

image:
	nerdctl build --namespace k8s.io -f docker/Dockerfile --tag spolab/petclinic-owner:v0.0.1 .
