package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	openai "github.com/openai/openai-go/v2"

	model "gitee.com/sichuan-shutong-zhihui-data/k-base/internal/common/models"
	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/kb_service/service"
	"github.com/gin-gonic/gin"
)

type DocumentHandler struct {
	documentService *service.DocumentService
}

func NewDocumentHandler(documentService *service.DocumentService) *DocumentHandler {
	return &DocumentHandler{
		documentService: documentService,
	}
}

// UploadDocument 上传文档
func (h *DocumentHandler) UploadDocument(c *gin.Context) {
	// 获取上传的文件
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, &model.APIResponse{
			Code:    http.StatusBadRequest,
			Message: "Failed to get uploaded file",
		})
		return
	}
	defer file.Close()

	// 获取其他参数
	fileName := c.PostForm("file_name")
	spaceIDStr := c.PostForm("space_id")
	subSpaceIDStr := c.PostForm("sub_space_id")
	classIDStr := c.PostForm("class_id")
	tags := c.PostForm("tags")
	summary := c.PostForm("summary")
	department := c.PostForm("department")
	needApprovalStr := c.PostForm("need_approval")
	version := c.PostForm("version")
	useType := c.PostForm("use_type")
	// 默认需要审批
	needApproval := true
	if needApprovalStr != "" {
		needApproval, err = strconv.ParseBool(needApprovalStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, &model.APIResponse{
				Code:    http.StatusBadRequest,
				Message: "Invalid need_approval",
			})
			return
		}
	}

	// 验证必需参数
	if spaceIDStr == "" {
		c.JSON(http.StatusBadRequest, &model.APIResponse{
			Code:    http.StatusBadRequest,
			Message: "space_id is required",
		})
		return
	}

	if subSpaceIDStr == "" {
		c.JSON(http.StatusBadRequest, &model.APIResponse{
			Code:    http.StatusBadRequest,
			Message: "sub_space_id is required",
		})
		return
	}

	spaceID, err := strconv.ParseUint(spaceIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, &model.APIResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid space_id",
		})
		return
	}

	subSpaceID, err := strconv.ParseUint(subSpaceIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, &model.APIResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid sub_space_id",
		})
		return
	}
	if classIDStr == "" {
		c.JSON(http.StatusBadRequest, &model.APIResponse{
			Code:    http.StatusBadRequest,
			Message: "class_id is required",
		})
		return
	}

	classID, err := strconv.ParseUint(classIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, &model.APIResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid class_id",
		})
		return
	}

	// 获取实际文件大小
	var actualSize int64
	if seeker, ok := file.(io.Seeker); ok {
		currentPos, _ := seeker.Seek(0, io.SeekCurrent)
		endPos, _ := seeker.Seek(0, io.SeekEnd)
		actualSize = endPos - currentPos
		log.Printf("实际文件大小: %d", actualSize)
		// 重置到开始位置
		seeker.Seek(currentPos, io.SeekStart)
	} else {
		actualSize = header.Size
	}

	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, &model.APIResponse{
			Code:    http.StatusUnauthorized,
			Message: "User not found",
		})
		return
	}

	// 构建服务请求
	req := &model.UploadDocumentRequest{
		File:            file,
		FileName:        fileName,
		FileSize:        actualSize,
		ContentType:     header.Header.Get("Content-Type"),
		SpaceID:         uint(spaceID),
		SubSpaceID:      uint(subSpaceID),
		ClassID:         uint(classID),
		Tags:            tags,
		Summary:         summary,
		CreatedBy:       user.(*model.User).ID,
		CreatorNickName: user.(*model.User).Nickname,
		Department:      department,
		NeedApproval:    needApproval,
		Version:         version,
		UseType:         model.UseType(useType),
	}

	// 调用服务层
	document, err := h.documentService.UploadDocument(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &model.APIResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to upload document: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &model.APIResponse{
		Code:    http.StatusOK,
		Message: "Document uploaded successfully",
		Data:    document,
	})
}

// GetDocument 获取文档详情
func (h *DocumentHandler) GetDocument(c *gin.Context) {
	documentIDStr := c.Param("id")
	documentID, err := strconv.ParseUint(documentIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, &model.APIResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid document ID",
		})
		return
	}

	document, err := h.documentService.GetDocument(c.Request.Context(), uint(documentID))
	if err != nil {
		if err.Error() == "document not found" {
			c.JSON(http.StatusNotFound, &model.APIResponse{
				Code:    http.StatusNotFound,
				Message: "Document not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, &model.APIResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get document: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &model.APIResponse{
		Code:    http.StatusOK,
		Message: "Success",
		Data:    document,
	})
}

// GetDocumentsBySpaceId 获取空间下的文档
func (h *DocumentHandler) GetDocumentsBySpaceId(c *gin.Context) {
	spaceIDStr := c.Param("id")
	spaceID, err := strconv.ParseUint(spaceIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, &model.APIResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid space ID",
		})
		return
	}

	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	documents, err := h.documentService.GetDocumentsBySpaceId(c.Request.Context(), uint(spaceID), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &model.APIResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get documents: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &model.APIResponse{
		Code:    http.StatusOK,
		Message: "Success",
		Data:    documents,
	})
}

