package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gatehide/gatehide-api/internal/models"
	"github.com/gatehide/gatehide-api/internal/services"
	"github.com/gatehide/gatehide-api/internal/utils"
	"github.com/gin-gonic/gin"
)

// GamenetHandler handles gamenet HTTP requests
type GamenetHandler struct {
	gamenetService services.GamenetServiceInterface
	fileUploader   *utils.FileUploader
}

// NewGamenetHandler creates a new gamenet handler
func NewGamenetHandler(gamenetService services.GamenetServiceInterface, fileUploader *utils.FileUploader) *GamenetHandler {
	return &GamenetHandler{
		gamenetService: gamenetService,
		fileUploader:   fileUploader,
	}
}

// GetAllGamenets handles GET /gamenets
func (h *GamenetHandler) GetAllGamenets(c *gin.Context) {
	// Check if search parameters are provided
	query := c.Query("query")
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")

	// Debug logging
	fmt.Printf("DEBUG: query=%s, pageStr=%s, pageSizeStr=%s\n", query, pageStr, pageSizeStr)

	// If search parameters are provided, use search endpoint
	if query != "" || pageStr != "" || pageSizeStr != "" {
		fmt.Printf("DEBUG: Using search path\n")
		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			page = 1
		}

		pageSize, err := strconv.Atoi(pageSizeStr)
		if err != nil || pageSize < 1 {
			pageSize = 10
		}
		if pageSize > 100 {
			pageSize = 100
		}

		searchReq := &models.GamenetSearchRequest{
			Query:    query,
			Page:     page,
			PageSize: pageSize,
		}

		result, err := h.gamenetService.Search(c.Request.Context(), searchReq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to search gamenets",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":    "Gamenets retrieved successfully",
			"data":       result.Data,
			"pagination": result.Pagination,
		})
		return
	}

	// Default behavior - get all gamenets
	gamenets, err := h.gamenetService.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve gamenets",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Gamenets retrieved successfully",
		"data":    gamenets,
	})
}

// GetGamenetByID handles GET /gamenets/:id
func (h *GamenetHandler) GetGamenetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid gamenet ID",
		})
		return
	}

	gamenet, err := h.gamenetService.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Gamenet not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Gamenet retrieved successfully",
		"data":    gamenet,
	})
}

// CreateGamenet handles POST /gamenets
func (h *GamenetHandler) CreateGamenet(c *gin.Context) {
	// Parse multipart form
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil { // 32MB max
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to parse multipart form",
		})
		return
	}

	// Get form values
	req := models.GamenetCreateRequest{
		Name:        c.PostForm("name"),
		OwnerName:   c.PostForm("owner_name"),
		OwnerMobile: c.PostForm("owner_mobile"),
		Address:     c.PostForm("address"),
		Email:       c.PostForm("email"),
	}

	// Handle license file upload
	file, fileHeader, err := c.Request.FormFile("license_attachment")
	if err == nil {
		defer file.Close()

		// Upload file
		uploadResult, err := h.fileUploader.UploadFile(fileHeader, "licenses")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to upload license file: " + err.Error(),
			})
			return
		}

		req.LicenseAttachment = &uploadResult.PublicURL
	}

	gamenet, err := h.gamenetService.Create(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Gamenet created successfully",
		"data":    gamenet,
	})
}

// UpdateGamenet handles PUT /gamenets/:id
func (h *GamenetHandler) UpdateGamenet(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid gamenet ID",
		})
		return
	}

	// Parse multipart form
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil { // 32MB max
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to parse multipart form",
		})
		return
	}

	// Get form values
	req := models.GamenetUpdateRequest{}

	if name := c.PostForm("name"); name != "" {
		req.Name = &name
	}
	if ownerName := c.PostForm("owner_name"); ownerName != "" {
		req.OwnerName = &ownerName
	}
	if ownerMobile := c.PostForm("owner_mobile"); ownerMobile != "" {
		req.OwnerMobile = &ownerMobile
	}
	if address := c.PostForm("address"); address != "" {
		req.Address = &address
	}
	if email := c.PostForm("email"); email != "" {
		req.Email = &email
	}

	// Handle license file upload
	file, fileHeader, err := c.Request.FormFile("license_attachment")
	if err == nil {
		defer file.Close()

		// Get current gamenet to check for existing license
		currentGamenet, err := h.gamenetService.GetByID(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Gamenet not found",
			})
			return
		}

		// Delete old license file if it exists
		if currentGamenet.LicenseAttachment != nil && *currentGamenet.LicenseAttachment != "" {
			// Extract file path from public URL
			oldFilePath := h.extractFilePathFromURL(*currentGamenet.LicenseAttachment)
			fmt.Printf("DEBUG: Original license URL: %s\n", *currentGamenet.LicenseAttachment)
			fmt.Printf("DEBUG: Extracted file path: %s\n", oldFilePath)
			if oldFilePath != "" {
				if err := h.fileUploader.DeleteFile(oldFilePath); err != nil {
					// Log error but don't fail the update
					fmt.Printf("Warning: Failed to delete old license file: %v\n", err)
				} else {
					fmt.Printf("DEBUG: Successfully deleted old file: %s\n", oldFilePath)
				}
			} else {
				fmt.Printf("DEBUG: Could not extract file path from URL: %s\n", *currentGamenet.LicenseAttachment)
			}
		}

		// Upload new file
		uploadResult, err := h.fileUploader.UploadFile(fileHeader, "licenses")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to upload license file: " + err.Error(),
			})
			return
		}

		req.LicenseAttachment = &uploadResult.PublicURL
	}

	gamenet, err := h.gamenetService.Update(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Gamenet updated successfully",
		"data":    gamenet,
	})
}

// DeleteGamenet handles DELETE /gamenets/:id
func (h *GamenetHandler) DeleteGamenet(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid gamenet ID",
		})
		return
	}

	err = h.gamenetService.Delete(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Gamenet deleted successfully",
	})
}

// extractFilePathFromURL extracts the file path from a public URL or relative path
func (h *GamenetHandler) extractFilePathFromURL(publicURL string) string {
	// Handle both formats:
	// 1. Full URL: http://localhost:8080/uploads/licenses/filename.ext
	// 2. Relative path: /uploads/licenses/filename.ext
	// We need to extract: ./uploads/licenses/filename.ext

	var relativePath string

	if strings.HasPrefix(publicURL, "http://") {
		// Full URL format
		prefix := "http://localhost:8080/uploads/"
		if !strings.HasPrefix(publicURL, prefix) {
			return ""
		}
		relativePath = strings.TrimPrefix(publicURL, prefix)
	} else if strings.HasPrefix(publicURL, "/uploads/") {
		// Relative path format
		relativePath = strings.TrimPrefix(publicURL, "/uploads/")
	} else {
		return ""
	}

	// Construct the full file path
	return fmt.Sprintf("./uploads/%s", relativePath)
}
