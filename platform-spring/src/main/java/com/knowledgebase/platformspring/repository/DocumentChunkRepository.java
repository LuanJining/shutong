package com.knowledgebase.platformspring.repository;

import org.springframework.data.r2dbc.repository.Query;
import org.springframework.data.r2dbc.repository.R2dbcRepository;
import org.springframework.stereotype.Repository;

import com.knowledgebase.platformspring.model.DocumentChunk;

import reactor.core.publisher.Flux;
import reactor.core.publisher.Mono;

@Repository
public interface DocumentChunkRepository extends R2dbcRepository<DocumentChunk, Long> {
    
    Flux<DocumentChunk> findByDocumentId(Long documentId);
    
    Flux<DocumentChunk> findByDocumentIdOrderByChunkIndexAsc(Long documentId);
    
    Mono<Long> countByDocumentId(Long documentId);
    
    @Query("DELETE FROM document_chunks WHERE document_id = :documentId")
    Mono<Void> deleteByDocumentId(Long documentId);
}

