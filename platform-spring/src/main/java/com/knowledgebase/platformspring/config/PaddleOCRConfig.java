package com.knowledgebase.platformspring.config;

import org.springframework.boot.context.properties.ConfigurationProperties;
import org.springframework.context.annotation.Configuration;

import lombok.Data;

@Data
@Configuration
@ConfigurationProperties(prefix = "app.paddle-ocr")
public class PaddleOCRConfig {
    
    private String url;
    private Boolean enabled = true;
}

