package com.knowledgebase.platformspring;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.boot.context.properties.EnableConfigurationProperties;
import org.springframework.data.r2dbc.repository.config.EnableR2dbcRepositories;

@SpringBootApplication
@EnableR2dbcRepositories
@EnableConfigurationProperties
public class PlatformSpringApplication {

    public static void main(String[] args) {
        SpringApplication.run(PlatformSpringApplication.class, args);
    }

}
