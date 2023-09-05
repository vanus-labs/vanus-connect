FROM --platform=$BUILDPLATFORM python:3.11 as builder
WORKDIR /build/vanus-connect

COPY . .
RUN python3 setup.py bdist_wheel

FROM --platform=$TARGETPLATFORM python:3.11
WORKDIR /vanus-connect

COPY --from=builder /build/vanus-connect/dist/vanus_cdk-*whl .
RUN python3 -m pip install vanus_cdk-*whl && rm vanus_cdk-*whl
