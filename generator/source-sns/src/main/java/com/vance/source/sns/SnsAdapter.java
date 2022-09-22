package com.vance.source.sns;

import com.linkall.vance.core.Adapter1;
import io.cloudevents.CloudEvent;
import io.cloudevents.core.builder.CloudEventBuilder;
import java.net.URI;

public class SnsAdapter implements Adapter1<SnsAdaptedSample>{

    public static final CloudEventBuilder template = CloudEventBuilder.v1();

    @Override
    public CloudEvent adapt(SnsAdaptedSample adaptedSample) {
        template.withId("xxxxx");

        URI uri = URI.create("xxxxx");
        template.withSource(uri);

        template.withType("xxxxx")
                .withDataContentType("application/json");

        return template.build();
    }

}