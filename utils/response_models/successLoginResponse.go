package response_models

type SuccessLoginResponse struct {
	StatusCode int `json:"status_code"`
	Message string `json:"message"`
	Token string `json:"token"`
	User interface{} `json:"user"`
}
