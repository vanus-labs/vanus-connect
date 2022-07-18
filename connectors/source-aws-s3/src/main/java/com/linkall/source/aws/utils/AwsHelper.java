package com.linkall.source.aws.utils;

import com.linkall.vance.common.env.SecretUtil;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import software.amazon.awssdk.regions.Region;

public class AwsHelper {
    private static final Logger LOGGER = LoggerFactory.getLogger(AwsHelper.class);

    /**
     * transform a region name to a Region entity
     * @param name
     * @return
     */
    public static Region getRegion(String name){
        switch (name){
            case "ap-south-1":
                return Region.AP_SOUTH_1;
            case "eu-south-1" :
                return Region.EU_SOUTH_1;
            case "us-gov-east-1":
                return Region.US_GOV_EAST_1;
            case "ca-central-1":
                return Region.CA_CENTRAL_1;
            case "eu-central-1" :
                return Region.EU_CENTRAL_1;
            case "us-west-1":
                return Region.US_WEST_1;
            case "us-west-2":
                return Region.US_WEST_2;
            case "af-south-1" :
                return Region.AF_SOUTH_1;
            case "eu-north-1":
                return Region.EU_NORTH_1;
            case "eu-west-3":
                return Region.EU_WEST_3;
            case "eu-west-2" :
                return Region.EU_WEST_2;
            case "eu-west-1":
                return Region.EU_WEST_1;
            case "ap-northeast-3":
                return Region.AP_NORTHEAST_3;
            case "ap-northeast-2" :
                return Region.AP_NORTHEAST_2;
            case "ap-northeast-1":
                return Region.AP_NORTHEAST_1;
            case "me-south-1":
                return Region.ME_SOUTH_1;
            case "sa-east-1" :
                return Region.SA_EAST_1;
            case "ap-east-1":
                return Region.AP_EAST_1;
            case "cn-north-1":
                return Region.CN_NORTH_1;
            case "us-gov-west-1" :
                return Region.US_GOV_WEST_1;
            case "ap-southeast-1":
                return Region.AP_SOUTHEAST_1;
            case "ap-southeast-2":
                return Region.AP_SOUTHEAST_2;
            case "us-iso-east-1" :
                return Region.US_ISO_EAST_1;
            case "us-east-1":
                return Region.US_EAST_1;
            case "us-east-2":
                return Region.US_EAST_2;
            case "cn-northwest-1" :
                return Region.CN_NORTHWEST_1;
            case "us-isob-east-1":
                return Region.US_ISOB_EAST_1;
            case "aws-global":
                return Region.AWS_GLOBAL;
            case "aws-cn-global" :
                return Region.AWS_CN_GLOBAL;
            case "aws-us-gov-global":
                return Region.AWS_US_GOV_GLOBAL;
            case "aws-iso-global":
                return Region.AWS_ISO_GLOBAL;
            case "aws-iso-b-global" :
                return Region.AWS_ISO_B_GLOBAL;
            default: return null;
        }
    }

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
