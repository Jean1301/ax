REQ += $(shell find src -name "*.proto")
REQ += src/ax/internal/bustest/handler.go
REQ += src/ax/internal/bustest/sender.go
REQ += src/ax/internal/bustest/transport.go
REQ += src/ax/internal/bustest/pipeline.go
REQ += src/ax/internal/persistencetest/datastore.go
REQ += src/ax/internal/persistencetest/transaction.go

-include artifacts/make/go/Makefile

%.pb.go: %.proto
	protoc --go_out=. $(@D)/*.proto

MOQ := $(GOPATH)/bin/moq
$(MOQ):
	go get -u github.com/matryer/moq

src/ax/internal/bustest/handler.go: src/ax/bus/handler.go | $(MOQ)
	$(MOQ) -out "$@" -pkg "bustest" src/ax/bus MessageHandler

src/ax/internal/bustest/sender.go: src/ax/bus/sender.go | $(MOQ)
	$(MOQ) -out "$@" -pkg "bustest" src/ax/bus MessageSender

src/ax/internal/bustest/transport.go: src/ax/bus/transport.go src/ax/bus/sender.go | $(MOQ)
	$(MOQ) -out "$@" -pkg "bustest" src/ax/bus Transport

src/ax/internal/bustest/pipeline.go: src/ax/bus/pipeline.go | $(MOQ)
	$(MOQ) -out "$@" -pkg "bustest" src/ax/bus InboundPipeline OutboundPipeline

src/ax/internal/persistencetest/datastore.go: src/ax/persistence/datastore.go | $(MOQ)
	$(MOQ) -out "$@" -pkg "persistencetest" src/ax/persistence DataStore

src/ax/internal/persistencetest/transaction.go: src/ax/persistence/transaction.go | $(MOQ)
	$(MOQ) -out "$@" -pkg "persistencetest" src/ax/persistence Tx Committer

artifacts/make/%/Makefile:
	curl -sf https://jmalloc.github.io/makefiles/fetch | bash /dev/stdin $*
