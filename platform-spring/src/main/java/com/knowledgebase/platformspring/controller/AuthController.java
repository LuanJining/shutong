package com.knowledgebase.platformspring.controller;

import org.springframework.security.core.Authentication;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import com.knowledgebase.platformspring.dto.ApiResponse;
import com.knowledgebase.platformspring.dto.LoginRequest;
import com.knowledgebase.platformspring.dto.LoginResponse;
import com.knowledgebase.platformspring.dto.RegisterRequest;
import com.knowledgebase.platformspring.model.User;
import com.knowledgebase.platformspring.service.AuthService;

import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import reactor.core.publisher.Mono;

@RestController
@RequestMapping("/api/v1/auth")
@RequiredArgsConstructor
public class AuthController {
    
    private final AuthService authService;
    
    @PostMapping("/login")
    public Mono<ApiResponse<LoginResponse>> login(@Valid @RequestBody LoginRequest request) {
        return authService.login(request)
                .map(response -> ApiResponse.success("Login successful", response));
    }
    
    @PostMapping("/register")
    public Mono<ApiResponse<User>> register(@Valid @RequestBody RegisterRequest request) {
        return authService.register(request)
                .map(user -> ApiResponse.success("Registration successful", user));
    }
    
    @PostMapping("/refresh")
    public Mono<ApiResponse<LoginResponse>> refreshToken(@RequestBody RefreshTokenRequest request) {
        return authService.refreshToken(request.refreshToken)
                .map(response -> ApiResponse.success("Token refreshed", response));
    }
    
    @GetMapping("/me")
    public Mono<ApiResponse<User>> getCurrentUser(Authentication authentication) {
        Long userId = (Long) authentication.getPrincipal();
        return authService.getCurrentUser(userId)
                .map(user -> ApiResponse.success(user));
    }
    
    private record RefreshTokenRequest(String refreshToken) {}
}

