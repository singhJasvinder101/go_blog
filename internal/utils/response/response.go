package response

import "github.com/singhJasvinder101/go_blog/internal/types"


const (
	StatusOK    = "ok"
	StatusError = "error"
)


func ErrorResponse(err error) types.APIResponse {
	return types.APIResponse{
		Status: StatusError,
		Error:  err.Error(),
	}
}


