package io.cloudevents.examples.vertx;

import org.apache.commons.codec.binary.Hex;
import org.apache.commons.codec.digest.DigestUtils;
import org.junit.Test;

import java.time.*;
import java.time.temporal.ChronoUnit;
import java.time.temporal.TemporalUnit;

public class createCodeTest {
    @Test
    public void test1(){
        String aaa = "ab56b4d92b40713acc5af89985d4b786";
        String code_sha256 = DigestUtils.sha256Hex("abcde");
        byte[] code_bytes = DigestUtils.sha256("abcde");
        String hex1 = Hex.encodeHexString(code_bytes);
        System.out.println(code_bytes.toString());
        System.out.println(hex1);
    }
    @Test
    public void test2(){
        String time = "2022-07-14T06:28:58Z";
        time = time.substring(0, time.length() - 1);
        LocalDateTime time1 = LocalDateTime.parse(time);
        LocalDateTime dt = LocalDateTime.now();
        ZoneId zone = ZoneId.systemDefault();
        ZonedDateTime zdt = dt.atZone(zone);
        ZoneOffset offset = zdt.getOffset();

        OffsetDateTime time2 = OffsetDateTime.of(time1, offset).plusHours(8);
        System.out.println(time2);
        System.out.println(OffsetDateTime.now());
        System.out.println(ZonedDateTime.now());
    }
    @Test
    public void test3(){
        long timestamp = System.currentTimeMillis()/1000;
        OffsetDateTime offsetDateTime = Instant.ofEpochSecond(timestamp).plus(8, ChronoUnit.HOURS).atOffset(ZoneOffset.ofHours(8));
        System.out.println(offsetDateTime);
    }
    @Test
    public void test4(){
        LocalDateTime time1 = LocalDateTime.now();
        LocalDateTime dt = LocalDateTime.now(ZoneId.of("Z"));
        Duration duration = Duration.between(time1, dt);
        OffsetDateTime time2 = OffsetDateTime.of(time1, ZoneOffset.UTC).plus(duration);
        System.out.println(time2);
    }
}
