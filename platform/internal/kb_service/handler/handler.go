package handler

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get uploaded file"})
		return
	}
	defer file.Close()

	// 获取其他参数
	spaceIDStr := c.PostForm("space_id")
	visibility := c.PostForm("visibility")
	urgency := c.PostForm("urgency")
	tags := c.PostForm("tags")
	summary := c.PostForm("summary")
	createdByStr := c.PostForm("created_by")
	department := c.PostForm("department")

	// 验证必需参数
	if spaceIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "space_id is required"})
		return
	}

	spaceID, err := strconv.ParseUint(spaceIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid space_id"})
		return
	}

	// 解析created_by
	var createdBy uint
	if createdByStr != "" {
		createdByUint, err := strconv.ParseUint(createdByStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid created_by"})
			return
		}
		createdBy = uint(createdByUint)
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

	// 构建服务请求
	req := &service.UploadDocumentRequest{
		File:        file,
		FileName:    header.Filename,
		FileSize:    actualSize,
		ContentType: header.Header.Get("Content-Type"),
		SpaceID:     uint(spaceID),
		Visibility:  visibility,
		Urgency:     urgency,
		Tags:        tags,
		Summary:     summary,
		CreatedBy:   createdBy,
		Department:  department,
	}

	// 调用服务层
	resp, err := h.documentService.UploadDocument(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetDocument 获取文档详情
func (h *DocumentHandler) GetDocument(c *gin.Context) {
	documentIDStr := c.Param("id")
	documentID, err := strconv.ParseUint(documentIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	document, err := h.documentService.GetDocument(c.Request.Context(), uint(documentID))
	if err != nil {
		if err.Error() == "document not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, document)
}

// DownloadDocument 下载文档
func (h *DocumentHandler) DownloadDocument(c *gin.Context) {
	documentIDStr := c.Param("id")
	documentID, err := strconv.ParseUint(documentIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	// 获取文档信息
	document, err := h.documentService.GetDocument(c.Request.Context(), uint(documentID))
	if err != nil {
		if err.Error() == "document not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 设置响应头
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", document.FileName))
	c.Header("Content-Type", document.MimeType)
	c.Header("Content-Length", fmt.Sprintf("%d", document.FileSize))

	// 添加调试信息
	fmt.Printf("Starting download: %s (%.2f KB)\n", document.FileName, float64(document.FileSize)/1024)

	// 下载文件
	fileReader, err := h.documentService.DownloadDocument(c.Request.Context(), uint(documentID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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

// PreviewDocument 预览文档（支持浏览器内嵌显示）
func (h *DocumentHandler) PreviewDocument(c *gin.Context) {
	documentIDStr := c.Param("id")
	documentID, err := strconv.ParseUint(documentIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	document, err := h.documentService.GetDocument(c.Request.Context(), uint(documentID))
	if err != nil {
		if err.Error() == "document not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
