package com.knowledgebase.platformspring.controller;

import org.springframework.web.bind.annotation.DeleteMapping;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.PutMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

import com.knowledgebase.platformspring.dto.ApiResponse;
import com.knowledgebase.platformspring.dto.PaginationResponse;
import com.knowledgebase.platformspring.dto.RegisterRequest;
import com.knowledgebase.platformspring.model.User;
import com.knowledgebase.platformspring.service.AuthService;
import com.knowledgebase.platformspring.service.UserService;

import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.security.SecurityRequirement;
import io.swagger.v3.oas.annotations.tags.Tag;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import reactor.core.publisher.Mono;

@Tag(name = "用户管理", description = "用户管理相关接口")
@RestController
@RequestMapping("/api/v1/users")
@RequiredArgsConstructor
@SecurityRequirement(name = "bearerAuth")
public class UserController {
    
    private final UserService userService;
    private final AuthService authService;
    
    @Operation(summary = "获取用户列表", description = "分页获取用户列表")
    @GetMapping
    public Mono<ApiResponse<PaginationResponse<java.util.List<User>>>> getUsers(
            @RequestParam(defaultValue = "1") Integer page,
            @RequestParam(defaultValue = "10") Integer pageSize) {
        return userService.getUsers(page, pageSize)
                .map(ApiResponse::success);
    }
    
    @Operation(summary = "获取用户详情", description = "根据ID获取用户详细信息")
    @GetMapping("/{id}")
    public Mono<ApiResponse<User>> getUser(@PathVariable Long id) {
        return userService.getUserById(id)
                .map(ApiResponse::success);
    }
    
    @Operation(summary = "创建用户", description = "创建新用户（管理员功能）")
    @PostMapping
    public Mono<ApiResponse<User>> createUser(@Valid @RequestBody RegisterRequest request) {
        return authService.register(request)
                .map(user -> ApiResponse.success("用户创建成功", user));
    }
    
    @Operation(summary = "更新用户", description = "更新用户信息")
    @PutMapping("/{id}")
    public Mono<ApiResponse<User>> updateUser(@PathVariable Long id, 
                                               @RequestBody UpdateUserRequest request) {
        return userService.updateUser(id, request)
                .map(user -> ApiResponse.success("用户更新成功", user));
    }
    
    @Operation(summary = "删除用户", description = "删除指定用户")
    @DeleteMapping("/{id}")
    public Mono<ApiResponse<Void>> deleteUser(@PathVariable Long id) {
        return userService.deleteUser(id)
                .then(Mono.just(ApiResponse.<Void>success("用户删除成功", null)));
    }
    
    public record UpdateUserRequest(
        String nickname,
        String department,
        String company,
        Integer status
    ) {}
}

