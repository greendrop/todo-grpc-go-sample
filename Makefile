BINDIR:=$(abspath bin)
PROTO_FILES:=$(shell find ./proto -path './proto/validate' -prune -o -path './proto/google' -prune -o -type f -name '*.proto' -print)

install-go-tools:
	@mkdir -p ${BINDIR}
	@cat tools.go | awk -F'"' '/_/ {print $$2}' | GOBIN=$(abspath ${BINDIR}) xargs -tI {} go install {}

lint-proto:
	bin/protolint lint proto/.

lint-fix-proto:
	bin/protolint lint -fix proto/.

build-proto:
	@for proto_file in ${PROTO_FILES}; do \
		echo $$proto_file; \
		protoc \
			-I ./proto \
			--plugin=protoc-gen-go=$(BINDIR)/protoc-gen-go \
			--plugin=protoc-gen-go-grpc=$(BINDIR)/protoc-gen-go-grpc \
			--plugin=protoc-gen-grpc-gateway=$(BINDIR)/protoc-gen-grpc-gateway \
			--plugin=protoc-gen-validate=$(BINDIR)/protoc-gen-validate \
			--go_out=./proto \
			--go_opt=paths=source_relative \
			--go-grpc_out=./proto \
			--go-grpc_opt=paths=source_relative \
			--grpc-gateway_out=./proto \
			--grpc-gateway_opt=paths=source_relative \
			--grpc-gateway_opt=generate_unbound_methods=true \
			--validate_out="lang=go,paths=source_relative:./proto" \
			$$proto_file; \
	done
