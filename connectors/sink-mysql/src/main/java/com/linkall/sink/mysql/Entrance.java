package com.linkall.sink.mysql;

import com.linkall.vance.core.VanceApplication;

public class Entrance {
  public static void main(String[] args) {
    VanceApplication.run(MySqlSink.class);
  }
}
