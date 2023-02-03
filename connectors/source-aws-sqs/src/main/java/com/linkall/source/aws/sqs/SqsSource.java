package com.linkall.source.aws.sqs;

import com.linkall.cdk.config.Config;
import com.linkall.cdk.connector.Element;
import com.linkall.cdk.connector.Source;
import com.linkall.cdk.connector.Tuple;
import com.linkall.source.aws.utils.AwsHelper;
import com.linkall.source.aws.utils.SQSUtil;
import io.cloudevents.CloudEvent;
import io.vertx.core.Vertx;
import io.vertx.ext.web.client.WebClient;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import software.amazon.awssdk.regions.Region;
import software.amazon.awssdk.services.sqs.SqsClient;
import software.amazon.awssdk.services.sqs.model.Message;

import java.util.List;
import java.util.concurrent.*;


public class SqsSource implements Source {
    private static final Logger LOGGER = LoggerFactory.getLogger(SqsSource.class);
    private final static Vertx vertx = Vertx.vertx();
    private final static WebClient webClient = WebClient.create(vertx);

    private BlockingQueue<Tuple> queue;
    private SqsConfig config;

    private SqsClient sqsClient;
    private String queueUrl;
    private String region;
    private String queueName;
    private ExecutorService executorService;
    private volatile boolean isRunning = true;

    public SqsSource() {
        queue = new LinkedBlockingQueue<>(100);
        executorService = Executors.newSingleThreadExecutor();
    }

    @SuppressWarnings("unchecked")
    public void start() {
        AwsHelper.checkCredentials(config.getSecretConfig().getAccessKeyID(), config.getSecretConfig().getSecretAccessKey());
        String sqsArn = config.getSqsArn();
        region = SQSUtil.getRegion(sqsArn);

        // get region
        sqsClient = SqsClient.builder().region(Region.of(region)).build();

        queueName = sqsArn.substring(sqsArn.lastIndexOf(":") + 1);
        queueUrl = SQSUtil.getQueueUrl(sqsClient, queueName);

        executorService.execute(this::runLoop);
    }

    public void runLoop() {
        SqsAdapter sqsAdapter = new SqsAdapter();
        while (isRunning) {
            List<Message> messages = SQSUtil.receiveLongPollMessages(sqsClient, queueUrl, 15, 5);
            for (Message message : messages) {
                LOGGER.info("[receive SQS msg]:{},{} ", message.messageId(), message.body());
                SqsContent sqsContent = new SqsContent(message.messageId(), message.body(), region, queueName);
                CloudEvent ce = sqsAdapter.adapt(sqsContent);
                Tuple tuple = new Tuple(new Element(ce, sqsContent), () -> {
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

    @Override
    public BlockingQueue<Tuple> queue() {
        return queue;
    }

    @Override
    public Class<? extends Config> configClass() {
        return SqsConfig.class;
    }

    @Override
    public void initialize(Config config) throws Exception {
        this.config = (SqsConfig) config;
        start();
    }

    @Override
    public String name() {
        return "AmazonSQSSource";
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