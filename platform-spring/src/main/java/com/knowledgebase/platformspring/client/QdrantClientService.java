package com.knowledgebase.platformspring.client;

import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.UUID;

import org.springframework.http.MediaType;
import org.springframework.stereotype.Service;
import org.springframework.web.reactive.function.client.WebClient;

import com.knowledgebase.platformspring.config.QdrantConfig;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import reactor.core.publisher.Flux;
import reactor.core.publisher.Mono;

@Slf4j
@Service
public class QdrantClientService {
    
    private final WebClient webClient;
    private final QdrantConfig qdrantConfig;
    
    public QdrantClientService(QdrantConfig qdrantConfig) {
        this.qdrantConfig = qdrantConfig;
        String baseUrl = (qdrantConfig.getUseHttps() ? "https://" : "http://") + 
                         qdrantConfig.getHost() + ":" + qdrantConfig.getPort();
        
        this.webClient = WebClient.builder()
                .baseUrl(baseUrl)
                .defaultHeaders(headers -> {
                    if (qdrantConfig.getApiKey() != null && !qdrantConfig.getApiKey().isEmpty()) {
                        headers.set("api-key", qdrantConfig.getApiKey());
                    }
                })
                .build();
        
        ensureCollectionExists().subscribe();
    }
    
    public Mono<Void> ensureCollectionExists() {
        Map<String, Object> collectionConfig = new HashMap<>();
        Map<String, Object> vectors = new HashMap<>();
        vectors.put("size", qdrantConfig.getVectorSize());
        vectors.put("distance", "Cosine");
        collectionConfig.put("vectors", vectors);
        
        return webClient.put()
                .uri("/collections/" + qdrantConfig.getCollectionName())
                .contentType(MediaType.APPLICATION_JSON)
                .bodyValue(collectionConfig)
                .retrieve()
                .bodyToMono(Void.class)
                .doOnSuccess(v -> log.info("Ensured Qdrant collection exists: {}", qdrantConfig.getCollectionName()))
                .onErrorResume(e -> {
                    log.warn("Collection may already exist: {}", e.getMessage());
                    return Mono.empty();
                });
    }
    
    public Mono<Void> upsertPoints(List<QdrantPoint> points) {
        Map<String, Object> request = new HashMap<>();
        request.put("points", points);
        
        return webClient.put()
                .uri("/collections/" + qdrantConfig.getCollectionName() + "/points")
                .contentType(MediaType.APPLICATION_JSON)
                .bodyValue(request)
                .retrieve()
                .bodyToMono(Void.class)
                .doOnSuccess(v -> log.debug("Upserted {} points to Qdrant", points.size()));
    }
    
    public Flux<QdrantSearchResult> searchPoints(List<Double> vector, int limit) {
        Map<String, Object> request = new HashMap<>();
        request.put("vector", vector);
        request.put("top", limit);
        request.put("with_payload", true);
        request.put("with_vector", false);
        
        return webClient.post()
                .uri("/collections/" + qdrantConfig.getCollectionName() + "/points/search")
                .contentType(MediaType.APPLICATION_JSON)
                .bodyValue(request)
                .retrieve()
                .bodyToMono(QdrantSearchResponse.class)
                .flatMapMany(response -> Flux.fromIterable(response.getResult()));
    }
    
    @Data
    @Builder
    @NoArgsConstructor
    @AllArgsConstructor
    public static class QdrantPoint {
        private String id;
        private List<Double> vector;
        private Map<String, Object> payload;
        
        public static QdrantPoint create(List<Double> vector, Map<String, Object> payload) {
            return QdrantPoint.builder()
                    .id(UUID.randomUUID().toString())
                    .vector(vector)
                    .payload(payload)
                    .build();
        }
    }
    
    @Data
    public static class QdrantSearchResult {
        private String id;
        private Double score;
        private Map<String, Object> payload;
    }
    
    @Data
    private static class QdrantSearchResponse {
        private List<QdrantSearchResult> result;
        private String status;
    }
}

