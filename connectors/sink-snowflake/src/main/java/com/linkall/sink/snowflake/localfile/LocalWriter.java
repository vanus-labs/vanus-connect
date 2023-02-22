package com.linkall.sink.snowflake.localfile;


import com.linkall.sink.snowflake.Writer;

import java.io.*;

public class LocalWriter implements Writer {

    private static final String fileExtension = ".json";
    private File file;
    private PrintWriter writer;
    private long size;

    public LocalWriter() throws IOException {
        file = new File(String.format("%s/%d%s", "data", System.currentTimeMillis(), fileExtension));
        if (!file.exists()) {
            file.getParentFile().mkdirs();
            file.createNewFile();
        }
        writer = new PrintWriter(new BufferedOutputStream(new FileOutputStream(file)));
    }

    @Override
    public void write(byte[] data) {
        writer.println(new String(data));
        size += data.length;
    }

    @Override
    public void flush() {
        writer.flush();
        writer.close();
    }

    @Override
    public void close() {
        file.delete();
    }

    @Override
    public String getFilepath() {
        return file.getAbsolutePath();
    }

    @Override
    public String getFilename() {
        return file.getName();
    }

    @Override
    public long size() {
        return size;
    }


}
