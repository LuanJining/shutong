package com.knowledgebase.platformspring.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import jakarta.validation.constraints.NotBlank;
import lombok.Builder;

import java.util.List;

@Builder
public record AcceptSuggestionRequest(
        @JsonProperty("sessionId")
        String sessionId,

        @JsonProperty("fileName")
        String fileName,

        @JsonProperty("fileType")
        String fileType,

        @JsonProperty("acceptedSuggestionIds")
        List<String> acceptedSuggestionIds,

        @JsonProperty("applyAll")
        Boolean applyAll
){}
