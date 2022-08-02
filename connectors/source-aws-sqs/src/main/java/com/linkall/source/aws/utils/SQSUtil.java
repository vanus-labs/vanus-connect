package com.linkall.source.aws.utils;

import com.linkall.vance.common.file.GenericFileUtil;
import io.vertx.core.json.JsonArray;
import io.vertx.core.json.JsonObject;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import software.amazon.awssdk.services.sqs.SqsClient;
import software.amazon.awssdk.services.sqs.model.*;

import java.util.HashMap;
import java.util.List;
import java.util.Map;

public class SQSUtil {
    private static final Logger LOGGER = LoggerFactory.getLogger(SQSUtil.class);

    public static String getRegion(String sqsArn){
        String[] arr = sqsArn.split(":");
        return arr[3];
    }

    public static boolean deleteMessage(SqsClient sqsClient, String queueUrl, Message message) {
        try {
            DeleteMessageRequest deleteMessageRequest = DeleteMessageRequest.builder()
                    .queueUrl(queueUrl)
                    .receiptHandle(message.receiptHandle())
                    .build();
            DeleteMessageResponse resp = sqsClient.deleteMessage(deleteMessageRequest);
            return resp.sdkHttpResponse().isSuccessful();
        } catch (SqsException e) {
            LOGGER.error(e.awsErrorDetails().errorMessage());
        }
        return false;
    }

    public static List<Message> receiveLongPollMessages(SqsClient sqsClient, String queUrl, Integer waitTime, Integer maxMsgNum) {

        List<Message> messages = null;

        try {

            // Enable long polling on a message receipt.
            ReceiveMessageRequest receiveRequest = ReceiveMessageRequest.builder()
                    .queueUrl(queUrl)
                    .waitTimeSeconds(waitTime)
                    .maxNumberOfMessages(maxMsgNum)
                    .build();

            messages = sqsClient.receiveMessage(receiveRequest).messages();


        } catch (SqsException e) {
            System.err.println(e.awsErrorDetails().errorMessage());
            System.exit(1);
        }
        return messages;
    }

    public static JsonObject getQueuePolicy(SqsClient sqsClient,String queUrl){
        LOGGER.info("====== Get the policy of a SQS queue [" + queUrl+"] Start ======" );
        GetQueueAttributesRequest req = GetQueueAttributesRequest.builder()
                .attributeNames(QueueAttributeName.POLICY)
                .queueUrl(queUrl).build();
        GetQueueAttributesResponse response = sqsClient.getQueueAttributes(req);
        String policy = response.attributes().get(QueueAttributeName.POLICY);
        LOGGER.info("====== Get the policy of a SQS queue [" + queUrl+"] End ======" );
        return new JsonObject(policy);
    }
    /**
     *
     * @param sqsClient
     * @param queueUrl
     * @param policy
     * @return queue Arn
     */
    public static boolean setQueuePolicy(SqsClient sqsClient, String queueUrl, JsonObject policy){

        //System.out.println(policy.encodePrettily());
        LOGGER.info("====== Set the policy of a SQS queue [" + queueUrl+"] Start ======" );
        Map<QueueAttributeName,String> map = new HashMap<>();
        map.put(QueueAttributeName.POLICY,policy.encodePrettily());
        SetQueueAttributesRequest req = SetQueueAttributesRequest.builder()
                .queueUrl(queueUrl)
                .attributes(map)
                .build();
        SetQueueAttributesResponse response = sqsClient.setQueueAttributes(req);
        LOGGER.info("====== Set the policy of a SQS queue [" + queueUrl+"] End ======" );
        return response.sdkHttpResponse().isSuccessful();
    }

    /**
     * get the sqs Arn by queueUrl
     * @param sqsClient
     * @param queueUrl
     * @return
     */
    public static String getQueueArn(SqsClient sqsClient, String queueUrl){
        GetQueueAttributesRequest request = GetQueueAttributesRequest.builder()
                .queueUrl(queueUrl)
                .attributeNames(QueueAttributeName.QUEUE_ARN)
                .build();
        GetQueueAttributesResponse response = sqsClient.getQueueAttributes(request);
        String queArn = response.attributes().get(QueueAttributeName.QUEUE_ARN);
        return  queArn;
    }

    public static List<String> listQueues(SqsClient sqsClient, String prefix) {

        try {
            ListQueuesRequest listQueuesRequest = ListQueuesRequest.builder().queueNamePrefix(prefix).build();
            ListQueuesResponse listQueuesResponse = sqsClient.listQueues(listQueuesRequest);

            if(listQueuesResponse.hasQueueUrls()) return  listQueuesResponse.queueUrls();
            for (String url : listQueuesResponse.queueUrls()) {
                System.out.println(url);
            }

        } catch (SqsException e) {
            System.err.println(e.awsErrorDetails().errorMessage());
            System.exit(1);
        }

        return null;
    }


    /**
     * Create a SQS queue
     * @param sqsClient
     * @param queueName
     * @return the queue Url
     */
    public static String createQueue(SqsClient sqsClient, String queueName) {
        LOGGER.info("====== Create a SQS queue [" + queueName+"] Start ======" );
        CreateQueueRequest createQueueRequest = CreateQueueRequest.builder()
                .queueName(queueName)
                .build();
        try {
            CreateQueueResponse response = sqsClient.createQueue(createQueueRequest);
            LOGGER.info("====== Create a SQS queue ["+queueName+"] successful ======");
            return response.queueUrl();
        } catch (SqsException e) {
            LOGGER.error("====== Create a SQS queue ["+queueName+"] failed ======");
            LOGGER.error(e.awsErrorDetails().errorMessage());
            System.exit(1);
        }
        return null;
    }

    /**
     * get a queue url by queue name
     * @param sqsClient
     * @param queueName
     * @return queUrl
     */
    public static String getQueueUrl(SqsClient sqsClient, String queueName){
        GetQueueUrlResponse getQueueUrlResponse =
                sqsClient.getQueueUrl(GetQueueUrlRequest.builder().queueName(queueName).build());
        String queueUrl = getQueueUrlResponse.queueUrl();
        return queueUrl;
    }
}
