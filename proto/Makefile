PROTO_ROOT=$(shell pwd)

generate-base:
	mkdir -p pkg/base
	protoc -I=${PROTO_ROOT} \
		   -I=${PROTO_ROOT}/thirds \
           --go_out=paths=source_relative:${PROTO_ROOT}/pkg \
           ${PROTO_ROOT}/base/*.proto

generate-database:
	mkdir -p pkg/database
	protoc -I=${PROTO_ROOT} \
		   -I=${PROTO_ROOT}/thirds \
           --go_out=paths=source_relative:${PROTO_ROOT}/pkg \
           ${PROTO_ROOT}/database/database.proto
