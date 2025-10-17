package com.knowledgebase.platformspring.exception;

import org.springframework.http.HttpStatus;

import lombok.Getter;

@Getter
public class BusinessException extends RuntimeException {
    
    private final HttpStatus status;
    
    public BusinessException(String message) {
        super(message);
        this.status = HttpStatus.BAD_REQUEST;
    }
    
    public BusinessException(String message, HttpStatus status) {
        super(message);
        this.status = status;
    }
    
    public BusinessException(String message, Throwable cause) {
        super(message, cause);
        this.status = HttpStatus.BAD_REQUEST;
    }
    
    public static BusinessException notFound(String message) {
        return new BusinessException(message, HttpStatus.NOT_FOUND);
    }
    
    public static BusinessException unauthorized(String message) {
        return new BusinessException(message, HttpStatus.UNAUTHORIZED);
    }
    
    public static BusinessException forbidden(String message) {
        return new BusinessException(message, HttpStatus.FORBIDDEN);
    }
}

