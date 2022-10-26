package com.linkall.source.mysql;

import com.linkall.vance.common.config.ConfigUtil;

// ["linkall",{"server":"linkall"}] ->
// {"ts_sec":1658130705,"file":"binlog.000010","pos":44602,"gtids":"46ce11bf-9497-11ec-97a5-0e18e1c0f63b:1-75","snapshot":true}
public class MySqlOffset {
  private Integer pos;
  private String file;
  private String gtids;

  public MySqlOffset() {
    this.gtids = ConfigUtil.getString("offset_binlog_gtids");
    this.file = ConfigUtil.getString("offset_binlog_file");
    String pos = ConfigUtil.getString("offset_binlog_pos");
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

  public String getGtids() {
    return gtids;
  }

  public void setGtids(String gtids) {
    this.gtids = gtids;
  }
}
