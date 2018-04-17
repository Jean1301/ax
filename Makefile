REQ += src/ax/internal/messagetest/messages.pb.go

-include artifacts/make/go/Makefile

src/ax/internal/messagetest/messages.pb.go: $(wildcard src/ax/internal/messagetest/*.proto)
	protoc --go_out=. $^

artifacts/make/%/Makefile:
	curl -sf https://jmalloc.github.io/makefiles/fetch | bash /dev/stdin $*
