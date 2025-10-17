package com.knowledgebase.platformspring.model;

import java.util.ArrayList;
import java.util.List;

import org.springframework.data.relational.core.mapping.Column;
import org.springframework.data.relational.core.mapping.Table;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.core.type.TypeReference;
import com.fasterxml.jackson.databind.ObjectMapper;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
@Table("space_members")
public class SpaceMember {
    
    @Column("space_id")
    private Long spaceId;
    
    @Column("user_id")
    private Long userId;
    
    @Column("roles")
    private String rolesJson; // JSON string storing list of roles
    
    // Helper methods to convert between List<String> and JSON string
    public List<String> getRoles() {
        if (rolesJson == null || rolesJson.isEmpty()) {
            return new ArrayList<>();
        }
        try {
            ObjectMapper mapper = new ObjectMapper();
            return mapper.readValue(rolesJson, new TypeReference<List<String>>() {});
        } catch (JsonProcessingException e) {
            return new ArrayList<>();
        }
    }
    
    public void setRoles(List<String> roles) {
        try {
            ObjectMapper mapper = new ObjectMapper();
            this.rolesJson = mapper.writeValueAsString(roles);
        } catch (JsonProcessingException e) {
            this.rolesJson = "[]";
        }
    }
    
    // Space member role constants
    public static final String ROLE_ADMIN = "admin";
    public static final String ROLE_APPROVER = "approver";
    public static final String ROLE_EDITOR = "editor";
    public static final String ROLE_READER = "reader";
}

