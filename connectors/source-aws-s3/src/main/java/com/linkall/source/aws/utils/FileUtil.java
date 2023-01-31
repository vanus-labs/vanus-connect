package com.linkall.source.aws.utils;


import java.io.IOException;

public class FileUtil {
    public static String readResource(String name) {
        byte[] b;
        try {
            b = FileUtil.class.getClassLoader().getResourceAsStream(name).readAllBytes();
        } catch (IOException e) {
            throw new RuntimeException(e);
        }
        return new String(b);
    }
}
