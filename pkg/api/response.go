package api

type ImageSuccessResponse struct {
	ProcessingTime int64 `json:"time"`
	Cached         bool  `json:"cached"`
}

type ImageErrorResponse struct {
	Error string `json:"error"`
}
