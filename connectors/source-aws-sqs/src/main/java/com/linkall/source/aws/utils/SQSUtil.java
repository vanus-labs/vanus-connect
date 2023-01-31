package com.linkall.source.aws.utils;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import software.amazon.awssdk.services.sqs.SqsClient;
import software.amazon.awssdk.services.sqs.model.*;

import java.util.List;

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
