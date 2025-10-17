package com.knowledgebase.platformspring.config;

import org.springframework.boot.context.properties.ConfigurationProperties;
import org.springframework.context.annotation.Configuration;

import lombok.Data;

@Data
@Configuration
@ConfigurationProperties(prefix = "spring.security.jwt")
public class JwtConfig {
    
    private String secret;
    private Long expiration; // milliseconds
    private Long refreshExpiration; // milliseconds
}