// GetTagCloud 获取标签云
func (h *DocumentHandler) GetTagCloud(c *gin.Context) {
	spaceIDStr := c.Query("space_id")
	subSpaceIDStr := c.Query("sub_space_id")
	limitStr := c.DefaultQuery("limit", "50")

	var (
		spaceID    uint64
		subSpaceID uint64
		err        error
	)

	if spaceIDStr != "" {
		spaceID, err = strconv.ParseUint(spaceIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, &model.APIResponse{
				Code:    http.StatusBadRequest,
				Message: "Invalid space_id",
			})
			return
		}
	}

	if subSpaceIDStr != "" {
		subSpaceID, err = strconv.ParseUint(subSpaceIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, &model.APIResponse{
				Code:    http.StatusBadRequest,
				Message: "Invalid sub_space_id",
			})
			return
		}
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, &model.APIResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid limit",
		})
		return
	}

	items, err := h.documentService.GetTagCloud(
		c.Request.Context(),
		uint(spaceID),
		uint(subSpaceID),
		limit,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &model.APIResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get tag cloud: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &model.APIResponse{
		Code:    http.StatusOK,
		Message: "Success",
		Data: model.TagCloudResponse{
			Items: items,
		},
	})
}

// SearchKnowledge 知识检索
func (h *DocumentHandler) SearchKnowledge(c *gin.Context) {
	var req model.KnowledgeSearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, &model.APIResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
		})
		return
	}

	results, err := h.documentService.SearchKnowledge(c.Request.Context(), &req)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrVectorClientNotConfigured):
			c.JSON(http.StatusInternalServerError, &model.APIResponse{
				Code:    http.StatusInternalServerError,
				Message: "Vector search not configured",
			})
		default:
			c.JSON(http.StatusInternalServerError, &model.APIResponse{
				Code:    http.StatusInternalServerError,
				Message: "Failed to search knowledge: " + err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, &model.APIResponse{
		Code:    http.StatusOK,
		Message: "Success",
		Data:    results,
	})
}

// RetryProcessDocument 重试处理文档
func (h *DocumentHandler) RetryProcessDocument(c *gin.Context) {
	var req model.RetryProcessRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, &model.APIResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	result, err := h.documentService.RetryProcessDocument(c.Request.Context(), req.DocumentID, req.ForceRetry)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &model.APIResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to retry process document: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &model.APIResponse{
		Code:    http.StatusOK,
		Message: "Success",
		Data:    result,
	})
}

// PreviewDocument 预览文档（支持浏览器内嵌显示）
func (h *DocumentHandler) PreviewDocument(c *gin.Context) {
	documentIDStr := c.Param("id")
	documentID, err := strconv.ParseUint(documentIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, &model.APIResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid document ID",
		})
		return
	}

	document, err := h.documentService.GetDocument(c.Request.Context(), uint(documentID))
	if err != nil {
		if err.Error() == "document not found" {
			c.JSON(http.StatusNotFound, &model.APIResponse{
				Code:    http.StatusNotFound,
				Message: "Document not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, &model.APIResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get document: " + err.Error(),
		})
		return
	}

	// 根据文件类型设置不同的Content-Type
	var contentType string
	switch document.FileType {
	case ".pdf":
		contentType = "application/pdf"
	case ".jpg", ".jpeg":
		contentType = "image/jpeg"
	case ".png":
		contentType = "image/png"
	case ".gif":
		contentType = "image/gif"
	case ".txt":
		contentType = "text/plain"
	case ".html", ".htm":
		contentType = "text/html"
	default:
		contentType = "application/octet-stream"
	}

	// 设置响应头，支持浏览器内嵌显示
	c.Header("Content-Type", contentType)
	c.Header("Content-Disposition", "inline; filename=\""+document.FileName+"\"")
	c.Header("Cache-Control", "public, max-age=3600")

	fileReader, err := h.documentService.DownloadDocument(c.Request.Context(), uint(documentID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, &model.APIResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to download document: " + err.Error(),
		})
		return
	}
	defer fileReader.Close()

	// 使用CopyBuffer提高性能，添加错误处理
	buffer := make([]byte, 64*1024) // 64KB buffer for better performance
	_, err = io.CopyBuffer(c.Writer, fileReader, buffer)
	if err != nil {
		// 检查是否是客户端取消或连接断开
		if err == io.EOF {
			// 正常结束
			return
		}
		// 检查context是否被取消
		select {
		case <-c.Request.Context().Done():
			// 客户端取消请求，不记录错误
			return
		default:
			// 其他错误
			fmt.Printf("Error copying file: %v\n", err)
		}
	}
}

