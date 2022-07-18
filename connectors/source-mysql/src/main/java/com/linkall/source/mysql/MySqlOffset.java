package com.linkall.source.mysql;

import com.linkall.vance.common.env.EnvUtil;

// ["linkall",{"server":"linkall"}] ->
// {"ts_sec":1658130705,"file":"binlog.000010","pos":44602,"gtids":"46ce11bf-9497-11ec-97a5-0e18e1c0f63b:1-75","snapshot":true}
public class MySqlOffset {
  private Integer pos;
  private String file;

  public MySqlOffset() {
    this.file = EnvUtil.getConfig("offset_binlog_file");
    String pos = EnvUtil.getConfig("offset_binlog_pos");
    if (pos != null && !pos.isEmpty()) {
      try {
        this.pos = Integer.parseInt(pos);
      } catch (NumberFormatException e) {
        throw new RuntimeException(e);
      }
    }
  }

  public Integer getPos() {
    return pos;
  }

  public void setPos(Integer pos) {
    this.pos = pos;
  }

  public String getFile() {
    return file;
  }

  public void setFile(String file) {
    this.file = file;
  }
}
