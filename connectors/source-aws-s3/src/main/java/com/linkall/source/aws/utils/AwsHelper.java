package com.linkall.source.aws.utils;

import com.linkall.vance.common.env.SecretUtil;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import software.amazon.awssdk.regions.Region;

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
