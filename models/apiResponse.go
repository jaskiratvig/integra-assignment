package models

// ApiResponse - Standard API response format
type ApiResponse struct {
    Success bool        `json:"success"`
    Message string      `json:"message"`
    Data    interface{} `json:"data"`
}