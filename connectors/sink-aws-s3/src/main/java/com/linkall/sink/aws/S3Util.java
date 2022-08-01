package com.linkall.sink.aws;

import com.linkall.vance.common.config.ConfigUtil;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import software.amazon.awssdk.core.sync.RequestBody;
import software.amazon.awssdk.services.s3.S3Client;
import software.amazon.awssdk.services.s3.model.*;

import java.io.File;
import java.net.URL;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

/**
 * S3Util encapsulate a set of common operation for S3
 */
public class S3Util {
    private static final Logger LOGGER = LoggerFactory.getLogger(S3Util.class);
    private static Map<String,String> metaData = new HashMap<>();
    private static final String SQS_NOTIFY_ID = "vance-s3-notification";

    public static boolean setNotifyConfig(S3Client s3, String bucketName, QueueConfiguration qConfig){
        LOGGER.info("====== Set the Notification Configuration of the bucket [" + bucketName+"] Start ======" );
        PutBucketNotificationConfigurationResponse resp = s3.putBucketNotificationConfiguration(PutBucketNotificationConfigurationRequest.builder()
                .bucket(bucketName)
                .notificationConfiguration(NotificationConfiguration.builder()
                        .queueConfigurations(qConfig).build())
                .build());
        LOGGER.info("====== Set the Notification Configuration of the bucket [" + bucketName+"] End ======" );
        return resp.sdkHttpResponse().isSuccessful();
    }

    public static QueueConfiguration buildQueConfig(String sqsArn){
        List<String> events = ConfigUtil.getStringArray("s3_events");
        QueueConfiguration qConfig = QueueConfiguration.builder().queueArn(sqsArn)
                .id(SQS_NOTIFY_ID)
                .eventsWithStrings(events)
                .build();
        return qConfig;
    }

    /**
     * return the url of an object with specific @param {keyName}
     * @param s3
     * @param bucketName
     * @param keyName
     * @return
     */
    public static URL getURL(S3Client s3, String bucketName, String keyName ) {

        try {
            GetUrlRequest request = GetUrlRequest.builder()
                    .bucket(bucketName)
                    .key(keyName)
                    .build();

            return s3.utilities().getUrl(request);

        } catch (S3Exception e) {
            System.err.println(e.awsErrorDetails().errorMessage());
            return null;
        }
    }


    public static boolean putS3Object(S3Client s3,
                                    String bucketName,
                                    String objectKey,
                                    File f ) {
        try {
            PutObjectRequest putOb = PutObjectRequest.builder()
                    .bucket(bucketName)
                    .key(objectKey)
                    .build();

            PutObjectResponse response = s3.putObject(putOb,
                    RequestBody.fromFile(f));

            return response.sdkHttpResponse().isSuccessful();

        } catch (S3Exception e) {
            System.err.println(e.getMessage());
            return false;
        }
    }

    public static boolean putS3Object(S3Client s3,
                                     String bucketName,
                                     String objectKey,
                                     byte[] data ) {

        try {
            PutObjectRequest putOb = PutObjectRequest.builder()
                    .bucket(bucketName)
                    .key(objectKey)
                    .build();

            PutObjectResponse response = s3.putObject(putOb,
                    RequestBody.fromBytes(data));

            return response.sdkHttpResponse().isSuccessful();

        } catch (S3Exception e) {
            LOGGER.error(e.getMessage());
            return false;
        }
    }


}
