package com.knowledgebase.platformspring.config;

import org.springframework.boot.context.properties.ConfigurationProperties;
import org.springframework.context.annotation.Configuration;

import lombok.Data;

@Data
@Configuration
@ConfigurationProperties(prefix = "app.qdrant")
public class QdrantConfig {
    
    private String host;
    private Integer port;
    private String apiKey;
    private String collectionName;
    private Integer vectorSize;
    private Boolean useHttps = false;
}

