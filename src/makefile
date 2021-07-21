.PHONY: apps router

apps:
	(cd cmd/daemon; go build)
	(cd unit/daemon; go build)
	(cd unit/halsimulation; go build)

router:
	(cd router/; go build)
