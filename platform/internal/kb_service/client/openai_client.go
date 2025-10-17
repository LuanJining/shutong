package client

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
	"sync"

	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/common/logger"
	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/kb_service/config"

	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"
	"github.com/openai/openai-go/v2/packages/ssestream"
)

type OpenAIClient struct {
	config *config.OpenAIConfig
	client *openai.Client
	once   sync.Once
	err    error
}

func NewOpenAIClient(config *config.OpenAIConfig) *OpenAIClient {
	return &OpenAIClient{config: config}
}

func (c *OpenAIClient) GetClient() (*openai.Client, error) {
	c.once.Do(func() {
		if c.config == nil {
			c.err = errors.New("config is nil")
			return
		}
		if c.config.ApiKey == "" {
			c.err = errors.New("api key is empty")
			return
		}
		client := openai.NewClient(
			option.WithAPIKey(c.config.ApiKey),
			option.WithBaseURL(c.config.Url),
		)
		c.client = &client
		c.err = nil
	})
	return c.client, c.err
}

// ChatWithFiles 使用提供的文件内容构造上下文并向 OpenAI 发起聊天请求。
func (c *OpenAIClient) ChatWithFiles(ctx context.Context, question string, fileContents []string) (string, error) {
	client, err := c.GetClient()
	if err != nil {
		return "", fmt.Errorf("failed to get OpenAI client: %w", err)
	}

	messages, err := buildChatMessages(question, fileContents)
	if err != nil {
		return "", err
	}
	log.Println("成功构建messages")
	response, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model:               c.config.Model,
		Messages:            messages,
		MaxCompletionTokens: openai.Int(2000),
		Temperature:         openai.Float(0.7),
	})
	if err != nil {
		log.Println("失败获取response", err)
		return "", fmt.Errorf("failed to create chat completion: %w", err)
	}
	log.Println("成功获取response", response)
	if len(response.Choices) == 0 {
		return "", errors.New("no response from OpenAI")
	}

	answer := strings.TrimSpace(response.Choices[0].Message.Content)
	if answer == "" {
		return "", errors.New("empty response from OpenAI")
	}

	log.Println("成功获取answer")
	return answer, nil
}

// ExtractTextFromReader 从 io.Reader 中提取纯文本。
func (c *OpenAIClient) ExtractTextFromReader(_ context.Context, reader io.Reader) (string, error) {
	content, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("failed to read content: %w", err)
	}

	text := strings.TrimSpace(string(content))
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")

	lines := strings.Split(text, "\n")
	cleaned := make([]string, 0, len(lines))
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			cleaned = append(cleaned, trimmed)
		}
	}

	return strings.Join(cleaned, "\n"), nil
}

// ChatWithMinioFiles 从 MinIO 下载文件内容后发起聊天请求。
func (c *OpenAIClient) ChatWithMinioFiles(ctx context.Context, question string, minioClient *S3Client, objectNames []string) (string, error) {
	fileContents, err := c.ExtractMinioFileContents(ctx, minioClient, objectNames)
	if err != nil {
		return "", err
	}

	return c.ChatWithFiles(ctx, question, fileContents)
}

// ChatWithFilesStream 使用 OpenAI 流式接口返回聊天流。
func (c *OpenAIClient) ChatWithFilesStream(ctx context.Context, question string, fileContents []string) (*ssestream.Stream[openai.ChatCompletionChunk], error) {
	client, err := c.GetClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get OpenAI client: %w", err)
	}

	messages, err := buildChatMessages(question, fileContents)
	if err != nil {
		return nil, err
	}

	log.Printf("model: %s", c.config.Model)
	stream := client.Chat.Completions.NewStreaming(ctx, openai.ChatCompletionNewParams{
		Model:               c.config.Model,
		Messages:            messages,
		MaxCompletionTokens: openai.Int(2000),
		Temperature:         openai.Float(0.7),
	})

	if err := stream.Err(); err != nil {
		stream.Close()
		return nil, fmt.Errorf("failed to create streaming chat completion: %w", err)
	}

	return stream, nil
}

