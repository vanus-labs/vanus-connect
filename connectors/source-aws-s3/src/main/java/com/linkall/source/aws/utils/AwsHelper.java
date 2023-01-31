package com.linkall.source.aws.utils;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class AwsHelper {
    private static final Logger LOGGER = LoggerFactory.getLogger(AwsHelper.class);

    /**
     * Check or Create aws credential file
     */
    public static void checkCredentials(String ak, String sk) {
        LOGGER.info("====== Check aws Credential start ======");
        System.setProperty("aws.accessKeyId", ak);
        System.setProperty("aws.secretAccessKey", sk);
        LOGGER.info("====== Check aws Credential end ======");
    }


}
