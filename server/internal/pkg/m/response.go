package m

type Response struct {
	Status string `json:"status"`
	Msg    string `json:"msg"`
}

const (
	Success           = "success"
	ErrFailed         = "failed"
	ErrInvalidParams  = "invalid_params"
	ErrInvalidRequest = "invalid_request"
)
