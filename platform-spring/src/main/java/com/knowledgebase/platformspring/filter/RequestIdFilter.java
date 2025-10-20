package com.knowledgebase.platformspring.filter;

import java.util.UUID;

import org.springframework.core.annotation.Order;
import org.springframework.stereotype.Component;
import org.springframework.web.server.ServerWebExchange;
import org.springframework.web.server.WebFilter;
import org.springframework.web.server.WebFilterChain;

import lombok.extern.slf4j.Slf4j;
import reactor.core.publisher.Mono;
import reactor.util.context.Context;

/**
 * Request ID 过滤器 - 为每个请求生成唯一ID，用于日志追踪
 */
@Slf4j
@Component
@Order(-100) // 确保在其他过滤器之前执行
public class RequestIdFilter implements WebFilter {
    
    public static final String REQUEST_ID_HEADER = "X-Request-Id";
    public static final String REQUEST_ID_KEY = "requestId";
    
    @Override
    public Mono<Void> filter(ServerWebExchange exchange, WebFilterChain chain) {
        // 从请求头获取或生成新的 requestId
        String requestId = exchange.getRequest().getHeaders().getFirst(REQUEST_ID_HEADER);
        if (requestId == null || requestId.trim().isEmpty()) {
            requestId = generateRequestId();
        }
        
        // 添加到响应头
        exchange.getResponse().getHeaders().add(REQUEST_ID_HEADER, requestId);
        
        // 记录请求开始
        String method = exchange.getRequest().getMethod().name();
        String path = exchange.getRequest().getPath().value();
        log.debug("[{}] Request started: {} {}", requestId, method, path);
        
        final String finalRequestId = requestId;
        
        // 将 requestId 存入 Reactor Context，供后续处理使用
        return chain.filter(exchange)
                .contextWrite(Context.of(REQUEST_ID_KEY, finalRequestId))
                .doOnSuccess(v -> log.debug("[{}] Request completed: {} {}", finalRequestId, method, path))
                .doOnError(e -> log.error("[{}] Request failed: {} {} - {}", finalRequestId, method, path, e.getMessage()));
    }
    
    private String generateRequestId() {
        return UUID.randomUUID().toString().replace("-", "").substring(0, 16);
    }
}

