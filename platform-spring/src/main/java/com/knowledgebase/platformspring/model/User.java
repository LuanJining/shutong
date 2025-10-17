package com.knowledgebase.platformspring.model;

import java.time.LocalDateTime;

import org.springframework.data.annotation.Id;
import org.springframework.data.relational.core.mapping.Column;
import org.springframework.data.relational.core.mapping.Table;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
@Table("users")
public class User {
    
    @Id
    private Long id;
    
    @Column("username")
    private String username;
    
    @Column("phone")
    private String phone;
    
    @Column("email")
    private String email;
    
    @Column("password")
    private String password;
    
    @Column("nickname")
    private String nickname;
    
    @Column("avatar")
    private String avatar;
    
    @Column("department")
    private String department;
    
    @Column("company")
    private String company;
    
    @Column("status")
    @Builder.Default
    private Integer status = 1; // 1-正常 0-禁用
    
    @Column("last_login")
    private LocalDateTime lastLogin;
    
    @Column("created_at")
    private LocalDateTime createdAt;
    
    @Column("updated_at")
    private LocalDateTime updatedAt;
    
    @Column("deleted_at")
    private LocalDateTime deletedAt;
}

