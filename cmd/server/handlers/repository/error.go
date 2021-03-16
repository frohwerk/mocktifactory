package repository

type errorResponse struct {
	Status  int    `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
}
