FROM --platform=$BUILDPLATFORM python:3.11 as builder

WORKDIR /build/vanus-connect

COPY . .
RUN python3 setup.py bdist_wheel

FROM --platform=$TARGETPLATFORM public.ecr.aws/vanus/connector/service-openai:latest

COPY --from=builder /build/vanus-connect/dist/vanus_connect_milvussink-*.whl .
RUN python3 -m pip install --no-cache-dir \
        vanus_connect_milvussink-*.whl \
        && rm vanus_connect_milvussink-*.whl

EXPOSE 8080

CMD ["python3", "-m", "vanus.connect.sink.milvus", "--config=/vanus-connect/config/config.yml"]
