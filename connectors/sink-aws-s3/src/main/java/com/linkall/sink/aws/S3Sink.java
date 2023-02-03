package com.linkall.sink.aws;

import com.linkall.cdk.config.Config;
import com.linkall.cdk.connector.Result;
import com.linkall.cdk.connector.Sink;
import com.linkall.cdk.util.EventUtil;
import io.cloudevents.CloudEvent;
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
    private S3Config.TimeInterval timeInterval;
    private int flushSize;
    private long scheduledInterval;
    private S3Client s3;

    private S3Config config;
    private ScheduledExecutorService executorService;


    public S3Sink() {
        executorService = Executors.newSingleThreadScheduledExecutor();
    }

    public void start() {
        AwsHelper.checkCredentials(config.getSecretConfig().getAccessKeyID(), config.getSecretConfig().getSecretAccessKey());

        //read config
        String strRegion = config.getRegion();
        String bucketName = config.getBucket();

        if (config.getFlushSize()!=null && config.getFlushSize() > 0) {
            flushSize = config.getFlushSize();
        } else {
            flushSize = 1000;
        }
        if (config.getScheduledInterval()!=null && config.getScheduledInterval() > 0) {
            scheduledInterval = config.getScheduledInterval();
        } else {
            scheduledInterval = 60;
        }

        //crate path for file to be upload
        pathCreateTime = getZeroTime(LocalDateTime.now());
        timeInterval = config.getTimeInterval();
        if (timeInterval==null) {
            timeInterval = S3Config.TimeInterval.HOURLY;
        }
        if (timeInterval==S3Config.TimeInterval.HOURLY) {
            pathName = pathCreateTime.getYear() + "/" + dateFormat.format(pathCreateTime.getMonthValue()) + "/"
                    + dateFormat.format(pathCreateTime.getDayOfMonth()) + "/" + dateFormat.format(pathCreateTime.getHour()) + "/";
        } else if (timeInterval==S3Config.TimeInterval.DAILY) {
            pathName = pathCreateTime.getYear() + "/" + dateFormat.format(pathCreateTime.getMonthValue()) + "/"
                    + dateFormat.format(pathCreateTime.getDayOfMonth()) + "/";
        }

        Region region = Region.of(strRegion);
        // get S3Client
        s3 = Objects.requireNonNull(S3Client.builder().region(region).build());

        fileCreateTime = LocalDateTime.now();


        //scheduled thread check upload condition
        executorService.scheduleAtFixedRate(() -> {
            pathName = checkAndGetPathName();
            long duration = Duration.between(fileCreateTime, LocalDateTime.now()).getSeconds();
            if ((fileSize.get() >= flushSize || duration >= scheduledInterval)) {
                uploadFile(bucketName);
            }
        }, 0L, 1L, TimeUnit.SECONDS);
    }

    private OffsetDateTime getZeroTime(LocalDateTime time) {
        LocalDateTime dt = LocalDateTime.now(ZoneId.of("Z"));
        Duration duration = Duration.between(time, dt);
        OffsetDateTime time2 = OffsetDateTime.of(time, ZoneOffset.UTC).plus(duration);
        return time2;
    }

    private String checkAndGetPathName() {
        String pathName = this.pathName;
        pathCreateTime = getZeroTime(LocalDateTime.now());
        String newPathNameDaily = pathCreateTime.getYear() + "/" + dateFormat.format(pathCreateTime.getMonthValue()) + "/"
                + dateFormat.format(pathCreateTime.getDayOfMonth()) + "/";
        if (timeInterval==S3Config.TimeInterval.HOURLY) {
            if (!pathName.equals(newPathNameDaily + dateFormat.format(pathCreateTime.getHour()) + "/")) {
                fileIdx.getAndSet(1);
                fileSize.getAndSet(0);
            }
            pathName = newPathNameDaily + dateFormat.format(pathCreateTime.getHour()) + "/";
        } else if (timeInterval==S3Config.TimeInterval.DAILY) {
            if (!pathName.equals(newPathNameDaily)) {
                fileIdx.getAndSet(1);
                fileSize.getAndSet(0);
            }
            pathName = newPathNameDaily;
        }
        return pathName;
    }

    private void uploadFile(String bucketName) {
        File uploadFile = new File("eventing-" + decimalFormat.format(fileIdx.get()));
        if (uploadFile.exists() && uploadFile.length()!=0) {
            int uploadFileIdx = fileIdx.get();
            fileSize.getAndSet(0);
            fileIdx.getAndAdd(1);
            boolean putOk = S3Util.putS3Object(s3, bucketName,
                    pathName + "eventing-" + decimalFormat.format(uploadFileIdx), uploadFile);
            try {
                Files.deleteIfExists(Paths.get("eventing-" + decimalFormat.format(uploadFileIdx)));
            } catch (IOException e) {
                e.printStackTrace();
            }
            if (putOk) {
                LOGGER.info("[upload file <" + "eventing-" + decimalFormat.format(uploadFileIdx) + "> completed");
            } else {
                LOGGER.info("[upload file <" + "eventing-" + decimalFormat.format(uploadFileIdx) + "> failed");
            }
        } else {
            LOGGER.info("no event ignore upload file,{}", uploadFile.getName());
        }
        fileCreateTime = LocalDateTime.now();
    }

    @Override
    public Result Arrived(CloudEvent... events) {
        for (CloudEvent event : events) {
            int num = eventNum.addAndGet(1);
            LOGGER.info("receive a new event, in total: " + num);
            File jsonFile = new File("eventing-" + decimalFormat.format(fileIdx.get()));
            try (BufferedWriter bw = new BufferedWriter(new OutputStreamWriter(new FileOutputStream(jsonFile, true)))) {
                bw.write(EventUtil.eventToJson(event));
                bw.write("\r\n");
            } catch (IOException e) {
                LOGGER.error("write event to file error,{},{}", jsonFile.getAbsolutePath(), event.getId(), e);
                throw new RuntimeException(e);
            }
            fileSize.getAndAdd(1);
        }
        return Result.SUCCESS;
    }

    @Override
    public Class<? extends Config> configClass() {
        return S3Config.class;
    }

    @Override
    public void initialize(Config config) throws Exception {
        this.config = (S3Config) config;
        start();
    }

    @Override
    public String name() {
        return "AmazonS3Sink";
    }

    @Override
    public void destroy() throws Exception {
        executorService.shutdown();
        pathName = checkAndGetPathName();
        uploadFile(config.getBucket());
        s3.close();
    }
}