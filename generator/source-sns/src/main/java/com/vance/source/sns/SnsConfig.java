package com.vance.source.sns;

import com.linkall.vance.common.annotation.Tag;

public class SnsConfig{

    @Tag(key = "v_target")
    private String vTarget;

    public void setVanceTarget(String vTarget){
        this.vTarget = vTarget;
    }

    public String getVanceTarget(){
        return vTarget;
    }

}