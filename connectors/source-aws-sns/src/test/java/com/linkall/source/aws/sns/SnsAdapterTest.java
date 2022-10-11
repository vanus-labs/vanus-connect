package com.linkall.source.aws.sns;

import com.google.common.collect.ArrayListMultimap;
import com.google.common.collect.HashMultimap;
import io.cloudevents.CloudEvent;
import io.cloudevents.core.builder.CloudEventBuilder;
import io.netty.handler.codec.http.HttpRequest;
import io.vertx.core.Future;
import io.vertx.core.Handler;
import io.vertx.core.MultiMap;
import io.vertx.core.buffer.Buffer;
import io.vertx.core.http.*;
import io.vertx.core.http.impl.headers.HeadersAdaptor;
import io.vertx.core.http.impl.headers.HeadersMultiMap;
import io.vertx.core.json.JsonObject;
import io.vertx.core.net.NetSocket;
import junit.framework.TestCase;
import org.junit.Test;

import javax.net.ssl.SSLPeerUnverifiedException;
import javax.security.cert.X509Certificate;
import java.net.URI;
import java.time.OffsetDateTime;
import java.util.Map;

import static org.junit.Assert.*;

public class SnsAdapterTest extends TestCase {

    private SnsAdapter snsAdapter;
    private HttpRequestImpl httpRequest;
    private Buffer buffer;
    private static final CloudEventBuilder template = CloudEventBuilder.v1();
    private CloudEvent cloudEvent;

    @Test
    public void testAdapt() {
        snsAdapter = new SnsAdapter();
        httpRequest = new HttpRequestImpl();
        httpRequest.headers().add("X-Amz-Sns-Message-Id","22b80b92-fdea-4c2c-8f9d-bdfb0c7bf324");
        httpRequest.headers().add("X-Amz-Sns-Subscription-Arn","arn:aws:sns:us-west-2:123456789012:MyTopic:c9135db0-26c4-47ec-8998-413945fb5a96");
        httpRequest.headers().add("X-Amz-Sns-Topic-Arn","arn:aws:sns:us-west-2:123456789012:MyTopic");
        JsonObject jsonObject = new JsonObject();
        jsonObject.put("Type", "Notification");
        jsonObject.put("Subject", "My First Message");
        jsonObject.put("Timestamp","2012-05-02T00:54:06.655Z");
        buffer = jsonObject.toBuffer();
        template.withId("22b80b92-fdea-4c2c-8f9d-bdfb0c7bf324");
        template.withType("com.amazonaws.sns.Notification");
        template.withSubject("My First Message");
        OffsetDateTime time = OffsetDateTime.parse("2012-05-02T00:54:06.655Z");
        template.withTime(time);
        URI uri = URI.create("arn:aws:sns:us-west-2:123456789012:MyTopic:c9135db0-26c4-47ec-8998-413945fb5a96");
        template.withSource(uri)
                .withDataContentType("application/json")
                .withData(buffer.getBytes());
        cloudEvent = template.build();
        CloudEvent testCE = snsAdapter.adapt(httpRequest, buffer);
        assertEquals(cloudEvent, testCE);
    }

}
class HttpRequestImpl implements HttpServerRequest{

    private MultiMap headers = new HeadersMultiMap();

    @Override
    public HttpServerRequest exceptionHandler(Handler<Throwable> handler) {
        return null;
    }

    @Override
    public HttpServerRequest handler(Handler<Buffer> handler) {
        return null;
    }

    @Override
    public HttpServerRequest pause() {
        return null;
    }

    @Override
    public HttpServerRequest resume() {
        return null;
    }

    @Override
    public HttpServerRequest fetch(long amount) {
        return null;
    }

    @Override
    public HttpServerRequest endHandler(Handler<Void> endHandler) {
        return null;
    }

    @Override
    public HttpVersion version() {
        return null;
    }

    @Override
    public HttpMethod method() {
        return null;
    }

    @Override
    public String scheme() {
        return null;
    }

    @Override
    public String uri() {
        return null;
    }

    @Override
    public String path() {
        return null;
    }

    @Override
    public String query() {
        return null;
    }

    @Override
    public String host() {
        return null;
    }

    @Override
    public long bytesRead() {
        return 0;
    }

    @Override
    public HttpServerResponse response() {
        return null;
    }

    @Override
    public MultiMap headers() {

        return headers;
    }

    @Override
    public MultiMap params() {
        return null;
    }

    @Override
    public X509Certificate[] peerCertificateChain() throws SSLPeerUnverifiedException {
        return new X509Certificate[0];
    }

    @Override
    public String absoluteURI() {
        return null;
    }

    @Override
    public Future<Buffer> body() {
        return null;
    }

    @Override
    public Future<Void> end() {
        return null;
    }

    @Override
    public Future<NetSocket> toNetSocket() {
        return null;
    }

    @Override
    public HttpServerRequest setExpectMultipart(boolean expect) {
        return null;
    }

    @Override
    public boolean isExpectMultipart() {
        return false;
    }

    @Override
    public HttpServerRequest uploadHandler(Handler<HttpServerFileUpload> uploadHandler) {
        return null;
    }

    @Override
    public MultiMap formAttributes() {
        return null;
    }

    @Override
    public String getFormAttribute(String attributeName) {
        return null;
    }

    @Override
    public Future<ServerWebSocket> toWebSocket() {
        return null;
    }

    @Override
    public boolean isEnded() {
        return false;
    }

    @Override
    public HttpServerRequest customFrameHandler(Handler<HttpFrame> handler) {
        return null;
    }

    @Override
    public HttpConnection connection() {
        return null;
    }

    @Override
    public HttpServerRequest streamPriorityHandler(Handler<StreamPriority> handler) {
        return null;
    }

    @Override
    public Map<String, Cookie> cookieMap() {
        return null;
    }
}