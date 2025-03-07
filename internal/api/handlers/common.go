package handlers

type ErrorResponse struct {
	Error string `json:"errors"`
}

var NotFoundResponse = ErrorResponse{Error: "not found"}
var AlreadyExistsResponse = ErrorResponse{Error: "already exists"}
var InvalidRequestResponse = ErrorResponse{Error: "invalid request"}
var ServerErrorResponse = ErrorResponse{Error: "server error"}
var BadGatewayResponse = ErrorResponse{Error: "external service error"}
