package com.knowledgebase.platformspring.controller;

import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import com.knowledgebase.platformspring.dto.ApiResponse;

import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.tags.Tag;
import reactor.core.publisher.Mono;

@RestController
@RequestMapping("/actuator")
@Tag(name = "Actuator", description = "Actuator endpoints")
public class ActuatorController {

    @Operation(summary = "Readiness probe", description = "Check if the application is ready to serve requests")
    @GetMapping("/health/readiness")
    public Mono<ApiResponse<String>> readiness() {
        return Mono.just(ApiResponse.success("OK"));
    }

    @Operation(summary = "Liveness probe", description = "Check if the application is still running")
    @GetMapping("/health/liveness")
    public Mono<ApiResponse<String>> liveness() {
        return Mono.just(ApiResponse.success("OK"));
        }
    }
