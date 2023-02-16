package com.linkall.sink.snowflake;


public interface Writer {

    void write(byte[] data);

    String getFilepath();

    String getFilename();

    long size();

    void flush();

    void close();

}
