package com.knowledgebase.platformspring.config;

import org.springframework.context.annotation.Configuration;
import org.springframework.transaction.annotation.EnableTransactionManagement;

/**
 * 事务配置 - 启用响应式事务管理
 */
@Configuration
@EnableTransactionManagement
public class TransactionConfig {
    // R2DBC 的事务管理器会自动配置
    // 注意：@EnableR2dbcRepositories 已在 PlatformSpringApplication 中声明，不需要重复
}

