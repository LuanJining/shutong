package com.knowledgebase.platformspring.config;

import org.springframework.boot.context.properties.ConfigurationProperties;
import org.springframework.context.annotation.Configuration;

import lombok.Data;

@Data
@Configuration
@ConfigurationProperties(prefix = "app.minio")
public class MinioConfig {
    
    private String endpoint;
    private String accessKey;
    private String secretKey;
    private String bucketName;
    private Boolean secure = false;
}

