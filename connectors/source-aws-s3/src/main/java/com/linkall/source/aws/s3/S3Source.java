package com.linkall.source.aws.s3;

import com.linkall.source.aws.utils.AwsHelper;
import com.linkall.source.aws.utils.S3Util;
import com.linkall.source.aws.utils.SQSUtil;
import com.linkall.vance.common.env.ConfigUtil;
import com.linkall.vance.common.json.JsonMapper;
import com.linkall.vance.core.Adapter;
import com.linkall.vance.core.Source;
import io.cloudevents.CloudEvent;
import io.cloudevents.http.vertx.VertxMessageFactory;
import io.cloudevents.jackson.JsonFormat;
import io.vertx.core.Future;
import io.vertx.core.Vertx;
import io.vertx.core.buffer.Buffer;
import io.vertx.core.json.JsonArray;
import io.vertx.core.json.JsonObject;
import io.vertx.ext.web.client.HttpResponse;
import io.vertx.ext.web.client.WebClient;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import software.amazon.awssdk.regions.Region;
import software.amazon.awssdk.services.s3.S3Client;
import software.amazon.awssdk.services.sqs.SqsClient;
import software.amazon.awssdk.services.sqs.model.Message;

import java.util.List;
import java.util.Objects;


public class S3Source implements Source {
    private static final Logger LOGGER = LoggerFactory.getLogger(S3Source.class);
    private static final String SQS_NAME = "vance-s3-source-sqs-axzvf";
    private final static Vertx vertx = Vertx.vertx();
    private final static WebClient webClient = WebClient.create(vertx);


    @SuppressWarnings("unchecked")
    @Override
    public void start(){
        AwsHelper.checkCredentials();
        S3Adapter s3Adapter = (S3Adapter) getAdapter();
        // get region
        String strRegion = ConfigUtil.getString("region");
        String sqsArn = ConfigUtil.getString("sqs_arn");
        if(null == strRegion){
            if(null == sqsArn){
                LOGGER.error("region and sqsArn cannot be null at the same time.");
                System.exit(1);
            }else{
                strRegion = SQSUtil.getRegion(sqsArn);
            }
        }
        Region region = Region.of(strRegion);

        // get S3Client
        S3Client s3 = Objects.requireNonNull(
                S3Client.builder().region(region).build());
        SqsClient sqsClient = SqsClient.builder().region(region).build();


        String s3Arn = ConfigUtil.getString("s3_bucket_arn");
        String v_target = ConfigUtil.getVanceSink();
        String queUrl = null;
        String bucketName = s3Arn.substring(s3Arn.indexOf(":::")+3);

        //if sqs_arn is omitted,create a default queue
        if(null==sqsArn||"".equals(sqsArn)){
            queUrl = SQSUtil.obtainVanceQueueUrl(sqsClient);

            // create a default SQS queue if queUrl is null
            if(queUrl ==null){
                //create a vance-sqs queue and return its queUrl
                queUrl = SQSUtil.createQueue(sqsClient,SQS_NAME);
                sqsArn = SQSUtil.getQueueArn(sqsClient,queUrl);
                //construct a policy
                boolean setQueuePolicyOK = SQSUtil.setQueuePolicy(sqsClient,queUrl, SQSUtil.buildPolicy(null,s3Arn,sqsArn));
                if(!setQueuePolicyOK){
                    LOGGER.error("set sqs policy failed");
                }

                boolean setNotifyConfigOK = S3Util.setNotifyConfig(s3,bucketName, S3Util.buildQueConfig(sqsArn));
                if(!setNotifyConfigOK){
                    LOGGER.error("set s3 bucket notify configuration failed");
                }

            }else{
                LOGGER.info("vance-sqs existed");
                sqsArn = SQSUtil.getQueueArn(sqsClient,queUrl);
                // Get sqs policy
                innerSetPolicy(s3, sqsClient, sqsArn, s3Arn, queUrl, bucketName);
            }
        }else{
            String queName = sqsArn.substring(sqsArn.lastIndexOf(":")+1);
            queUrl = SQSUtil.getQueueUrl(sqsClient,queName);
            // Get sqs policy
            innerSetPolicy(s3, sqsClient, sqsArn, s3Arn, queUrl, bucketName);
            //JsonObject policy = SQSUtil.getQueuePolicy(sqsClient,queUrl);
        }
        s3.close();
        Runtime.getRuntime().addShutdownHook(new Thread(()->{
            if(null!=sqsClient) sqsClient.close();
        }));
        final String qurl =queUrl;
        while (true){
            List<Message> messages = SQSUtil.receiveLongPollMessages(sqsClient,queUrl,15,5);
            for (Message message : messages) {
                //System.out.println(message.body());
                LOGGER.info("[receive S3 events]: "+message.body().toString());
                JsonObject body = new JsonObject(message.body());
                //delete testEvent
                if("s3:TestEvent".equals(body.getString("Event"))){
                    SQSUtil.deleteMessage(sqsClient,qurl,message);
                    LOGGER.info("[delete s3:TestEvent completed]");
                }
                JsonArray records = body.getJsonArray("Records");
                if(null!=records){
                    for (int i = 0; i < records.size(); i++) {
                        CloudEvent ce = s3Adapter.adapt(records.getJsonObject(i));

                        Future<HttpResponse<Buffer>> responseFuture = VertxMessageFactory.createWriter(webClient.postAbs(v_target))
                                .writeStructured(ce, JsonFormat.CONTENT_TYPE);
                        responseFuture
                                .onSuccess((resp)->{
                                    //ret.set(resp.bodyAsString());
                                    JsonObject sendObj = JsonMapper.wrapCloudEvent(ce);
                                    LOGGER.info("[deliver cloud events]: "+ sendObj.getString("type")+" "+sendObj.getJsonObject("data").getJsonObject("object").getString("key"));
                                    LOGGER.info("[response: "+resp.bodyAsString()+"]");
                                    if(resp.statusCode()==200){
                                        SQSUtil.deleteMessage(sqsClient,qurl,message);
                                        LOGGER.info("[sqs delete message completed]");
                                    }
                                }) // Print the received message
                                .onFailure(System.err::println);
                    }
                }
            }
        }

    }

    private void innerSetPolicy(S3Client s3, SqsClient sqsClient, String sqsArn, String s3Arn, String queUrl, String bucketName) {
        JsonObject policy = SQSUtil.getQueuePolicy(sqsClient,queUrl);

        boolean setQueuePolicyOK = SQSUtil.setQueuePolicy(sqsClient,queUrl, SQSUtil.buildPolicy(policy,s3Arn,sqsArn));
        if(!setQueuePolicyOK){
            LOGGER.error("update sqs policy failed");
        }
        boolean setNotifyConfigOK = S3Util.setNotifyConfig(s3,bucketName, S3Util.buildQueConfig(sqsArn));
        if(!setNotifyConfigOK){
            LOGGER.error("set s3 bucket notify configuration failed");
        }
    }

    @Override
    public Adapter getAdapter() {
        return new S3Adapter();
    }
}