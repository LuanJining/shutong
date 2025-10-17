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
@Table("document_chunks")
public class DocumentChunk {
    
    @Id
    private Long id;
    
    @Column("document_id")
    private Long documentId;
    
    @Column("chunk_index")
    private Integer chunkIndex;
    
    @Column("content")
    private String content;
    
    @Column("vector_id")
    private String vectorId; // UUID in Qdrant
    
    @Column("metadata")
    private String metadata; // JSON string for additional metadata
    
    @Column("created_at")
    private LocalDateTime createdAt;
    
    @Column("updated_at")
    private LocalDateTime updatedAt;
}

