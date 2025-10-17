package com.knowledgebase.platformspring.client;

import java.io.InputStream;

import org.springframework.stereotype.Service;

import com.knowledgebase.platformspring.config.MinioConfig;

import io.minio.BucketExistsArgs;
import io.minio.GetObjectArgs;
import io.minio.MakeBucketArgs;
import io.minio.MinioClient;
import io.minio.PutObjectArgs;
import io.minio.RemoveObjectArgs;
import lombok.extern.slf4j.Slf4j;
import reactor.core.publisher.Mono;

@Slf4j
@Service
public class MinioClientService {
    
    private final MinioClient minioClient;
    private final MinioConfig minioConfig;
    
    public MinioClientService(MinioConfig minioConfig) {
        this.minioConfig = minioConfig;
        this.minioClient = MinioClient.builder()
                .endpoint(minioConfig.getEndpoint())
                .credentials(minioConfig.getAccessKey(), minioConfig.getSecretKey())
                .build();
        
        ensureBucketExists().subscribe();
    }
    
    private Mono<Void> ensureBucketExists() {
        return Mono.fromCallable(() -> {
            boolean found = minioClient.bucketExists(
                BucketExistsArgs.builder()
                    .bucket(minioConfig.getBucketName())
                    .build()
            );
            
            if (!found) {
                minioClient.makeBucket(
                    MakeBucketArgs.builder()
                        .bucket(minioConfig.getBucketName())
                        .build()
                );
                log.info("Created MinIO bucket: {}", minioConfig.getBucketName());
            }
            return null;
        }).then();
    }
    
    public Mono<String> uploadFile(String objectName, InputStream inputStream, long size, String contentType) {
        return Mono.fromCallable(() -> {
            minioClient.putObject(
                PutObjectArgs.builder()
                    .bucket(minioConfig.getBucketName())
                    .object(objectName)
                    .stream(inputStream, size, -1)
                    .contentType(contentType)
                    .build()
            );
            return objectName;
        }).doOnSuccess(name -> log.debug("Uploaded file to MinIO: {}", name))
          .doOnError(e -> log.error("Failed to upload file to MinIO: {}", objectName, e));
    }
    
    public Mono<InputStream> downloadFile(String objectName) {
        return Mono.fromCallable(() -> {
            return (InputStream) minioClient.getObject(
                GetObjectArgs.builder()
                    .bucket(minioConfig.getBucketName())
                    .object(objectName)
                    .build()
            );
        }).doOnError(e -> log.error("Failed to download file from MinIO: {}", objectName, e));
    }
    
    public Mono<Void> deleteFile(String objectName) {
        return Mono.fromRunnable(() -> {
            try {
                minioClient.removeObject(
                    RemoveObjectArgs.builder()
                        .bucket(minioConfig.getBucketName())
                        .object(objectName)
                        .build()
                );
            } catch (Exception e) {
                throw new RuntimeException("Failed to delete file from MinIO", e);
            }
        }).then();
    }
    
    public String getBucketName() {
        return minioConfig.getBucketName();
    }
}

