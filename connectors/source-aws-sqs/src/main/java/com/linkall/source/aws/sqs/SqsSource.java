package com.linkall.source.aws.sqs;

import com.linkall.source.aws.utils.AwsHelper;
import com.linkall.source.aws.utils.SQSUtil;
import com.linkall.vance.common.config.ConfigUtil;
import com.linkall.vance.common.json.JsonMapper;
import com.linkall.vance.core.Adapter;
import com.linkall.vance.core.Source;
import com.linkall.vance.core.http.HttpClient;
import io.cloudevents.CloudEvent;
import io.cloudevents.http.vertx.VertxMessageFactory;
import io.cloudevents.jackson.JsonFormat;
import io.vertx.core.Future;
import io.vertx.core.Vertx;
import io.vertx.core.buffer.Buffer;
import io.vertx.core.json.JsonObject;
import io.vertx.ext.web.client.HttpResponse;
import io.vertx.ext.web.client.WebClient;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import software.amazon.awssdk.regions.Region;
import software.amazon.awssdk.services.sqs.SqsClient;
import software.amazon.awssdk.services.sqs.model.Message;

import java.util.List;
import java.util.concurrent.atomic.AtomicInteger;


public class SqsSource implements Source {
    private static final Logger LOGGER = LoggerFactory.getLogger(SqsSource.class);
    private final static Vertx vertx = Vertx.vertx();
    private final static WebClient webClient = WebClient.create(vertx);

    @SuppressWarnings("unchecked")
    @Override
    public void start(){
        AwsHelper.checkCredentials();
        SqsAdapter sqsAdapter = (SqsAdapter) getAdapter();
        String sqsArn = ConfigUtil.getString("sqs_arn");
        String strRegion = SQSUtil.getRegion(sqsArn);

        // get region
        Region region = Region.of(strRegion);
        SqsClient sqsClient = SqsClient.builder().region(region).build();

        String v_target = ConfigUtil.getVanceSink();
        String queUrl = null;
        String queName = sqsArn.substring(sqsArn.lastIndexOf(":")+1);
        queUrl = SQSUtil.getQueueUrl(sqsClient,queName);
        final String qurl = queUrl;
        Runtime.getRuntime().addShutdownHook(new Thread(()->{
            if(null!=sqsClient) sqsClient.close();
        }));
        //final String qurl =queUrl;
        while (true){
            List<Message> messages = SQSUtil.receiveLongPollMessages(sqsClient,queUrl,15,5);
            for (Message message : messages) {

                LOGGER.info("[receive SQS msg]: "+message);
                SqsContent sqsContent = new SqsContent(message.messageId(),message.body(),strRegion,queName);
                CloudEvent ce = sqsAdapter.adapt(sqsContent);
                Future<HttpResponse<Buffer>> responseFuture = VertxMessageFactory.createWriter(webClient.postAbs(v_target))
                        .writeStructured(ce, JsonFormat.CONTENT_TYPE);
                responseFuture
                        .onSuccess((resp)->{
                            LOGGER.info("[response: "+resp.bodyAsString()+"]");
                            SQSUtil.deleteMessage(sqsClient,qurl,message);
                            LOGGER.info("[sqs delete message completed]");

                        }) // Print the received message
                        .onFailure(System.err::println);
            }
        }

    }
    @Override
    public Adapter getAdapter() {
        return new SqsAdapter();
    }
}