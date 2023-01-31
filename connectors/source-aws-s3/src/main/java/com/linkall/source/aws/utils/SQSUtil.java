package com.linkall.source.aws.utils;

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

    private static final String SQS_POLICY = "vanus-connect-s3-sqs-policy";


    public static String getRegion(String sqsArn) {
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

    public static JsonObject getQueuePolicy(SqsClient sqsClient, String queUrl) {
        LOGGER.info("====== Get the policy of a SQS queue [" + queUrl + "] Start ======");
        GetQueueAttributesRequest req = GetQueueAttributesRequest.builder()
                .attributeNames(QueueAttributeName.POLICY)
                .queueUrl(queUrl).build();
        GetQueueAttributesResponse response = sqsClient.getQueueAttributes(req);
        String policy = response.attributes().get(QueueAttributeName.POLICY);
        LOGGER.info("====== Get the policy of a SQS queue [" + queUrl + "] End ======");
        return new JsonObject(policy);
    }

    /**
     * @param sqsClient
     * @param queueUrl
     * @param policy
     * @return queue Arn
     */
    public static boolean setQueuePolicy(SqsClient sqsClient, String queueUrl, JsonObject policy) {

        //System.out.println(policy.encodePrettily());
        LOGGER.info("====== Set the policy of a SQS queue [" + queueUrl + "] Start ======");
        Map<QueueAttributeName, String> map = new HashMap<>();
        map.put(QueueAttributeName.POLICY, policy.encodePrettily());
        SetQueueAttributesRequest req = SetQueueAttributesRequest.builder()
                .queueUrl(queueUrl)
                .attributes(map)
                .build();
        SetQueueAttributesResponse response = sqsClient.setQueueAttributes(req);
        LOGGER.info("====== Set the policy of a SQS queue [" + queueUrl + "] End ======");
        return response.sdkHttpResponse().isSuccessful();
    }

    /**
     * get the sqs Arn by queueUrl
     *
     * @param sqsClient
     * @param queueUrl
     * @return
     */
    public static String getQueueArn(SqsClient sqsClient, String queueUrl) {
        GetQueueAttributesRequest request = GetQueueAttributesRequest.builder()
                .queueUrl(queueUrl)
                .attributeNames(QueueAttributeName.QUEUE_ARN)
                .build();
        GetQueueAttributesResponse response = sqsClient.getQueueAttributes(request);
        String queArn = response.attributes().get(QueueAttributeName.QUEUE_ARN);
        return queArn;
    }

    public static List<String> listQueues(SqsClient sqsClient, String prefix) {

        try {
            ListQueuesRequest listQueuesRequest = ListQueuesRequest.builder().queueNamePrefix(prefix).build();
            ListQueuesResponse listQueuesResponse = sqsClient.listQueues(listQueuesRequest);

            if (listQueuesResponse.hasQueueUrls()) {
                return listQueuesResponse.queueUrls();
            }
            for (String url : listQueuesResponse.queueUrls()) {
                LOGGER.error("sqs queue url:{}", url);
            }

        } catch (SqsException e) {
            LOGGER.error("list sqs queue error,{}", e.awsErrorDetails());
            System.exit(1);
        }

        return null;
    }

    /**
     * check the vance-sqs exists or not
     *
     * @return vance-sqs queueUrl if it exists
     */
    public static String obtainVanceQueueUrl(SqsClient sqsClient, String sqsName) {
        List<String> queueUrls = listQueues(sqsClient, sqsName);
        if (queueUrls!=null) {
            return queueUrls.get(0);
        }
        return null;
    }

    /**
     * Create a SQS queue
     *
     * @param sqsClient
     * @param queueName
     * @return the queue Url
     */
    public static String createQueue(SqsClient sqsClient, String queueName) {
        LOGGER.info("====== Create a SQS queue [" + queueName + "] Start ======");
        CreateQueueRequest createQueueRequest = CreateQueueRequest.builder()
                .queueName(queueName)
                .build();
        try {
            CreateQueueResponse response = sqsClient.createQueue(createQueueRequest);
            LOGGER.info("====== Create a SQS queue [" + queueName + "] successful ======");
            return response.queueUrl();
        } catch (SqsException e) {
            LOGGER.error("====== Create a SQS queue [" + queueName + "] failed ======");
            LOGGER.error(e.awsErrorDetails().errorMessage());
            System.exit(1);
        }
        return null;
    }

    /**
     * get a queue url by queue name
     *
     * @param sqsClient
     * @param queueName
     * @return queUrl
     */
    public static String getQueueUrl(SqsClient sqsClient, String queueName) {
        GetQueueUrlResponse getQueueUrlResponse =
                sqsClient.getQueueUrl(GetQueueUrlRequest.builder().queueName(queueName).build());
        String queueUrl = getQueueUrlResponse.queueUrl();
        return queueUrl;
    }

    public static JsonObject buildPolicy(JsonObject existedPolicy, String s3Arn, String sqsArn) {
        if (null==existedPolicy) {
            JsonObject quePolicy = new JsonObject(FileUtil.readResource("queue_policy.json"));
            quePolicy.put("Id", SQS_POLICY);
            JsonObject statement = quePolicy.getJsonArray("Statement").getJsonObject(0);
            buildStatement(statement, s3Arn, sqsArn);
            return quePolicy;
        } else {
            JsonArray statements = existedPolicy.getJsonArray("Statement");
            boolean findSameStatement = false;
            for (int i = 0; i < statements.size(); i++) {
                if (statements.getJsonObject(i).getString("Sid").equals(s3Arn)) {
                    findSameStatement = true;
                    continue;
                }
            }
            // if can't find same statement, just create a new one and add it to statements
            if (!findSameStatement) {
                JsonObject statement = new JsonObject();
                JsonObject condition = new JsonObject();
                statement.put("Condition", condition);
                JsonObject arnEquals = new JsonObject();
                condition.put("ArnEquals", arnEquals);
                buildStatement(statement, s3Arn, sqsArn);
                statements.add(statement);
            }
            return existedPolicy;
        }

    }

    private static void buildStatement(JsonObject statement, String s3Arn, String sqsArn) {
        statement.put("Sid", s3Arn);
        statement.put("Resource", sqsArn);
        statement.put("Effect", "Allow");
        JsonObject principal = new JsonObject();
        principal.put("Service", "s3.amazonaws.com");
        statement.put("Principal", principal);
        statement.put("Action", "sqs:SendMessage");
        JsonObject condition = statement.getJsonObject("Condition");
        //JsonObject stringEquals = condition.getJsonObject("StringEquals");
        JsonObject arnEquals = condition.getJsonObject("ArnEquals");
        //stringEquals.put("aws:SourceAccount",accountId);
        arnEquals.put("aws:SourceArn", s3Arn);
    }

    public static void putVancePolicy() {

    }
}
