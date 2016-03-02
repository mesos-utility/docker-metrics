default: help
#docker_image: dockerized/.dockerized.created ## Create the calico/mesos-calico image

# TODO: maybe change this so docker runs and handles the caching itself,
# instead of relying on the .created file.
# dockerized/.dockerized.created:
#	./control build
#	docker build -t mesos-utility/docker-metrics .
#	touch dockerized/.dockerized_image.created

## Make bin for docker-metrics
bin:
	./control build

## Clean everything (including stray volumes)
clean:
#	find . -name '*.created' -exec rm -f {} +
	-rm -rf var
	-rm -f docker-metrics
#	-docker rmi mesos-utility/docker-metrics

help: # Some kind of magic from https://gist.github.com/rcmachado/af3db315e31383502660
	$(info Available targets)
	@awk '/^[a-zA-Z\-\_0-9]+:/ {                                   \
		nb = sub( /^## /, "", helpMsg );                             \
		if(nb == 0) {                                                \
			helpMsg = $$0;                                             \
			nb = sub( /^[^:]*:.* ## /, "", helpMsg );                  \
		}                                                            \
		if (nb)                                                      \
			printf "\033[1;31m%-" width "s\033[0m %s\n", $$1, helpMsg; \
	}                                                              \
	{ helpMsg = $$0 }'                                             \
	width=$$(grep -o '^[a-zA-Z_0-9]\+:' $(MAKEFILE_LIST) | wc -L)  \
	$(MAKEFILE_LIST)

