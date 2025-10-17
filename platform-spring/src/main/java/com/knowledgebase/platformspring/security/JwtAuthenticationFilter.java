package com.knowledgebase.platformspring.security;

import org.springframework.http.HttpHeaders;
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.core.context.ReactiveSecurityContextHolder;
import org.springframework.stereotype.Component;
import org.springframework.web.server.ServerWebExchange;
import org.springframework.web.server.WebFilter;
import org.springframework.web.server.WebFilterChain;

import lombok.RequiredArgsConstructor;
import reactor.core.publisher.Mono;

@Component
@RequiredArgsConstructor
public class JwtAuthenticationFilter implements WebFilter {
    
    private final JwtUtil jwtUtil;
    
    @Override
    public Mono<Void> filter(ServerWebExchange exchange, WebFilterChain chain) {
        String path = exchange.getRequest().getPath().value();
        
        // Skip JWT validation for public endpoints
        if (path.startsWith("/api/v1/auth/login") || 
            path.startsWith("/api/v1/auth/register") ||
            path.startsWith("/api/v1/auth/refresh") ||
            path.startsWith("/actuator")) {
            return chain.filter(exchange);
        }
        
        String authHeader = exchange.getRequest().getHeaders().getFirst(HttpHeaders.AUTHORIZATION);
        
        if (authHeader != null && authHeader.startsWith("Bearer ")) {
            String token = authHeader.substring(7);
            
            if (jwtUtil.validateToken(token)) {
                Long userId = jwtUtil.getUserIdFromToken(token);
                
                UsernamePasswordAuthenticationToken authentication = 
                    new UsernamePasswordAuthenticationToken(userId, null, null);
                
                return chain.filter(exchange)
                        .contextWrite(ReactiveSecurityContextHolder.withAuthentication(authentication));
            }
        }
        
        return chain.filter(exchange);
    }
}

