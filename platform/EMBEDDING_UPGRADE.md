# 向量化升级说明

## 问题背景

之前的实现使用了`simpleEmbedding`函数，只是简单统计文本特征（字符数、单词数等7个统计特征），完全没有语义信息。这导致向量搜索效果很差，无法真正理解文本的语义。

## 解决方案

### Qdrant是否自带嵌入模型？

**答案：不自带**。Qdrant只是一个向量数据库，负责存储和搜索向量，但不负责生成向量。你需要自己提供嵌入模型来把文本转换成向量。

### 升级内容

现在已经集成了OpenAI的Embeddings API（使用SDK方式），将文本转换为真正具有语义的向量：

1. **OpenAI Embeddings API**
   - 模型：`text-embedding-3-small`
   - 维度：1536维
   - 优点：语义理解能力强，性价比高

2. **降级策略**
   - 优先使用OpenAI生成向量
   - 如果OpenAI不可用，自动降级到简单统计特征（保证系统可用性）

## 代码修改

### 1. OpenAI Client 新增方法

```go
// CreateEmbedding 单个文本生成向量
func (c *OpenAIClient) CreateEmbedding(ctx context.Context, text string) ([]float64, error)

// CreateEmbeddingBatch 批量生成向量（更高效）
func (c *OpenAIClient) CreateEmbeddingBatch(ctx context.Context, texts []string) ([][]float64, error)
```

**SDK正确用法**：
```go
// 单个文本
response, err := client.Embeddings.New(ctx, openai.EmbeddingNewParams{
    Input: openai.EmbeddingNewParamsInputUnion{
        OfString: openai.String(text),  // 关键：使用OfString字段
    },
    Model: openai.EmbeddingModelTextEmbedding3Small,
})

// 多个文本（批量）
response, err := client.Embeddings.New(ctx, openai.EmbeddingNewParams{
    Input: openai.EmbeddingNewParamsInputUnion{
        OfArrayOfStrings: []string{"text1", "text2"},  // 关键：使用OfArrayOfStrings字段
    },
    Model: openai.EmbeddingModelTextEmbedding3Small,
})
```

### 2. DocumentService 新增方法

```go
// generateEmbedding 生成单个文本的向量（带降级策略）
func (s *DocumentService) generateEmbedding(ctx context.Context, text string) ([]float64, error)

// generateEmbeddingBatch 批量生成向量（带降级策略）
func (s *DocumentService) generateEmbeddingBatch(ctx context.Context, texts []string) ([][]float64, error)
```

### 3. 配置更新

**配置文件** (`internal/kb_service/config/localtest.yaml`):
```yaml
vector:
  base_url: http://localhost:6333
  api_key: ""
  collection: kb_documents
  vector_size: 1536  # 从7改为1536，匹配OpenAI embedding维度
  distance: Cosine
```

**代码默认值** (`internal/kb_service/client/qdrant_client.go`):
```go
vectorSize := cfg.VectorSize
if vectorSize <= 0 {
    vectorSize = 1536 // 默认使用OpenAI text-embedding-3-small 向量维度
}
```

## 使用流程

### 文档上传和向量化

```
用户上传文档
    ↓
文档解析/OCR提取文本
    ↓
文本分块（800字符，重叠120字符）
    ↓
批量生成向量（OpenAI Embeddings API）
    ↓
存储到Qdrant向量数据库
```

### 知识搜索

```
用户输入查询
    ↓
生成查询向量（OpenAI Embeddings API）
    ↓
在Qdrant中进行向量相似度搜索
    ↓
返回最相关的文档片段
```

## 重要提示

### ⚠️ 需要重新处理现有文档

**重要**：因为向量维度从7改为1536，**所有已存在的文档向量都需要重新生成**。

有两个选择：

**方案1：删除Qdrant中的旧集合（推荐）**
```bash
# 删除旧集合，会自动用新维度重建
curl -X DELETE "http://localhost:6333/collections/kb_documents"
```

**方案2：重新处理所有文档**
- 遍历数据库中的所有文档
- 调用`ProcessDocument`重新生成向量

### 💰 成本考虑

使用OpenAI Embeddings API会产生费用：
- `text-embedding-3-small`: $0.02 / 1M tokens
- 示例：1000个文档，每个800字 ≈ 200K tokens ≈ $0.004

### 🔄 降级策略

系统设计了降级策略，保证可用性：
```go
if s.openaiClient != nil {
    embedding, err := s.openaiClient.CreateEmbedding(ctx, text)
    if err != nil {
        log.Printf("OpenAI embedding failed, falling back to simple embedding: %v", err)
        return simpleEmbedding(text), nil  // 降级到统计特征
    }
    return embedding, nil
}
return simpleEmbedding(text), nil  // OpenAI客户端未配置时使用统计特征
```

## 其他可选方案

如果不想使用OpenAI，还有以下替代方案：

### 1. 本地模型（完全免费）

使用开源embedding模型，如：
- **sentence-transformers**（Python）
- **text2vec**（Go）
- **bge-small-zh**（中文优化）

### 2. FastEmbed（Qdrant官方推荐）

Qdrant提供了FastEmbed库（Python），支持多种本地模型。

### 3. 其他API服务

- Cohere Embeddings
- Google Vertex AI
- Azure OpenAI
- 本地部署的Ollama + embedding模型

## SDK关键知识点

OpenAI Go SDK v2的`EmbeddingNewParamsInputUnion`是一个联合类型，包含以下字段：

- `OfString`: 单个字符串（`param.Opt[string]`）
- `OfArrayOfStrings`: 字符串数组（`[]string`）
- `OfArrayOfTokens`: token数组（`[]int64`）
- `OfArrayOfTokenArrays`: token数组的数组（`[][]int64`）

**不能直接赋值字符串**，必须设置对应的字段：
```go
// ✅ 正确
Input: openai.EmbeddingNewParamsInputUnion{
    OfString: openai.String(text),
}

// ❌ 错误
Input: text  // 类型不匹配
```

## 测试建议

1. 上传一个新文档，检查向量是否正确生成（维度1536）
2. 进行知识搜索，对比之前的搜索效果
3. 检查日志，确认没有降级到simpleEmbedding

## 总结

通过这次升级：
- ✅ 使用OpenAI Embeddings API生成真正的语义向量
- ✅ 向量维度从7维升级到1536维
- ✅ 使用官方SDK而非HTTP直接调用
- ✅ 支持降级策略，保证系统可用性
- ✅ 批量处理提高效率

向量搜索效果将得到**显著提升**！🎉

