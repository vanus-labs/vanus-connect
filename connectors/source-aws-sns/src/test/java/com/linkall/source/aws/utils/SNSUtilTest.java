package com.linkall.source.aws.utils;

import com.amazonaws.SdkClientException;

import com.amazonaws.services.sns.message.SnsMessageManager;

import io.vertx.core.json.JsonObject;

import org.junit.AfterClass;

import org.junit.BeforeClass;
import org.junit.Test;
import org.mockito.Mockito;
import software.amazon.awssdk.http.SdkHttpResponse;
import software.amazon.awssdk.services.sns.SnsClient;
import software.amazon.awssdk.services.sns.model.*;

import java.io.ByteArrayInputStream;
import java.io.InputStream;
import java.util.List;
import java.util.Map;
import java.util.Optional;
import java.util.UUID;

import static org.junit.Assert.*;

public class SNSUtilTest{

    private static SnsClient snsClient;
    private static String snsTopicArn = "arn:aws:sns:us-west-2:843378899134:TestTopic";
    private static String region = "us-west-2";
    private static String endpoint = "https://105a-2408-8207-2537-f360-65b6-8b9a-7ed0-8db1.jp.ngrok.io";
    private static String protocol = "https";
    private static String subscriptionArn = snsTopicArn+":abc";

    @BeforeClass
    public static void beforeInit(){
        snsClient = Mockito.mock(SnsClient.class);
    }

    @Test
    public void testGetRegion() {
        String retRegion = SNSUtil.getRegion(snsTopicArn);
        assertEquals(region, retRegion);
    }

    @Test
    public void testSubHTTPS() {
        SubscribeRequest request = SubscribeRequest.builder()
                .protocol(protocol)
                .endpoint(endpoint)
                .returnSubscriptionArn(true)
                .topicArn(snsTopicArn)
                .build();
        SubscribeResponse response = (SubscribeResponse) SubscribeResponse.builder()
                .subscriptionArn(subscriptionArn)
                .sdkHttpResponse(new SdkHttpResponse() {
                    @Override
                    public Optional<String> statusText() {
                        return Optional.empty();
                    }

                    @Override
                    public int statusCode() {
                        return 200;
                    }

                    @Override
                    public Map<String, List<String>> headers() {
                        return null;
                    }

                    @Override
                    public Builder toBuilder() {
                        return null;
                    }
                })
                .build();
        Mockito.when(snsClient.subscribe(request)).thenReturn(response);
        String subHTTPS = SNSUtil.subHTTPS(snsClient, snsTopicArn, endpoint, protocol);
        assertEquals(subscriptionArn, subHTTPS);
    }

    @Test
    public void testConfirmSubHTTPS() {
        String subscriptionToken = UUID.randomUUID().toString();
        ConfirmSubscriptionRequest request = ConfirmSubscriptionRequest.builder()
                .token(subscriptionToken)
                .topicArn(snsTopicArn)
                .build();
        ConfirmSubscriptionResponse response = (ConfirmSubscriptionResponse) ConfirmSubscriptionResponse.builder()
                .subscriptionArn(subscriptionArn)
                .sdkHttpResponse(new SdkHttpResponse() {
                    @Override
                    public Optional<String> statusText() {
                        return Optional.empty();
                    }

                    @Override
                    public int statusCode() {
                        return 200;
                    }

                    @Override
                    public Map<String, List<String>> headers() {
                        return null;
                    }

                    @Override
                    public Builder toBuilder() {
                        return null;
                    }
                })
                .build();
        Mockito.when(snsClient.confirmSubscription(request)).thenReturn(response);
        SNSUtil.confirmSubHTTPS(snsClient, subscriptionToken, snsTopicArn);
    }

    @Test
    public void testUnSubHTTPS() {
        UnsubscribeRequest request = UnsubscribeRequest.builder()
                .subscriptionArn(subscriptionArn)
                .build();
        UnsubscribeResponse response = (UnsubscribeResponse) UnsubscribeResponse.builder()
                .sdkHttpResponse(new SdkHttpResponse() {
                    @Override
                    public Optional<String> statusText() {
                        return Optional.empty();
                    }

                    @Override
                    public int statusCode() {
                        return 200;
                    }

                    @Override
                    public Map<String, List<String>> headers() {
                        return null;
                    }

                    @Override
                    public Builder toBuilder() {
                        return null;
                    }
                })
                .build();
        Mockito.when(snsClient.unsubscribe(request)).thenReturn(response);
        SNSUtil.unSubHTTPS(snsClient, subscriptionArn);
    }

    @Test
    public void testVerifySignatrue() {
        SnsMessageManager manager = Mockito.mock(SnsMessageManager.class);
        JsonObject jsonObject = new JsonObject();
        jsonObject.put("Type","Notification");
        jsonObject.put("MessageId","arn:aws:sns:us-west-2:123456789012:MyTopic");
        jsonObject.put("Subject","My First Message");
        jsonObject.put("Message","Hello world!");
        jsonObject.put("Timestamp","2012-05-02T00:54:06.655Z");
        jsonObject.put("Signature","EXAMPLEw6JRN...");
        InputStream inputStream = new ByteArrayInputStream(jsonObject.toBuffer().getBytes());
        Mockito.when(manager.parseMessage(inputStream)).
                thenThrow(new SdkClientException("An error occurred while verifying the signature."));
        boolean ret = SNSUtil.verifySignatrue(manager, inputStream);
        assertEquals(false, ret);
    }


    @AfterClass
    public static void close(){
        snsClient.close();
    }
}