package com.linkall.source.aws.s3;

import com.linkall.cdk.config.Config;
import com.linkall.cdk.connector.Element;
import com.linkall.cdk.connector.Source;
import com.linkall.cdk.connector.Tuple;
import com.linkall.source.aws.utils.AwsHelper;
import com.linkall.source.aws.utils.S3Util;
import com.linkall.source.aws.utils.SQSUtil;
import io.cloudevents.CloudEvent;
import io.vertx.core.json.JsonArray;
import io.vertx.core.json.JsonObject;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import software.amazon.awssdk.regions.Region;
import software.amazon.awssdk.services.s3.S3Client;
import software.amazon.awssdk.services.sqs.SqsClient;
import software.amazon.awssdk.services.sqs.model.Message;

import java.util.List;
import java.util.Objects;
import java.util.concurrent.*;


public class S3Source implements Source {
    private static final Logger LOGGER = LoggerFactory.getLogger(S3Source.class);
    private static final String SQS_NAME = "vanus-connect-s3-source-sqs";

    private BlockingQueue<Tuple> queue;
    private S3Config config;
    private SqsClient sqsClient;
    private String queueUrl;
    private ExecutorService executorService;
    private volatile boolean isRunning = true;

    public S3Source() {
        queue = new LinkedBlockingQueue<>(100);
        executorService = Executors.newSingleThreadExecutor();
    }

    @SuppressWarnings("unchecked")
    public void start() {
        AwsHelper.checkCredentials(config.getSecretConfig().getAccessKeyID(), config.getSecretConfig().getSecretAccessKey());
        // get region
        String strRegion = config.getRegion();
        String sqsArn = config.getSqsArn();
        if (null==strRegion) {
            if (null==sqsArn) {
                LOGGER.error("region and sqsArn cannot be null at the same time.");
                System.exit(1);
            } else {
                strRegion = SQSUtil.getRegion(sqsArn);
            }
        }
        Region region = Region.of(strRegion);

        // get S3Client
        S3Client s3 = Objects.requireNonNull(
                S3Client.builder().region(region).build());
        sqsClient = SqsClient.builder().region(region).build();


        String s3Arn = config.getBucketArn();
        String queUrl = null;
        String bucketName = s3Arn.substring(s3Arn.indexOf(":::") + 3);

        //if sqs_arn is omitted,create a default queue
        if (null==sqsArn || "".equals(sqsArn)) {
            queUrl = SQSUtil.obtainVanceQueueUrl(sqsClient,SQS_NAME);

            // create a default SQS queue if queUrl is null
            if (queUrl==null) {
                //create a vance-sqs queue and return its queUrl
                queUrl = SQSUtil.createQueue(sqsClient, SQS_NAME);
                sqsArn = SQSUtil.getQueueArn(sqsClient, queUrl);
                //construct a policy
                boolean setQueuePolicyOK = SQSUtil.setQueuePolicy(sqsClient, queUrl, SQSUtil.buildPolicy(null, s3Arn, sqsArn));
                if (!setQueuePolicyOK) {
                    LOGGER.error("set sqs policy failed");
                }

                boolean setNotifyConfigOK = S3Util.setNotifyConfig(s3, bucketName, S3Util.buildQueConfig(sqsArn, config.getS3Events()));
                if (!setNotifyConfigOK) {
                    LOGGER.error("set s3 bucket notify configuration failed");
                }

            } else {
                LOGGER.info("vance-sqs existed");
                sqsArn = SQSUtil.getQueueArn(sqsClient, queUrl);
                // Get sqs policy
                innerSetPolicy(s3, sqsClient, sqsArn, s3Arn, queUrl, bucketName);
            }
        } else {
            String queName = sqsArn.substring(sqsArn.lastIndexOf(":") + 1);
            queUrl = SQSUtil.getQueueUrl(sqsClient, queName);
            // Get sqs policy
            innerSetPolicy(s3, sqsClient, sqsArn, s3Arn, queUrl, bucketName);
        }
        s3.close();
        this.queueUrl = queUrl;
        executorService.execute(this::runLoop);
    }

    public void runLoop() {
        S3Adapter s3Adapter = new S3Adapter();
        while (isRunning) {
            List<Message> messages = SQSUtil.receiveLongPollMessages(sqsClient, queueUrl, 15, 5);
            for (Message message : messages) {
                LOGGER.info("[receive S3 events]: " + message.body().toString());
                JsonObject body = new JsonObject(message.body());
                //delete testEvent
                if ("s3:TestEvent".equals(body.getString("Event"))) {
                    SQSUtil.deleteMessage(sqsClient, queueUrl, message);
                    LOGGER.info("[delete s3:TestEvent completed]");
                }
                JsonArray records = body.getJsonArray("Records");
                if (null!=records) {
                    for (int i = 0; i < records.size(); i++) {
                        JsonObject record = records.getJsonObject(i);
                        CloudEvent ce = s3Adapter.adapt(record);
                        Tuple tuple = new Tuple(new Element(ce, record), () -> {
                            LOGGER.info("send event success,{}", ce.getId());
                            SQSUtil.deleteMessage(sqsClient, queueUrl, message);
                            LOGGER.info("[sqs delete message completed]");
                        }, (success, failed, msg) -> {
                            LOGGER.info("send event failed,{},{}", ce.getId(), msg);
                        });
                        try {
                            queue.put(tuple);
                        } catch (InterruptedException e) {
                            LOGGER.warn("put event interrupted");
                        }
                    }
                }
            }
        }
    }

    private void innerSetPolicy(S3Client s3, SqsClient sqsClient, String sqsArn, String s3Arn, String queUrl, String bucketName) {
        JsonObject policy = SQSUtil.getQueuePolicy(sqsClient, queUrl);

        boolean setQueuePolicyOK = SQSUtil.setQueuePolicy(sqsClient, queUrl, SQSUtil.buildPolicy(policy, s3Arn, sqsArn));
        if (!setQueuePolicyOK) {
            LOGGER.error("update sqs policy failed");
        }
        boolean setNotifyConfigOK = S3Util.setNotifyConfig(s3, bucketName, S3Util.buildQueConfig(sqsArn, config.getS3Events()));
        if (!setNotifyConfigOK) {
            LOGGER.error("set s3 bucket notify configuration failed");
        }
    }

    @Override
    public BlockingQueue<Tuple> queue() {
        return queue;
    }

    @Override
    public Class<? extends Config> configClass() {
        return S3Config.class;
    }

    @Override
    public void initialize(Config config) throws Exception {
        this.config = (S3Config) config;
        start();
    }

    @Override
    public String name() {
        return "AmazonS3Source";
    }

    @Override
    public void destroy() throws Exception {
        isRunning = false;
        executorService.shutdown();
        try {
            executorService.awaitTermination(10, TimeUnit.SECONDS);
        } catch (InterruptedException e) {
            LOGGER.error("awaitTermination", e);
        }
        if (sqsClient!=null) {
            sqsClient.close();
        }
    }
}