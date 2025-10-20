package com.knowledgebase.platformspring.util;

import reactor.core.publisher.Mono;
import reactor.util.context.ContextView;

/**
 * Request ID 工具类 - 用于获取当前请求的 requestId
 */
public class RequestIdUtil {
    
    private static final String REQUEST_ID_KEY = "requestId";
    
    /**
     * 从 Reactor Context 中获取 requestId
     */
    public static Mono<String> getRequestId() {
        return Mono.deferContextual(ctx -> {
            String requestId = ctx.getOrDefault(REQUEST_ID_KEY, "unknown");
            return Mono.just(requestId);
        });
    }
    
    /**
     * 从 ContextView 中获取 requestId
     */
    public static String getRequestIdFromContext(ContextView context) {
        return context.getOrDefault(REQUEST_ID_KEY, "unknown");
    }
}

