package com.linkall.source.mysql;

import com.linkall.vance.common.config.ConfigUtil;
import com.linkall.vance.core.VanceApplication;

import java.io.File;

public class Entrance {
  public static void main(String[] args) throws Exception {
    String storeFile = ConfigUtil.getEnvOrConfigOrDefault("v_store_file");
    if (storeFile != null && storeFile != "") {
      File f = new File(storeFile);
      if (!f.exists()) {
        f.createNewFile();
      }
    }
    VanceApplication.run(MySqlSource.class);
  }
}
