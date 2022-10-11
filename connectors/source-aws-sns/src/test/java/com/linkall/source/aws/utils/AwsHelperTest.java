package com.linkall.source.aws.utils;

import com.linkall.vance.common.config.SecretUtil;
import junit.framework.TestCase;
import org.junit.Test;

import static org.junit.Assert.*;

public class AwsHelperTest extends TestCase {

    @Test
    public void testCheckCredentials() {
        AwsHelper.checkCredentials();
        assertEquals(SecretUtil.getString("awsAccessKeyID"), System.getProperty("aws.accessKeyId"));
        assertEquals(SecretUtil.getString("awsSecretAccessKey"), System.getProperty("aws.secretAccessKey"));
    }
}