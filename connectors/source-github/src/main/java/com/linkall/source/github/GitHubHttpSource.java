package com.linkall.source.github;

import com.linkall.cdk.config.Config;
import com.linkall.cdk.connector.Element;
import com.linkall.cdk.connector.Source;
import com.linkall.cdk.connector.Tuple;
import io.cloudevents.CloudEvent;
import io.vertx.core.Vertx;
import org.apache.commons.codec.digest.HmacAlgorithms;
import org.apache.commons.codec.digest.HmacUtils;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.concurrent.BlockingQueue;
import java.util.concurrent.LinkedBlockingQueue;
import java.util.concurrent.atomic.AtomicInteger;

public class GitHubHttpSource implements Source {

    private static final Logger LOGGER = LoggerFactory.getLogger(GitHubHttpSource.class);
    private static final AtomicInteger eventNum = new AtomicInteger(0);
    private static final Vertx vertx = Vertx.vertx();
    private BlockingQueue<Tuple> queue;
    private GitHubConfig config;

    public GitHubHttpSource() {
        queue = new LinkedBlockingQueue<>(100);
    }

    public void start() {
        String secret;
        if (config.getSecretConfig()!=null) {
            secret = config.getSecretConfig().getGithubWebHookSecret();
        } else {
            secret = null;
        }
        GitHubHttpAdapter adapter = new GitHubHttpAdapter();
        int port = config.getPort();
        if (port <= 0) {
            port = 8080;
        }
        vertx.createHttpServer()
                .exceptionHandler(System.err::println)
                .requestHandler(request -> {
                    //process github ping event
                    if (request.getHeader("X-Github-Event").equals("ping")) {
                        request.response().setStatusCode(200);
                        request.response().end("receive ping event success");
                        return;
                    }
                    String signature = request.getHeader("X-Hub-Signature-256");
                    request.bodyHandler(body -> {
                        // if header contains signature we verify it
                        eventNum.addAndGet(1);
                        LOGGER.info("receive a request " + eventNum.get());
                        String bodyStr = body.toJsonObject().toString();
                        if (null!=signature) {
                            // if signature existed in headers but user didn't provide secret, program exits.
                            if (null==secret) {
                                request.response().setStatusCode(505);
                                request.response().end("signature verified failed");
                                System.exit(1);
                            }
                            // verify signature failed
                            if (!verifySignature(signature, bodyStr, secret)) {
                                LOGGER.info("signature verified failed");
                                request.response().setStatusCode(505);
                                request.response().end("signature verified failed");
                                return;
                            }
                        }
                        //transform http payload to CloudEvent
                        CloudEvent ce = adapter.adapt(request, body);
                        Tuple tuple = new Tuple(new Element(ce, body), () -> {
                            LOGGER.info("send task success");
                            request.response().setStatusCode(200);
                            request.response().end("Receive success, deliver CloudEvents to"
                                    + config.getTarget() + "success");
                        }, (success, failed, msg) -> {
                            LOGGER.warn("send task failed,{}", msg);
                            request.response().setStatusCode(504);
                            request.response().end("Receive success, deliver CloudEvents to"
                                    + config.getTarget() + "failed");
                        });
                        try {
                            queue.put(tuple);
                        } catch (InterruptedException e) {
                            LOGGER.warn("put event interrupted");
                        }
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

    public boolean verifySignature(String signature, String bodyStr, String secret) {
        String hex = "sha256=" + new HmacUtils(HmacAlgorithms.HMAC_SHA_256, secret).hmacHex(bodyStr);

        if (!hex.equals(signature)) {
            return false;
        }
        return true;
    }

    @Override
    public BlockingQueue<Tuple> queue() {
        return queue;
    }

    @Override
    public Class<? extends Config> configClass() {
        return GitHubConfig.class;
    }

    @Override
    public void initialize(Config config) throws Exception {
        this.config = (GitHubConfig) config;
        start();
    }

    @Override
    public String name() {
        return "GitHubSource";
    }

    @Override
    public void destroy() throws Exception {
        vertx.close();
    }
}
