FROM --platform=$BUILDPLATFORM python:3.11 as builder

WORKDIR /build/vanus-connect

COPY . .
RUN python3 setup.py bdist_wheel

FROM --platform=$TARGETPLATFORM public.ecr.aws/vanus/connector/python:latest

COPY --from=builder /build/vanus-connect/dist/vanus_connect_openaiservice-*.whl .
RUN python3 -m pip install --no-cache-dir \
        vanus_connect_openaiservice-*.whl \
        && rm vanus_connect_openaiservice-*.whl
