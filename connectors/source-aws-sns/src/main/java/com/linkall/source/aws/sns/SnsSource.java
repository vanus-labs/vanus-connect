package com.linkall.source.aws.sns;

import com.linkall.cdk.config.Config;
import com.linkall.cdk.connector.Element;
import com.linkall.cdk.connector.Source;
import com.linkall.cdk.connector.Tuple;
import com.linkall.source.aws.utils.AwsHelper;
import com.linkall.source.aws.utils.SNSUtil;
import io.cloudevents.CloudEvent;
import io.vertx.core.Vertx;
import io.vertx.core.json.JsonObject;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import software.amazon.awssdk.regions.Region;
import software.amazon.awssdk.services.sns.SnsClient;
import software.amazon.awssdk.services.sns.model.SnsException;

import java.io.ByteArrayInputStream;
import java.util.concurrent.BlockingQueue;
import java.util.concurrent.LinkedBlockingQueue;
import java.util.concurrent.atomic.AtomicInteger;

public class SnsSource implements Source {

    private static final Logger LOGGER = LoggerFactory.getLogger(SnsSource.class);
    private static final AtomicInteger eventNum = new AtomicInteger(0);
    private static final Vertx vertx = Vertx.vertx();


    private BlockingQueue<Tuple> queue;
    private SnsConfig config;

    private SnsClient snsClient;
    private String subscribeArn;

    public SnsSource(){
        queue = new LinkedBlockingQueue<>(100);
    }

    public void start() {
        AwsHelper.checkCredentials(config.getSecretConfig().getAccessKeyID(), config.getSecretConfig().getSecretAccessKey());
        SnsAdapter adapter = new SnsAdapter();

        String snsTopicArn = config.getSnsArn();
        String region = SNSUtil.getRegion(snsTopicArn);
        String host = config.getEndpoint();
        String protocol = config.getProtocol();
        int port = config.getPort();
        if (port <= 0) {
            port = 8080;
        }

        snsClient = SnsClient.builder().region(Region.of(region)).build();
        try {
            subscribeArn = SNSUtil.subHTTPS(snsClient, snsTopicArn, host, protocol);
        } catch (SnsException e) {
            LOGGER.error(e.awsErrorDetails().errorMessage());
            snsClient.close();
            System.exit(1);
        }

        vertx.createHttpServer()
                .exceptionHandler(failed->{
                    LOGGER.error("unknown error",failed.getCause());
                })
                .requestHandler(request -> {
                    String messageType = request.getHeader("x-amz-sns-message-type");
                    request.bodyHandler(body -> {
                        JsonObject jsonObject = body.toJsonObject();
                        String token = jsonObject.getString("Token");
                        if (!SNSUtil.verifySignatrue(new ByteArrayInputStream(body.getBytes()), region)) {
                            request.response().setStatusCode(505);
                            request.response().end("signature verified failed");
                        } else {
                            //confirm sub or unSub
                            LOGGER.info("verify signature successful");
                            if (messageType.equals("SubscriptionConfirmation") || messageType.equals("UnsubscribeConfirmation")) {
                                try {
                                    SNSUtil.confirmSubHTTPS(snsClient, token, snsTopicArn);
                                } catch (SnsException e) {
                                    LOGGER.error(e.awsErrorDetails().errorMessage());
                                    snsClient.close();
                                    System.exit(1);
                                }
                            }

                            CloudEvent ce = adapter.adapt(request, body);
                            Tuple tuple = new Tuple(new Element(ce, jsonObject), () -> {
                                LOGGER.info("send event success,{}", ce.getId());
                                eventNum.getAndAdd(1);
                                LOGGER.info("send " + eventNum + " CloudEvents in total");
                                request.response().setStatusCode(200);
                                request.response().end("Receive success, deliver CloudEvents to");
                            }, (success, failed, msg) -> {
                                LOGGER.info("send event failed,{},{}", ce.getId(), msg);
                                request.response().setStatusCode(500);
                                request.response().end("Receive success, deliver CloudEvents failed " + msg);
                            });
                            try {
                                queue.put(tuple);
                            } catch (InterruptedException e) {
                                LOGGER.warn("put event interrupted");
                            }
                        }

                    });
                }).listen(port, (server) -> {
                    if (server.succeeded()) {
                        LOGGER.info("HttpServer is listening on port: " + (server.result()).actualPort());
                    } else {
                        LOGGER.error(server.cause().getMessage());
                    }
                });

        String finalSubscribeArn = subscribeArn;
        Runtime.getRuntime().addShutdownHook(new Thread(() -> {
            try {
                SNSUtil.unSubHTTPS(snsClient, finalSubscribeArn);
            } catch (SnsException e) {
                LOGGER.error(e.awsErrorDetails().errorMessage());
            }
            snsClient.close();

            LOGGER.info("shut down!");
        }));
    }


    @Override
    public BlockingQueue<Tuple> queue() {
        return queue;
    }

    @Override
    public Class<? extends Config> configClass() {
        return SnsConfig.class;
    }

    @Override
    public void initialize(Config config) throws Exception {
        this.config = (SnsConfig) config;
        start();
    }

    @Override
    public String name() {
        return "AmazonSNSSource";
    }

    @Override
    public void destroy() {
        if (snsClient==null) {
            return;
        }
        try {
            SNSUtil.unSubHTTPS(snsClient, subscribeArn);
        } catch (SnsException e) {
            LOGGER.error(e.awsErrorDetails().errorMessage());
        }
        snsClient.close();
    }
}
