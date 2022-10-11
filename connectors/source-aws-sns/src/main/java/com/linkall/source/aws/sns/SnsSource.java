package com.linkall.source.aws.sns;

import com.amazonaws.services.sns.message.SnsMessageManager;
import com.linkall.source.aws.utils.AwsHelper;
import com.linkall.source.aws.utils.SNSUtil;
import com.linkall.vance.common.config.ConfigUtil;
import com.linkall.vance.core.Adapter;
import com.linkall.vance.core.Source;
import com.linkall.vance.core.http.HttpClient;

import com.linkall.vance.core.http.HttpResponseInfo;
import io.cloudevents.CloudEvent;
import io.cloudevents.http.vertx.VertxMessageFactory;
import io.vertx.core.Future;
import io.vertx.core.Vertx;
import io.vertx.core.buffer.Buffer;
import io.vertx.core.http.HttpServer;
import io.vertx.core.json.JsonObject;
import io.vertx.ext.web.Router;
import io.vertx.ext.web.RoutingContext;
import io.vertx.ext.web.client.HttpResponse;
import io.vertx.ext.web.client.WebClient;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import software.amazon.awssdk.regions.Region;
import software.amazon.awssdk.services.sns.SnsClient;
import software.amazon.awssdk.services.sns.model.SnsException;

import java.io.ByteArrayInputStream;
import java.util.concurrent.atomic.AtomicInteger;

public class SnsSource implements Source {

    private static final Logger LOGGER = LoggerFactory.getLogger(SnsSource.class);
    private static final AtomicInteger eventNum = new AtomicInteger(0);
    private static final Vertx vertx = Vertx.vertx();
    private static Router router;
    private HttpServer httpServer;
    private HttpResponseInfo handlerRI;
    private static final WebClient webClient = WebClient.create(vertx);
    private String snsTopicArn;
    private String region;
    private String endPoint;
    private String protocol;
    private SnsAdapter adapter;
    private SnsClient snsClient;

    public void init(){
        AwsHelper.checkCredentials();
        this.adapter = (SnsAdapter) getAdapter();

        this.snsTopicArn = ConfigUtil.getString("topic_arn");
        this.region = SNSUtil.getRegion(snsTopicArn);
        this.endPoint = ConfigUtil.getString("endpoint");
        this.protocol = ConfigUtil.getString("protocol");

        this.snsClient = SnsClient.builder().region(Region.of(this.region)).build();

        this.httpServer = vertx.createHttpServer();
        this.router = Router.router(vertx);
        this.handlerRI = new HttpResponseInfo(200, "Receive success, deliver CloudEvents to"
                + ConfigUtil.getVanceSink() + "success", 500, "Receive success, deliver CloudEvents to"
                + ConfigUtil.getVanceSink() + "failure");
    }

    @Override
    public void start(){
        this.init();

        String subscribeArn = subscribe(this.snsClient, this.snsTopicArn, this.endPoint, this.protocol);

        this.router.route("/").handler(request-> {
            request.request().bodyHandler(body->{
                JsonObject jsonObject = body.toJsonObject();
                SnsMessageManager manager = new SnsMessageManager(region);

                boolean verifyResult = verifySignature(manager, snsClient, jsonObject, snsTopicArn);

                if(verifyResult){
                    CloudEvent ce = adapter.adapt(request.request(), body);

                    Future<HttpResponse<Buffer>> responseFuture;
                    String vanceSink = ConfigUtil.getVanceSink();
                    responseFuture = VertxMessageFactory.createWriter(webClient.postAbs(vanceSink))
                            .writeStructured(ce, "application/cloudevents+json");

                    responseFuture.onSuccess(resp->{
                       LOGGER.info("send CloudEvent to " + vanceSink + " success");
                       eventNum.getAndAdd(1);
                       LOGGER.info("send " + eventNum + " CloudEvents in total");
                       HttpResponseInfo info = this.handlerRI;
                       request.response().setStatusCode(info.getSuccessCode());
                       request.response().end(info.getSuccessChunk());
                    });
                    responseFuture.onFailure(resp->{
                        LOGGER.error("send CloudEvent to " + vanceSink + " failure");
                        LOGGER.info("send " + eventNum + " CloudEvents in total");
                    });

                }

            });
        });

        this.httpServer.requestHandler(this.router);

        int port = Integer.parseInt(ConfigUtil.getPort());
        this.httpServer.listen(port, (server) -> {
            if (server.succeeded()) {
                LOGGER.info("HttpServer is listening on port: " + ((HttpServer)server.result()).actualPort());
            } else {
                LOGGER.error(server.cause().getMessage());
            }
        });

        String finalSubscribeArn = subscribeArn;
        Runtime.getRuntime().addShutdownHook(new Thread(()->{
            try{
                SNSUtil.unSubHTTPS(snsClient, finalSubscribeArn);
            }catch (SnsException e){
                LOGGER.error(e.awsErrorDetails().errorMessage());
            }
            snsClient.close();

            LOGGER.info("shut down!");
        }));

    }

    public String subscribe(SnsClient snsClient, String snsTopicArn, String endPoint, String protocol){
        String subscribeArn = "";
        try {
            subscribeArn =  SNSUtil.subHTTPS(snsClient, snsTopicArn, endPoint, protocol);
        }catch (SnsException e){
            LOGGER.error(e.awsErrorDetails().errorMessage());
            return subscribeArn;
        }
        return subscribeArn;
    }

    public boolean verifySignature(SnsMessageManager manager, SnsClient snsClient, JsonObject jsonObject, String snsTopicArn){
        String messageType = jsonObject.getString("Type");
        String token = jsonObject.getString("Token");
        if(!SNSUtil.verifySignatrue(manager, new ByteArrayInputStream(jsonObject.toBuffer().getBytes()))){
            LOGGER.error("An error occurred while verifying the signature.");
            return false;
        }else{
            //confirm sub or unSub
            LOGGER.info("verify signature successful");
            if (messageType.equals("SubscriptionConfirmation") || messageType.equals("UnsubscribeConfirmation")) {
                try {
                    SNSUtil.confirmSubHTTPS(snsClient, token, snsTopicArn);
                } catch (SnsException e) {
                    LOGGER.error(e.awsErrorDetails().errorMessage());
                    LOGGER.error("an error occurred while confirming subscription");
                    return false;
                }
            }
        }
        return true;
    }

    @Override
    public Adapter getAdapter() {
        return new SnsAdapter();
    }

    public String getTopicArn(){
        return this.snsTopicArn;
    }

    public String getRegion(){
        return this.region;
    }

    public String getEndPoint(){
        return this.endPoint;
    }

    public String getProtocol(){
        return this.protocol;
    }
}
