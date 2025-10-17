package com.knowledgebase.platformspring.service;

import java.time.LocalDateTime;

import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Service;

import com.knowledgebase.platformspring.dto.LoginRequest;
import com.knowledgebase.platformspring.dto.LoginResponse;
import com.knowledgebase.platformspring.dto.RegisterRequest;
import com.knowledgebase.platformspring.exception.BusinessException;
import com.knowledgebase.platformspring.model.User;
import com.knowledgebase.platformspring.repository.UserRepository;
import com.knowledgebase.platformspring.security.JwtUtil;

import lombok.RequiredArgsConstructor;
import reactor.core.publisher.Mono;

@Service
@RequiredArgsConstructor
public class AuthService {
    
    private final UserRepository userRepository;
    private final PasswordEncoder passwordEncoder;
    private final JwtUtil jwtUtil;
    
    public Mono<LoginResponse> login(LoginRequest request) {
        // 支持用户名、手机号、邮箱登录
        return userRepository.findByUsername(request.getLogin())
                .switchIfEmpty(userRepository.findByPhone(request.getLogin()))
                .switchIfEmpty(userRepository.findByEmail(request.getLogin()))
                .switchIfEmpty(Mono.error(BusinessException.unauthorized("用户名或密码错误")))
                .flatMap(user -> {
                    if (!passwordEncoder.matches(request.getPassword(), user.getPassword())) {
                        return Mono.error(BusinessException.unauthorized("用户名或密码错误"));
                    }
                    
                    if (user.getStatus() == 0) {
                        return Mono.error(BusinessException.forbidden("用户已被禁用"));
                    }
                    
                    user.setLastLogin(LocalDateTime.now());
                    return userRepository.save(user)
                            .map(savedUser -> {
                                String accessToken = jwtUtil.generateToken(savedUser.getId(), savedUser.getUsername());
                                String refreshToken = jwtUtil.generateRefreshToken(savedUser.getId());
                                
                                LocalDateTime accessExpires = LocalDateTime.now()
                                        .plusSeconds(jwtUtil.getAccessTokenExpiration() / 1000);
                                LocalDateTime refreshExpires = LocalDateTime.now()
                                        .plusSeconds(jwtUtil.getRefreshTokenExpiration() / 1000);
                                
                                // Clear password before returning
                                savedUser.setPassword(null);
                                
                                return LoginResponse.builder()
                                        .accessToken(accessToken)
                                        .refreshToken(refreshToken)
                                        .user(savedUser)
                                        .accessTokenExpiresAt(accessExpires)
                                        .refreshTokenExpiresAt(refreshExpires)
                                        .build();
                            });
                });
    }
    
    public Mono<User> register(RegisterRequest request) {
        return userRepository.existsByUsername(request.getUsername())
                .flatMap(exists -> {
                    if (exists) {
                        return Mono.error(new BusinessException("Username already exists"));
                    }
                    return userRepository.existsByPhone(request.getPhone());
                })
                .flatMap(exists -> {
                    if (exists) {
                        return Mono.error(new BusinessException("Phone number already exists"));
                    }
                    
                    User user = User.builder()
                            .username(request.getUsername())
                            .phone(request.getPhone())
                            .email(request.getEmail())
                            .password(passwordEncoder.encode(request.getPassword()))
                            .nickname(request.getNickname())
                            .department(request.getDepartment())
                            .company(request.getCompany())
                            .status(1)
                            .createdAt(LocalDateTime.now())
                            .updatedAt(LocalDateTime.now())
                            .build();
                    
                    return userRepository.save(user)
                            .doOnSuccess(savedUser -> savedUser.setPassword(null));
                });
    }
    
    public Mono<LoginResponse> refreshToken(String refreshToken) {
        if (!jwtUtil.validateToken(refreshToken)) {
            return Mono.error(BusinessException.unauthorized("Invalid refresh token"));
        }
        
        Long userId = jwtUtil.getUserIdFromToken(refreshToken);
        
        return userRepository.findById(userId)
                .switchIfEmpty(Mono.error(BusinessException.notFound("User not found")))
                .map(user -> {
                    String newAccessToken = jwtUtil.generateToken(user.getId(), user.getUsername());
                    String newRefreshToken = jwtUtil.generateRefreshToken(user.getId());
                    
                    user.setPassword(null);
                    
                    return LoginResponse.builder()
                            .accessToken(newAccessToken)
                            .refreshToken(newRefreshToken)
                            .user(user)
                            .build();
                });
    }
    
    public Mono<User> getCurrentUser(Long userId) {
        return userRepository.findById(userId)
                .switchIfEmpty(Mono.error(BusinessException.notFound("User not found")))
                .doOnSuccess(user -> user.setPassword(null));
    }
    
    public Mono<Void> changePassword(Long userId, String oldPassword, String newPassword) {
        return userRepository.findById(userId)
                .switchIfEmpty(Mono.error(BusinessException.notFound("User not found")))
                .flatMap(user -> {
                    if (!passwordEncoder.matches(oldPassword, user.getPassword())) {
                        return Mono.error(new BusinessException("原密码错误"));
                    }
                    
                    user.setPassword(passwordEncoder.encode(newPassword));
                    user.setUpdatedAt(LocalDateTime.now());
                    
                    return userRepository.save(user).then();
                });
    }
}

