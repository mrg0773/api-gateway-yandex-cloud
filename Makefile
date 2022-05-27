all:
	make build-http-handler
	make build-mq-handler

http-handler:
	zip -r http-handler.zip handler.go queue.go types.go utils.go go.mod go.sum

mq-trigger:
	zip -r mq-trigger.zip trigger.go queue.go types.go utils.go go.mod go.sum