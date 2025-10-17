package com.knowledgebase.platformspring.client;

import java.time.Duration;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

import org.springframework.http.MediaType;
import org.springframework.stereotype.Service;
import org.springframework.web.reactive.function.client.WebClient;

import com.knowledgebase.platformspring.config.OpenAIConfig;

import lombok.Data;
import lombok.extern.slf4j.Slf4j;
import reactor.core.publisher.Mono;

@Slf4j
@Service
public class OpenAIClientService {
    
    private final WebClient webClient;
    private final OpenAIConfig openAIConfig;
    
    public OpenAIClientService(OpenAIConfig openAIConfig) {
        this.openAIConfig = openAIConfig;
        
        this.webClient = WebClient.builder()
                .baseUrl(openAIConfig.getBaseUrl())
                .defaultHeader("Authorization", "Bearer " + openAIConfig.getApiKey())
                .defaultHeader("Content-Type", "application/json")
                .build();
    }
    
    public Mono<String> chat(String question, List<String> context) {
        StringBuilder systemPrompt = new StringBuilder();
        systemPrompt.append("你是一个智能助手，可以基于提供的知识库内容回答问题。请根据以下文件内容来回答用户的问题：\n\n");
        
        for (int i = 0; i < context.size(); i++) {
            systemPrompt.append("文件 ").append(i + 1).append(" 内容：\n");
            systemPrompt.append(context.get(i)).append("\n\n");
        }
        
        systemPrompt.append("请基于以上知识库内容回答用户的问题。如果问题与文件内容无关，请说明无法从提供的文件中找到相关信息。");
        
        List<Map<String, String>> messages = new ArrayList<>();
        Map<String, String> systemMessage = new HashMap<>();
        systemMessage.put("role", "system");
        systemMessage.put("content", systemPrompt.toString());
        messages.add(systemMessage);
        
        Map<String, String> userMessage = new HashMap<>();
        userMessage.put("role", "user");
        userMessage.put("content", question);
        messages.add(userMessage);
        
        Map<String, Object> requestBody = new HashMap<>();
        requestBody.put("model", openAIConfig.getModel());
        requestBody.put("messages", messages);
        requestBody.put("temperature", 0.7);
        requestBody.put("max_tokens", 2000);
        
        return webClient.post()
                .uri("/chat/completions")
                .contentType(MediaType.APPLICATION_JSON)
                .bodyValue(requestBody)
                .retrieve()
                .bodyToMono(ChatCompletionResponse.class)
                .map(response -> {
                    if (response.getChoices() != null && !response.getChoices().isEmpty()) {
                        return response.getChoices().get(0).getMessage().getContent();
                    }
                    return "";
                })
                .timeout(Duration.ofSeconds(60))
                .doOnError(e -> log.error("OpenAI chat request failed", e));
    }
    
    public Mono<List<Double>> createEmbedding(String text) {
        Map<String, Object> requestBody = new HashMap<>();
        requestBody.put("model", openAIConfig.getEmbeddingModel());
        requestBody.put("input", text);
        
        return webClient.post()
                .uri("/embeddings")
                .contentType(MediaType.APPLICATION_JSON)
                .bodyValue(requestBody)
                .retrieve()
                .bodyToMono(EmbeddingResponse.class)
                .map(response -> {
                    if (response.getData() != null && !response.getData().isEmpty()) {
                        List<Double> embedding = response.getData().get(0).getEmbedding();
                        return embedding != null ? embedding : new ArrayList<Double>();
                    }
                    return new ArrayList<Double>();
                })
                .timeout(Duration.ofSeconds(30))
                .doOnError(e -> log.error("OpenAI embedding request failed", e));
    }
    
    @Data
    private static class ChatCompletionResponse {
        private List<Choice> choices;
    }
    
    @Data
    private static class Choice {
        private Message message;
    }
    
    @Data
    private static class Message {
        private String role;
        private String content;
    }
    
    @Data
    private static class EmbeddingResponse {
        private List<EmbeddingData> data;
    }
    
    @Data
    private static class EmbeddingData {
        private List<Double> embedding;
    }
}

