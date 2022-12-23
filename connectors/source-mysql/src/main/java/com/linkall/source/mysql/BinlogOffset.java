package com.linkall.source.mysql;

/*
 ["linkall",{"server":"linkall"}] ->
 {"ts_sec":1658130705,"file":"binlog.000010","pos":44602,"gtids":"46ce11bf-9497-11ec-97a5-0e18e1c0f63b:1-75","snapshot":true}
*/
public class BinlogOffset {
  private Integer pos;
  private String file;
  private String gtids;

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
