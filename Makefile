all:
	make build-http-handler
	make build-mq-handler

build-http-handler:
	zip -r http-handler.zip handler.go queue.go types.go utils.go go.mod go.sum

build-mq-trigger:
	zip -r mq-trigger.zip trigger.go types.go utils.go go.mod go.sum