// TODO 提交文档
func (h *DocumentHandler) SubmitDocument(c *gin.Context) {
	documentIDStr := c.Param("id")
	documentID, err := strconv.ParseUint(documentIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, &model.APIResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid document ID",
		})
		return
	}

	document, err := h.documentService.GetDocument(c.Request.Context(), uint(documentID))
	if err != nil {
		if err.Error() == "document not found" {
			c.JSON(http.StatusNotFound, &model.APIResponse{
				Code:    http.StatusNotFound,
				Message: "Document not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, &model.APIResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get document: " + err.Error(),
		})
		return
	}

	// TODO 不需要审批 直接OCR然后向量化，最后发布
	if !document.NeedApproval {
		c.JSON(http.StatusOK, &model.APIResponse{
			Code:    http.StatusOK,
			Message: "Document does not need approval",
		})
		return
	}

	// 需要审批 创建审批流程
	document, err = h.documentService.CreateWorkflow(c.Request.Context(), document)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &model.APIResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to create workflow: " + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, document)
}

// PublishDocument 发布文档
func (h *DocumentHandler) PublishDocument(c *gin.Context) {
	documentIDStr := c.Param("id")
	documentID, err := strconv.ParseUint(documentIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, &model.APIResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid document ID",
		})
		return
	}

	document, err := h.documentService.GetDocument(c.Request.Context(), uint(documentID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, &model.APIResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get document: " + err.Error(),
		})
		return
	}

	if document.Status != model.DocumentStatusPendingPublish {
		c.JSON(http.StatusBadRequest, &model.APIResponse{
			Code:    http.StatusBadRequest,
			Message: "Document is not pending approval",
		})
		return
	}

	document.Status = model.DocumentStatusPublished
	err = h.documentService.PublishDocument(c.Request.Context(), document)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &model.APIResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to publish document: " + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, &model.APIResponse{
		Code:    http.StatusOK,
		Message: "Document published successfully",
		Data:    document,
	})
}

func (h *DocumentHandler) UnpublishDocument(c *gin.Context) {
	documentIDStr := c.Param("id")
	documentID, err := strconv.ParseUint(documentIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, &model.APIResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid document ID",
		})
		return
	}

	document, err := h.documentService.GetDocument(c.Request.Context(), uint(documentID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, &model.APIResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get document: " + err.Error(),
		})
		return
	}

	if document.Status != model.DocumentStatusPublished {
		c.JSON(http.StatusBadRequest, &model.APIResponse{
			Code:    http.StatusBadRequest,
			Message: "Document is not published",
		})
		return
	}

	document.Status = model.DocumentStatusPendingPublish
	err = h.documentService.UnpublishDocument(c.Request.Context(), document)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &model.APIResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to unpublish document: " + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, &model.APIResponse{
		Code:    http.StatusOK,
		Message: "Document unpublished successfully",
		Data:    document,
	})
}

