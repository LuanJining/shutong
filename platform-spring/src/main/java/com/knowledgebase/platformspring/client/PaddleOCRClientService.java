package com.knowledgebase.platformspring.client;

import java.time.Duration;
import java.util.Base64;
import java.util.HashMap;
import java.util.Map;

import org.springframework.http.MediaType;
import org.springframework.stereotype.Service;
import org.springframework.web.reactive.function.client.WebClient;

import com.knowledgebase.platformspring.config.PaddleOCRConfig;

import lombok.Data;
import lombok.extern.slf4j.Slf4j;
import reactor.core.publisher.Mono;

@Slf4j
@Service
public class PaddleOCRClientService {
    
    private final WebClient webClient;
    private final PaddleOCRConfig ocrConfig;
    
    public PaddleOCRClientService(PaddleOCRConfig ocrConfig) {
        this.ocrConfig = ocrConfig;
        
        if (ocrConfig.getEnabled() && ocrConfig.getUrl() != null) {
            this.webClient = WebClient.builder()
                    .baseUrl(ocrConfig.getUrl())
                    .build();
        } else {
            this.webClient = null;
        }
    }
    
    public Mono<String> recognize(String fileName, byte[] data) {
        if (webClient == null || !ocrConfig.getEnabled()) {
            return Mono.error(new RuntimeException("PaddleOCR client is not configured"));
        }
        
        Map<String, String> requestBody = new HashMap<>();
        requestBody.put("file_name", fileName);
        requestBody.put("content_base64", Base64.getEncoder().encodeToString(data));
        requestBody.put("language", "ch");
        
        return webClient.post()
                .uri("/v1/ocr")
                .contentType(MediaType.APPLICATION_JSON)
                .bodyValue(requestBody)
                .retrieve()
                .bodyToMono(OCRResponse.class)
                .map(OCRResponse::getText)
                .timeout(Duration.ofMinutes(10))
                .doOnError(e -> log.error("PaddleOCR request failed for file: {}", fileName, e));
    }
    
    @Data
    private static class OCRResponse {
        private String text;
        private String error;
    }
}