// ChatWithMinioFilesStream 从 MinIO 下载文件后发起流式聊天请求。
func (c *OpenAIClient) ChatWithMinioFilesStream(ctx context.Context, question string, minioClient *S3Client, objectNames []string) (*ssestream.Stream[openai.ChatCompletionChunk], error) {
	fileContents, err := c.ExtractMinioFileContents(ctx, minioClient, objectNames)
	if err != nil {
		return nil, err
	}

	return c.ChatWithFilesStream(ctx, question, fileContents)
}

// ExtractMinioFileContents 下载 MinIO 文件并提取文本内容。
func (c *OpenAIClient) ExtractMinioFileContents(ctx context.Context, minioClient *S3Client, objectNames []string) ([]string, error) {
	if minioClient == nil {
		return nil, errors.New("minio client is nil")
	}
	if len(objectNames) == 0 {
		return nil, errors.New("no object names provided")
	}

	fileContents := make([]string, 0, len(objectNames))
	for _, objectName := range objectNames {
		if strings.TrimSpace(objectName) == "" {
			continue
		}

		reader, err := minioClient.DownloadFile(ctx, objectName)
		if err != nil {
			return nil, fmt.Errorf("failed to download file %s: %w", objectName, err)
		}
		log.Println("下载成功")
		content, extractErr := c.ExtractTextFromReader(ctx, reader)
		log.Println("提取成功")
		closeErr := reader.Close()
		log.Println("关闭成功")
		if extractErr != nil {
			return nil, fmt.Errorf("failed to extract text from file %s: %w", objectName, extractErr)
		}
		if closeErr != nil {
			return nil, fmt.Errorf("failed to close reader for file %s: %w", objectName, closeErr)
		}

		fileContents = append(fileContents, content)
	}

	if len(fileContents) == 0 {
		return nil, errors.New("no file contents available")
	}

	return fileContents, nil
}

// CreateEmbedding 调用OpenAI Embeddings API生成向量（使用SDK，支持独立的embedding配置）
func (c *OpenAIClient) CreateEmbedding(ctx context.Context, text string) ([]float64, error) {
	if c.config == nil {
		return nil, errors.New("config is nil")
	}

	text = strings.TrimSpace(text)
	if text == "" {
		return nil, errors.New("text is empty")
	}

	// 如果配置了独立的embedding服务，使用独立配置
	embeddingURL := c.config.EmbeddingURL
	embeddingAPIKey := c.config.EmbeddingAPIKey
	embeddingModel := c.config.EmbeddingModel

	// 如果没有配置独立的embedding服务，使用chat服务的配置
	if embeddingURL == "" {
		embeddingURL = c.config.Url
	}
	if embeddingAPIKey == "" {
		embeddingAPIKey = c.config.ApiKey
	}
	if embeddingModel == "" {
		embeddingModel = openai.EmbeddingModelTextEmbedding3Small
	}

	logger.Debugf("📝 CreateEmbedding - URL: %s, Model: %s, TextLength: %d", embeddingURL, embeddingModel, len(text))

	// 创建独立的embedding客户端
	embClient := openai.NewClient(
		option.WithAPIKey(embeddingAPIKey),
		option.WithBaseURL(embeddingURL),
	)

	logger.Debugf("🔄 Calling embedding API...")
	// 使用SDK的正确方式：EmbeddingNewParamsInputUnion的OfString字段
	response, err := embClient.Embeddings.New(ctx, openai.EmbeddingNewParams{
		Input: openai.EmbeddingNewParamsInputUnion{
			OfString: openai.String(text),
		},
		Model: embeddingModel,
	})
	if err != nil {
		log.Printf("❌ Embedding API call failed - URL: %s, Model: %s, Error: %v", embeddingURL, embeddingModel, err)
		return nil, fmt.Errorf("failed to create embedding: %w", err)
	}
	logger.Debugf("✅ Embedding API call succeeded")

	if len(response.Data) == 0 {
		return nil, errors.New("no embedding data returned")
	}

	embedding := response.Data[0].Embedding
	if len(embedding) == 0 {
		return nil, errors.New("empty embedding returned")
	}

	return embedding, nil
}

