package types

import (
	"vkyc-backend/models"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    *gin.H `json:"data"`
}

type OCRResponseData struct {
	OCRData *map[string]interface{} `json:"OCRdata"`
	MRZData *map[string]interface{} `json:"MRZdata"`
}

type OCRResponse struct {
	Status  string           `json:"Status"`
	Message string           `json:"Message"`
	Data    *OCRResponseData `json:"data"`
}

type FacematchResponse struct {
	Status bool                   `json:"status"`
	Data   map[string]interface{} `json:"data"`
}

type LivenessResponse struct {
	Status bool                    `json:"status"`
	Data   *map[string]interface{} `json:"data"`
}

type Pagination struct {
	PageSize    int `json:"pageSize"`
	CurrentPage int `json:"currentPage"`
	TotalPages  int `json:"totalPages"`
}

type OCRVerificationResponse struct {
	Document    models.Document
	OCRResponse OCRResponse
}
