package com.vance.source.sns;

import com.linkall.vance.core.Source;
import com.linkall.vance.core.Adapter;

public class SnsSource implements Source {

    @Override
    public void start(){

    }

    @Override
    public Adapter getAdapter() {
       return new SnsAdapter();
    }
}