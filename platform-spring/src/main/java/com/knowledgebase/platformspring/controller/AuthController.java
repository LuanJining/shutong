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

import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.security.SecurityRequirement;
import io.swagger.v3.oas.annotations.tags.Tag;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import reactor.core.publisher.Mono;

@Tag(name = "认证管理", description = "用户认证相关接口")
@RestController
@RequestMapping("/api/v1/auth")
@RequiredArgsConstructor
public class AuthController {
    
    private final AuthService authService;
    
    @Operation(summary = "用户登录", description = "使用用户名和密码登录系统，返回访问令牌")
    @PostMapping("/login")
    public Mono<ApiResponse<LoginResponse>> login(@Valid @RequestBody LoginRequest request) {
        return authService.login(request)
                .map(response -> ApiResponse.success("Login successful", response));
    }
    
    @Operation(summary = "用户注册", description = "注册新用户账号")
    @PostMapping("/register")
    public Mono<ApiResponse<User>> register(@Valid @RequestBody RegisterRequest request) {
        return authService.register(request)
                .map(user -> ApiResponse.success("Registration successful", user));
    }
    
    @Operation(summary = "刷新Token", description = "使用 Refresh Token 获取新的 Access Token")
    @PostMapping("/refresh")
    public Mono<ApiResponse<LoginResponse>> refreshToken(@RequestBody RefreshTokenRequest request) {
        return authService.refreshToken(request.refreshToken)
                .map(response -> ApiResponse.success("Token refreshed", response));
    }
    
    @Operation(summary = "获取当前用户", description = "获取当前登录用户的详细信息",
               security = @SecurityRequirement(name = "bearerAuth"))
    @GetMapping("/me")
    public Mono<ApiResponse<User>> getCurrentUser(Authentication authentication) {
        Long userId = (Long) authentication.getPrincipal();
        return authService.getCurrentUser(userId)
                .map(ApiResponse::success);
    }
    
    @Operation(summary = "修改密码", description = "修改当前用户的登录密码",
               security = @SecurityRequirement(name = "bearerAuth"))
    @io.swagger.v3.oas.annotations.parameters.RequestBody(
        description = "修改密码请求",
        required = true
    )
    @org.springframework.web.bind.annotation.PatchMapping("/change-password")
    public Mono<ApiResponse<Void>> changePassword(
            @RequestBody com.knowledgebase.platformspring.dto.ChangePasswordRequest request,
            Authentication authentication) {
        Long userId = (Long) authentication.getPrincipal();
        return authService.changePassword(userId, request.getOldPassword(), request.getNewPassword())
                .then(Mono.just(ApiResponse.<Void>success("密码修改成功", null)));
    }
    
    private record RefreshTokenRequest(String refreshToken) {}
}

