package com.knowledgebase.platformspring.config;

import org.springframework.boot.context.properties.ConfigurationProperties;
import org.springframework.context.annotation.Configuration;

import lombok.Data;

@Data
@Configuration
@ConfigurationProperties(prefix = "app.openai")
public class OpenAIConfig {
    
    private String apiKey;
    private String baseUrl;
    private String model;
    private String embeddingModel;
    private String timeout;
}

