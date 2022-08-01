package com.linkall.sink.aws;

import com.linkall.vance.common.config.SecretUtil;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class AwsHelper {
    private static final Logger LOGGER = LoggerFactory.getLogger(AwsHelper.class);

    /**
     * Check or Create aws credential file
     */
    public static void checkCredentials(){
        LOGGER.info("====== Check aws Credential start ======");
        System.setProperty("aws.accessKeyId", SecretUtil.getString("awsAccessKeyID"));
        System.setProperty("aws.secretAccessKey", SecretUtil.getString("awsSecretAccessKey"));
        LOGGER.info("====== Check aws Credential end ======");
    }
}
