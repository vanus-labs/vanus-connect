package com.linkall.source.aws.utils;

import com.amazonaws.SdkClientException;
import com.amazonaws.arn.Arn;
import com.amazonaws.services.sns.message.SnsMessageManager;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import software.amazon.awssdk.services.sns.SnsClient;
import software.amazon.awssdk.services.sns.model.*;

import java.io.InputStream;


public class SNSUtil {

    private static final Logger LOGGER = LoggerFactory.getLogger(SNSUtil.class);

    public static String getRegion(String snsTopicArn) {
        return Arn.fromString(snsTopicArn).getRegion();
    }

    public static String subHTTPS(SnsClient snsClient, String topicArn, String url, String protocol) throws SnsException {
        String subscriptionArn = "";
        SubscribeRequest request = SubscribeRequest.builder()
                .protocol(protocol)
                .endpoint(url)
                .returnSubscriptionArn(true)
                .topicArn(topicArn)
                .build();

        SubscribeResponse result = snsClient.subscribe(request);
        LOGGER.info("subscribe ARN is {}, status:{}", result.subscriptionArn(), result.sdkHttpResponse().statusCode());
        subscriptionArn = result.subscriptionArn();
        return subscriptionArn;
    }

    public static void confirmSubHTTPS(SnsClient snsClient, String subscriptionToken, String topicArn) throws SnsException {
        ConfirmSubscriptionRequest request = ConfirmSubscriptionRequest.builder()
                .token(subscriptionToken)
                .topicArn(topicArn)
                .build();

        ConfirmSubscriptionResponse result = snsClient.confirmSubscription(request);
        LOGGER.info("confirm subscription ARN is {}, status:{}", result.subscriptionArn(), result.sdkHttpResponse().statusCode());
    }

    public static void unSubHTTPS(SnsClient snsClient, String subscriptionArn) throws SnsException {

        UnsubscribeRequest request = UnsubscribeRequest.builder()
                .subscriptionArn(subscriptionArn)
                .build();

        UnsubscribeResponse result = snsClient.unsubscribe(request);

        LOGGER.info("unsubscribe ARN is {}, status:{}", subscriptionArn, result.sdkHttpResponse().statusCode());

    }

    public static boolean verifySignatrue(InputStream message, String region) {
        SnsMessageManager manager = new SnsMessageManager(region);
        boolean verifyResult = true;
        try {
            manager.parseMessage(message);
        } catch (SdkClientException e) {
            LOGGER.error("The signature is not legal", e);
            verifyResult = false;
        }
        return verifyResult;
    }

}
