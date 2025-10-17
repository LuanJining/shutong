package com.knowledgebase.platformspring.dto;

import jakarta.validation.constraints.NotBlank;
import lombok.Data;

@Data
public class LoginRequest {
    
    @NotBlank(message = "Login is required")
    private String login; // 支持用户名、手机号、邮箱登录
    
    @NotBlank(message = "Password is required")
    private String password;
}

