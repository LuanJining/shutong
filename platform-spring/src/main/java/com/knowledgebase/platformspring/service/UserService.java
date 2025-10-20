package com.knowledgebase.platformspring.service;

import java.time.LocalDateTime;
import java.util.List;

import org.springframework.stereotype.Service;

import com.knowledgebase.platformspring.dto.PaginationResponse;
import com.knowledgebase.platformspring.exception.BusinessException;
import com.knowledgebase.platformspring.model.User;
import com.knowledgebase.platformspring.repository.UserRepository;

import lombok.RequiredArgsConstructor;
import reactor.core.publisher.Mono;

@Service
@RequiredArgsConstructor
public class UserService {
    
    private final UserRepository userRepository;
    
    public Mono<PaginationResponse<List<User>>> getUsers(Integer page, Integer pageSize) {
        return userRepository.findAll()
                .collectList()
                .map(allUsers -> {
                    // 清除密码
                    allUsers.forEach(user -> user.setPassword(null));
                    
                    long total = allUsers.size();
                    int offset = (page - 1) * pageSize;
                    List<User> pagedUsers = allUsers.stream()
                            .skip(offset)
                            .limit(pageSize)
                            .toList();
                    
                    return PaginationResponse.of(pagedUsers, total, page, pageSize);
                });
    }
    
    public Mono<User> getUserById(Long id) {
        return userRepository.findById(id)
                .switchIfEmpty(Mono.error(BusinessException.notFound("用户不存在")))
                .doOnSuccess(user -> user.setPassword(null));
    }
    
    public Mono<User> updateUser(Long id, com.knowledgebase.platformspring.controller.UserController.UpdateUserRequest request) {
        return userRepository.findById(id)
                .switchIfEmpty(Mono.error(BusinessException.notFound("用户不存在")))
                .flatMap(user -> {
                    if (request.nickname() != null) {
                        user.setNickname(request.nickname());
                    }
                    if (request.department() != null) {
                        user.setDepartment(request.department());
                    }
                    if (request.company() != null) {
                        user.setCompany(request.company());
                    }
                    if (request.status() != null) {
                        user.setStatus(request.status());
                    }
                    user.setUpdatedAt(LocalDateTime.now());
                    
                    return userRepository.save(user)
                            .doOnSuccess(savedUser -> savedUser.setPassword(null));
                });
    }
    
    public Mono<Void> deleteUser(Long id) {
        return userRepository.findById(id)
                .switchIfEmpty(Mono.error(BusinessException.notFound("用户不存在")))
                .flatMap(userRepository::delete);
    }
}

