package com.linkall.source.aws.utils;

import com.amazonaws.SdkClientException;
import com.amazonaws.services.sns.message.SnsMessage;
import com.amazonaws.services.sns.message.SnsMessageManager;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import software.amazon.awssdk.services.sns.SnsClient;
import software.amazon.awssdk.services.sns.model.*;
import java.io.InputStream;


public class SNSUtil {

    private static final Logger LOGGER = LoggerFactory.getLogger(SNSUtil.class);

    public static  String subscriptionArn;

    public static String getRegion(String snsTopicArn){
        String[] arr = snsTopicArn.split(":");
        return arr[3];
    }

    public static void subHTTPS(SnsClient snsClient, String topicArn, String url, String protocol) {
        try {
            SubscribeRequest request = SubscribeRequest.builder()
                    .protocol(protocol)
                    .endpoint(url)
                    .returnSubscriptionArn(true)
                    .topicArn(topicArn)
                    .build();

            SubscribeResponse result = snsClient.subscribe(request);
            LOGGER.info("Subscription ARN is " + result.subscriptionArn() + "\n\n Status is " + result.sdkHttpResponse().statusCode());
            subscriptionArn = result.subscriptionArn();
        } catch (SnsException e) {
            LOGGER.error(e.awsErrorDetails().errorMessage());
            System.exit(1);
        }
    }

    public static void confirmSubHTTPS(SnsClient snsClient, String subscriptionToken, String topicArn ) {
        try {
            ConfirmSubscriptionRequest request = ConfirmSubscriptionRequest.builder()
                    .token(subscriptionToken)
                    .topicArn(topicArn)
                    .build();

            ConfirmSubscriptionResponse result = snsClient.confirmSubscription(request);
            LOGGER.info("\n\nStatus was " + result.sdkHttpResponse().statusCode() + "\n\nSubscription Arn: \n\n" + result.subscriptionArn());

        } catch (SnsException e) {
            LOGGER.error(e.awsErrorDetails().errorMessage());
            System.exit(1);
        }
    }

    public static void unSubHTTPS(SnsClient snsClient){
        unSubHTTPS(snsClient, subscriptionArn);
    }

    public static void unSubHTTPS(SnsClient snsClient, String subscriptionArn) {

        try {
            UnsubscribeRequest request = UnsubscribeRequest.builder()
                    .subscriptionArn(subscriptionArn)
                    .build();

            UnsubscribeResponse result = snsClient.unsubscribe(request);

            System.out.println("\n\nStatus was " + result.sdkHttpResponse().statusCode()
                    + "\n\nSubscription was removed for " + request.subscriptionArn());

        } catch (SnsException e) {
            System.out.println(e.awsErrorDetails().errorMessage());
            System.exit(1);
        }
    }

    public static boolean verifySignatrue(InputStream message, String region){
        SnsMessageManager manager = new SnsMessageManager(region);
        boolean verifyResult = true;
        try{
            manager.parseMessage(message);
        }catch(SdkClientException e){
            e.printStackTrace();
            LOGGER.error("The signature is not legal!");
            verifyResult = false;
        }
        return verifyResult;
    }

}