func (h *DocumentHandler) ChatDocument(c *gin.Context) {
	spaceIDStr := c.Param("id")
	spaceID, err := strconv.ParseUint(spaceIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, &model.APIResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid space ID",
		})
		return
	}

	var chatReq model.ChatDocumentRequest
	if err := c.ShouldBindJSON(&chatReq); err != nil {
		c.JSON(http.StatusBadRequest, &model.APIResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
		})
		return
	}

	resp, err := h.documentService.ChatDocument(c.Request.Context(), uint(spaceID), &chatReq)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrEmptyChatQuestion):
			c.JSON(http.StatusBadRequest, &model.APIResponse{
				Code:    http.StatusBadRequest,
				Message: "Empty chat question",
			})
		case errors.Is(err, model.ErrNoDocumentsAvailable):
			c.JSON(http.StatusNotFound, &model.APIResponse{
				Code:    http.StatusNotFound,
				Message: "No documents available for chat",
			})
		case errors.Is(err, model.ErrOpenAIClientNotConfigured):
			c.JSON(http.StatusInternalServerError, &model.APIResponse{
				Code:    http.StatusInternalServerError,
				Message: "OpenAI client not configured: " + err.Error(),
			})
		default:
			c.JSON(http.StatusInternalServerError, &model.APIResponse{
				Code:    http.StatusInternalServerError,
				Message: "Failed to chat document: " + err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *DocumentHandler) ChatDocumentStream(c *gin.Context) {
	spaceIDStr := c.Param("id")
	spaceID, err := strconv.ParseUint(spaceIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, &model.APIResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid space ID",
		})
		return
	}

	var chatReq model.ChatDocumentRequest
	if err := c.ShouldBindJSON(&chatReq); err != nil {
		c.JSON(http.StatusBadRequest, &model.APIResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
		})
		return
	}

	result, err := h.documentService.ChatDocumentStream(c.Request.Context(), uint(spaceID), &chatReq)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrEmptyChatQuestion):
			c.JSON(http.StatusBadRequest, &model.APIResponse{
				Code:    http.StatusBadRequest,
				Message: "Empty chat question",
			})
		case errors.Is(err, model.ErrNoDocumentsAvailable):
			c.JSON(http.StatusNotFound, &model.APIResponse{
				Code:    http.StatusNotFound,
				Message: "No documents available for chat",
			})
		case errors.Is(err, model.ErrOpenAIClientNotConfigured):
			c.JSON(http.StatusInternalServerError, &model.APIResponse{
				Code:    http.StatusInternalServerError,
				Message: "OpenAI client not configured: " + err.Error(),
			})
		default:
			c.JSON(http.StatusInternalServerError, &model.APIResponse{
				Code:    http.StatusInternalServerError,
				Message: "Failed to chat document: " + err.Error(),
			})
		}
		return
	}

	writer := c.Writer
	flusher, ok := writer.(http.Flusher)
	if !ok {
		c.JSON(http.StatusInternalServerError, &model.APIResponse{
			Code:    http.StatusInternalServerError,
			Message: "Streaming not supported",
		})
		return
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")
	writer.WriteHeader(http.StatusOK)

	if err := writeSSE(writer, flusher, "sources", result.Sources); err != nil {
		log.Printf("failed to write sources SSE: %v", err)
		return
	}

	stream := result.Stream
	defer stream.Close()
	acc := openai.ChatCompletionAccumulator{}

	for stream.Next() {
		chunk := stream.Current()
		_ = acc.AddChunk(chunk)

		if len(chunk.Choices) == 0 {
			continue
		}

		for _, choice := range chunk.Choices {
			if delta := choice.Delta.Content; delta != "" {
				if err := writeSSE(writer, flusher, "token", map[string]string{"content": delta}); err != nil {
					log.Printf("failed to write token SSE: %v", err)
					return
				}
			}
			if refusal := choice.Delta.Refusal; refusal != "" {
				if err := writeSSE(writer, flusher, "refusal", map[string]string{"content": refusal}); err != nil {
					log.Printf("failed to write refusal SSE: %v", err)
					return
				}
			}
		}
	}

	if err := stream.Err(); err != nil {
		_ = writeSSE(writer, flusher, "error", map[string]string{"message": err.Error()})
		return
	}

	finalMessage := ""
	if len(acc.Choices) > 0 {
		finalMessage = strings.TrimSpace(acc.Choices[0].Message.Content)
	}

	if err := writeSSE(writer, flusher, "done", map[string]any{
		"message": finalMessage,
		"sources": result.Sources,
	}); err != nil {
		log.Printf("failed to write done SSE: %v", err)
	}
}

func writeSSE(w http.ResponseWriter, flusher http.Flusher, event string, payload any) error {
	var builder strings.Builder
	if event != "" {
		builder.WriteString("event: ")
		builder.WriteString(event)
		builder.WriteString("\n")
	}

	switch v := payload.(type) {
	case string:
		lines := strings.Split(v, "\n")
		if len(lines) == 0 {
			lines = []string{""}
		}
		for _, line := range lines {
			builder.WriteString("data: ")
			builder.WriteString(line)
			builder.WriteString("\n")
		}
	default:
		data, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		builder.WriteString("data: ")
		builder.Write(data)
		builder.WriteString("\n")
	}

	builder.WriteString("\n")

	if _, err := w.Write([]byte(builder.String())); err != nil {
		return err
	}

	flusher.Flush()
	return nil
}

func (h *DocumentHandler) DeleteDocument(c *gin.Context) {
	documentIDStr := c.Param("id")
	documentID, err := strconv.ParseUint(documentIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, &model.APIResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid document ID",
		})
		return
	}

	err = h.documentService.DeleteDocument(c.Request.Context(), uint(documentID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, &model.APIResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to delete document: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &model.APIResponse{
		Code:    http.StatusNoContent,
		Message: "Document deleted successfully",
	})
}

func (h *DocumentHandler) ApproveDocument(c *gin.Context) {
	documentIDStr := c.Param("id")
	documentID, err := strconv.ParseUint(documentIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, &model.APIResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid document ID",
		})
		return
	}
	log.Printf("documentID: %d", documentID)
}

func (h *DocumentHandler) GetHomepageDocuments(c *gin.Context) {
	documents, err := h.documentService.GetHomepageDocuments(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, &model.APIResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get homepage documents: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &model.APIResponse{
		Code:    http.StatusOK,
		Message: "Success",
		Data:    documents,
	})
}