// CreateEmbeddingBatch 批量生成向量（使用SDK的批量API，支持独立的embedding配置）
func (c *OpenAIClient) CreateEmbeddingBatch(ctx context.Context, texts []string) ([][]float64, error) {
	if c.config == nil {
		return nil, errors.New("config is nil")
	}

	if len(texts) == 0 {
		return nil, errors.New("texts is empty")
	}

	// 过滤空文本
	validTexts := make([]string, 0, len(texts))
	for _, text := range texts {
		if strings.TrimSpace(text) != "" {
			validTexts = append(validTexts, text)
		}
	}

	if len(validTexts) == 0 {
		return nil, errors.New("no valid texts after filtering")
	}

	// 如果配置了独立的embedding服务，使用独立配置
	embeddingURL := c.config.EmbeddingURL
	embeddingAPIKey := c.config.EmbeddingAPIKey
	embeddingModel := c.config.EmbeddingModel

	// 如果没有配置独立的embedding服务，使用chat服务的配置
	if embeddingURL == "" {
		embeddingURL = c.config.Url
	}
	if embeddingAPIKey == "" {
		embeddingAPIKey = c.config.ApiKey
	}
	if embeddingModel == "" {
		embeddingModel = openai.EmbeddingModelTextEmbedding3Small
	}

	logger.Debugf("📝 CreateEmbeddingBatch - URL: %s, Model: %s, BatchSize: %d", embeddingURL, embeddingModel, len(validTexts))
	logger.Debugf("🔍 First text preview (max 100 chars): %s", truncateString(validTexts[0], 100))

	// 创建独立的embedding客户端
	embClient := openai.NewClient(
		option.WithAPIKey(embeddingAPIKey),
		option.WithBaseURL(embeddingURL),
	)

	logger.Debugf("🔄 Calling batch embedding API...")
	// 使用SDK的批量方式：EmbeddingNewParamsInputUnion的OfArrayOfStrings字段
	response, err := embClient.Embeddings.New(ctx, openai.EmbeddingNewParams{
		Input: openai.EmbeddingNewParamsInputUnion{
			OfArrayOfStrings: validTexts,
		},
		Model: embeddingModel,
	})
	if err != nil {
		log.Printf("❌ Batch embedding API failed - URL: %s, Model: %s, BatchSize: %d, Error: %v", embeddingURL, embeddingModel, len(validTexts), err)
		return nil, fmt.Errorf("failed to create embeddings: %w", err)
	}
	logger.Debugf("✅ Batch embedding API succeeded, got %d embeddings", len(response.Data))

	if len(response.Data) != len(validTexts) {
		return nil, fmt.Errorf("expected %d embeddings but got %d", len(validTexts), len(response.Data))
	}

	embeddings := make([][]float64, len(response.Data))
	for i, data := range response.Data {
		embeddings[i] = data.Embedding
	}

	return embeddings, nil
}

func buildChatMessages(question string, fileContents []string) ([]openai.ChatCompletionMessageParamUnion, error) {
	question = strings.TrimSpace(question)
	if question == "" {
		return nil, errors.New("question is empty")
	}

	var promptBuilder strings.Builder
	promptBuilder.WriteString("你是一个智能助手，可以基于提供的知识库内容回答问题。请根据以下文件内容来回答用户的问题：\n\n")
	for idx, content := range fileContents {
		trimmed := strings.TrimSpace(content)
		if trimmed == "" {
			continue
		}
		fmt.Fprintf(&promptBuilder, "文件 %d 内容：\n%s\n\n", idx+1, trimmed)
	}
	promptBuilder.WriteString("请基于以上知识库内容回答用户的问题。如果问题与文件内容无关，请说明无法从提供的文件中找到相关信息。\n\n")
	promptBuilder.WriteString("用户问题：\n")
	promptBuilder.WriteString(question)

	log.Println("成功构建prompt")

	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(promptBuilder.String()),
		openai.UserMessage(question),
	}
	return messages, nil
}

// truncateString 截断字符串到指定长度
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
