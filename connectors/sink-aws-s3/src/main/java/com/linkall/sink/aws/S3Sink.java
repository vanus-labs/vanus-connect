package com.linkall.sink.aws;

import com.linkall.vance.common.config.ConfigUtil;
import com.linkall.vance.common.json.JsonMapper;
import com.linkall.vance.core.Sink;
import com.linkall.vance.core.http.HttpServer;
import io.vertx.core.json.JsonObject;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import software.amazon.awssdk.regions.Region;
import software.amazon.awssdk.services.s3.S3Client;

import java.io.*;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.text.DecimalFormat;
import java.time.*;
import java.util.Objects;
import java.util.concurrent.Executors;
import java.util.concurrent.ScheduledExecutorService;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicInteger;


public class S3Sink implements Sink {
    private static final Logger LOGGER = LoggerFactory.getLogger(S3Sink.class);
    private static final AtomicInteger eventNum = new AtomicInteger(0);
    private static final DecimalFormat decimalFormat = new DecimalFormat("0000000");
    private static final DecimalFormat dateFormat = new DecimalFormat("00");
    private LocalDateTime fileCreateTime;
    private OffsetDateTime pathCreateTime;
    private static final AtomicInteger fileIdx = new AtomicInteger(1);
    private static final AtomicInteger fileSize = new AtomicInteger(0);
    private String pathName;
    private String timeIntervalUnit;
    private int flushSize = 1000;
    private long scheduledInterval = 30;
    private S3Client s3;

    @Override
    public void start(){
        AwsHelper.checkCredentials();

        //read config
        String strRegion = ConfigUtil.getString("region");
        String bucketName = ConfigUtil.getString("bucketName");

        if(ConfigUtil.getString("flushSize")!=null){
            flushSize = Integer.parseInt(ConfigUtil.getString("flushSize"));
        }
        if(ConfigUtil.getString("scheduledInterval")!=null){
            scheduledInterval = Long.parseLong(ConfigUtil.getString("scheduledInterval"));
        }

        //crate path for file to be upload
        pathCreateTime = getZeroTime(LocalDateTime.now());
        timeIntervalUnit = ConfigUtil.getString("timeInterval");
        if(timeIntervalUnit.equals("HOURLY")){
            pathName = pathCreateTime.getYear()+"/"+dateFormat.format(pathCreateTime.getMonthValue())+"/"
                    +dateFormat.format(pathCreateTime.getDayOfMonth())+"/"+dateFormat.format(pathCreateTime.getHour())+"/";
        }else if(timeIntervalUnit.equals("DAILY")){
            pathName = pathCreateTime.getYear()+"/"+dateFormat.format(pathCreateTime.getMonthValue())+"/"
                    +dateFormat.format(pathCreateTime.getDayOfMonth())+"/";
        }

        Region region = Region.of(strRegion);
        // get S3Client
        s3 = Objects.requireNonNull(
                S3Client.builder().region(region).build());

        HttpServer server = HttpServer.createHttpServer();

        fileCreateTime = LocalDateTime.now();


        //scheduled thread check upload condition
        ScheduledExecutorService threadPool = Executors.newScheduledThreadPool(10);
        threadPool.scheduleAtFixedRate(new Runnable() {
            @Override
            public void run() {
                pathName = checkAndGetPathName();
                long duration = Duration.between(fileCreateTime, LocalDateTime.now()).getSeconds();
                if((fileSize.get() >= flushSize || duration >= scheduledInterval)){
                    uploadFile(bucketName);
                }
            }
        }, 0L, 1L, TimeUnit.SECONDS);

        //write ce into file
        server.ceHandler(event -> {
            int num = eventNum.addAndGet(1);
            LOGGER.info("receive a new event, in total: "+num);

            JsonObject js = JsonMapper.wrapCloudEvent(event);

            File jsonFile = new File("eventing-"+decimalFormat.format(fileIdx.get()));

            try (BufferedWriter bw = new BufferedWriter(new OutputStreamWriter(new FileOutputStream(jsonFile, true)))){
                bw.write(js.toString());
                bw.write("\r\n");
            } catch (IOException e) {
                e.printStackTrace();
            }
            fileSize.getAndAdd(1);
        });
        server.listen();

        Runtime.getRuntime().addShutdownHook(new Thread(()->{
            pathName = checkAndGetPathName();
            uploadFile(bucketName);
            System.out.println("shut down");
        }));
    }

    private OffsetDateTime getZeroTime(LocalDateTime time){
        LocalDateTime dt = LocalDateTime.now(ZoneId.of("Z"));
        Duration duration = Duration.between(time, dt);
        OffsetDateTime time2 = OffsetDateTime.of(time, ZoneOffset.UTC).plus(duration);
        return time2;
    }

    private String checkAndGetPathName(){
        String pathName = this.pathName;
        pathCreateTime = getZeroTime(LocalDateTime.now());
        String newPathNameDaily = pathCreateTime.getYear()+"/"+dateFormat.format(pathCreateTime.getMonthValue())+"/"
                +dateFormat.format(pathCreateTime.getDayOfMonth())+"/";
        if(timeIntervalUnit.equals("HOURLY")){
            if(!pathName.equals(newPathNameDaily+dateFormat.format(pathCreateTime.getHour())+"/")){
                fileIdx.getAndSet(1);
                fileSize.getAndSet(0);
            }
            pathName = newPathNameDaily + dateFormat.format(pathCreateTime.getHour()) + "/";
        }else if(timeIntervalUnit.equals("DAILY")){
            if(!pathName.equals(newPathNameDaily)){
                fileIdx.getAndSet(1);
                fileSize.getAndSet(0);
            }
            pathName = newPathNameDaily;
        }
        return pathName;
    }

    private void uploadFile(String bucketName){
        File uploadFile = new File("eventing-"+decimalFormat.format(fileIdx.get()));
        if(null != uploadFile && uploadFile.length() != 0){
            int uploadFileIdx = fileIdx.get();
            fileSize.getAndSet(0);
            fileIdx.getAndAdd(1);
            boolean putOk = S3Util.putS3Object(s3,bucketName,
                    pathName+"eventing-"+decimalFormat.format(uploadFileIdx), uploadFile);
            try {
                Files.deleteIfExists(Paths.get("eventing-"+decimalFormat.format(uploadFileIdx) ));
            } catch (IOException e) {
                e.printStackTrace();
            }
            if(putOk){
                LOGGER.info("[upload file <" + "eventing-"+decimalFormat.format(uploadFileIdx) + "> completed");
            }else{
                LOGGER.info("[upload file <" + "eventing-"+decimalFormat.format(uploadFileIdx) + "> failed");
            }
        }else{
            LOGGER.info("invalid data format, upload failed");
        }
        fileCreateTime = LocalDateTime.now();
    }

}