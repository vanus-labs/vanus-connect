WORKDIR=$(shell pwd)
PROTO_ROOT=${WORKDIR}/../../schemas

generate-pb:
	mkdir -p proto
	protoc -I=${PROTO_ROOT} \
		   -I=${PROTO_ROOT}/thirds \
           --go_out=paths=:proto \
           ${PROTO_ROOT}/database/mongodb.proto