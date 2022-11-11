package com.linkall.source.github;

import com.linkall.vance.common.config.ConfigUtil;
import com.linkall.vance.common.config.SecretUtil;
import com.linkall.vance.core.Adapter;
import com.linkall.vance.core.Adapter2;
import com.linkall.vance.core.Source;
import com.linkall.vance.core.http.HttpClient;
import com.linkall.vance.core.http.HttpResponseInfo;
import io.cloudevents.CloudEvent;
import io.cloudevents.core.message.MessageReader;
import io.cloudevents.http.vertx.VertxMessageFactory;
import io.cloudevents.jackson.JsonFormat;
import io.vertx.core.Future;
import io.vertx.core.Vertx;
import io.vertx.core.buffer.Buffer;
import io.vertx.core.http.HttpServerRequest;
import io.vertx.ext.web.client.HttpResponse;
import io.vertx.ext.web.client.WebClient;
import org.apache.commons.codec.digest.HmacAlgorithms;
import org.apache.commons.codec.digest.HmacUtils;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.concurrent.atomic.AtomicInteger;

public class GitHubHttpSource implements Source {

    private static final Logger LOGGER = LoggerFactory.getLogger(GitHubHttpSource.class);
    private static final AtomicInteger eventNum = new AtomicInteger(0);
    private static final Vertx vertx = Vertx.vertx();
    private WebClient webClient = WebClient.create(vertx);
    @Override
    public Adapter getAdapter() {
        return new GitHubHttpAdapter();
    }

    @Override
    public void start() {
        String secret = SecretUtil.getString("githubWebHookSecret");
        Adapter2<HttpServerRequest, Buffer> adapter = (Adapter2<HttpServerRequest, Buffer>) getAdapter();
        String vanceSink = ConfigUtil.getVanceSink();
        int port = Integer.parseInt(ConfigUtil.getPort());
        vertx.createHttpServer()
                .exceptionHandler(System.err::println)
                .requestHandler(request -> {
                    //process github ping event
                    if(request.getHeader("X-Github-Event").equals("ping")){
                        request.response().setStatusCode(200);
                        request.response().end("receive ping event success");
                        return;
                    }
                    String signature = request.getHeader("X-Hub-Signature-256");
                    request.bodyHandler(body -> {
                        // if header contains signature we verify it
                        eventNum.addAndGet(1);
                        LOGGER.info("receive a request "+eventNum.get());
                        String bodyStr = body.toJsonObject().toString();
                        if(null != signature){
                            // if signature existed in headers but user didn't provide secret, program exits.
                            if(null == secret){
                                request.response().setStatusCode(505);
                                request.response().end("signature verified failed");
                                System.exit(1);
                            }
                            // verify signature failed
                            if (!verifySignature(signature, bodyStr,secret)) {
                                LOGGER.info("signature verified failed");
                                request.response().setStatusCode(505);
                                request.response().end("signature verified failed");
                                return;
                            }
                        }
                        //transform http payload to CloudEvent
                        CloudEvent ce = adapter.adapt(request, body);
                        Future<HttpResponse<Buffer>> responseFuture = VertxMessageFactory.createWriter(webClient.postAbs(vanceSink))
                                .writeStructured(ce, JsonFormat.CONTENT_TYPE);
                        responseFuture.onSuccess(resp-> {
                            LOGGER.info("send task success");
                            request.response().setStatusCode(200);
                            request.response().end("Receive success, deliver CloudEvents to"
                                    + ConfigUtil.getVanceSink() + "success");
                        }).onFailure(t->{
                            LOGGER.info("send task failed");
                            request.response().setStatusCode(504);
                            request.response().end("Receive success, deliver CloudEvents to"
                                    + ConfigUtil.getVanceSink() + "failed");
                        });
                    });
                })
                .listen(port, server -> {
                    if (server.succeeded()) {
                        LOGGER.info("Server listening on port: " + server.result().actualPort());

                    } else {
                        LOGGER.error(server.cause().getMessage());
                    }
                });
    }
    public boolean verifySignature(String signature, String bodyStr,String secret){
        String hex = "sha256=" + new HmacUtils(HmacAlgorithms.HMAC_SHA_256,secret).hmacHex(bodyStr);

        if(!hex.equals(signature)){
            return false;
        }
        return true;
    }
}
