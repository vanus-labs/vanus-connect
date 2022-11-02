WORKDIR=$(shell pwd)
PROTO_ROOT=${WORKDIR}/../../proto

generate-pb:
	mkdir -p proto
	protoc -I=${PROTO_ROOT} \
		   -I=${PROTO_ROOT}/thirds \
           --go_out=paths=source_relative:proto \
           ${PROTO_ROOT}/database/mongodb.proto