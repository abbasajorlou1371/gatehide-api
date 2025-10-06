package models

import "time"

// Gamenet represents a gaming center in the system
type Gamenet struct {
	ID                int       `json:"id" db:"id"`
	Name              string    `json:"name" db:"name"`
	OwnerName         string    `json:"owner_name" db:"owner_name"`
	OwnerMobile       string    `json:"owner_mobile" db:"owner_mobile"`
	Address           string    `json:"address" db:"address"`
	Email             string    `json:"email" db:"email"`
	Password          string    `json:"-" db:"password"` // Hidden from JSON
	LicenseAttachment *string   `json:"license_attachment" db:"license_attachment"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
}

// GamenetCreateRequest represents a gamenet creation request
type GamenetCreateRequest struct {
	Name              string  `json:"name" binding:"required"`
	OwnerName         string  `json:"owner_name" binding:"required"`
	OwnerMobile       string  `json:"owner_mobile" binding:"required"`
	Address           string  `json:"address" binding:"required"`
	Email             string  `json:"email" binding:"required,email"`
	LicenseAttachment *string `json:"license_attachment"`
}

// GamenetUpdateRequest represents a gamenet update request
type GamenetUpdateRequest struct {
	Name              *string `json:"name"`
	OwnerName         *string `json:"owner_name"`
	OwnerMobile       *string `json:"owner_mobile"`
	Address           *string `json:"address"`
	Email             *string `json:"email"`
	Password          *string `json:"-"` // Hidden from JSON
	LicenseAttachment *string `json:"license_attachment"`
}

// GamenetResponse represents a gamenet response
type GamenetResponse struct {
	ID                int       `json:"id"`
	Name              string    `json:"name"`
	OwnerName         string    `json:"owner_name"`
	OwnerMobile       string    `json:"owner_mobile"`
	Address           string    `json:"address"`
	Email             string    `json:"email"`
	LicenseAttachment *string   `json:"license_attachment"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// GamenetSearchRequest represents a gamenet search request
type GamenetSearchRequest struct {
	Query    string `form:"query" json:"query"`
	Page     int    `form:"page" json:"page" binding:"min=1"`
	PageSize int    `form:"page_size" json:"page_size" binding:"min=1,max=100"`
}

// GamenetSearchResponse represents a paginated gamenet search response
type GamenetSearchResponse struct {
	Data       []GamenetResponse `json:"data"`
	Pagination PaginationInfo    `json:"pagination"`
}

// PaginationInfo represents pagination metadata
type PaginationInfo struct {
	CurrentPage int   `json:"current_page"`
	PageSize    int   `json:"page_size"`
	TotalItems  int64 `json:"total_items"`
	TotalPages  int   `json:"total_pages"`
	HasNext     bool  `json:"has_next"`
	HasPrev     bool  `json:"has_prev"`
}

// ToResponse converts Gamenet to GamenetResponse
func (g *Gamenet) ToResponse() GamenetResponse {
	return GamenetResponse{
		ID:                g.ID,
		Name:              g.Name,
		OwnerName:         g.OwnerName,
		OwnerMobile:       g.OwnerMobile,
		Address:           g.Address,
		Email:             g.Email,
		LicenseAttachment: g.LicenseAttachment,
		CreatedAt:         g.CreatedAt,
		UpdatedAt:         g.UpdatedAt,
	}
}
