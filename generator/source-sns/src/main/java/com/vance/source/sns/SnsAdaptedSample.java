package com.vance.source.sns;

import io.vertx.core.buffer.Buffer;
import io.vertx.core.http.HttpServerRequest;

public class SnsAdaptedSample{

    //wrap attributes used by Adapter to transform data into CloudEvents

    //example:if you need http request to get data for CloudEvents
    public HttpServerRequest httpServerRequest;
    public Buffer buffer;

    public SnsAdaptedSample(){}

    public SnsAdaptedSample(HttpServerRequest httpServerRequest, Buffer buffer){
        this.httpServerRequest = httpServerRequest;
        this.buffer = buffer;
    }

}