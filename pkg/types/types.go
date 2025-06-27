package types

type RPCRequest struct {
	JSONRPC string                 `json:"jsonrpc"`
	Method  string                 `json:"method"`
	Params  map[string]interface{} `json:"params"`
	ID      interface{}            `json:"id"`
}

type RPCResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *RPCError   `json:"error,omitempty"`
}

type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewSuccessResponse(id interface{}, result interface{}) RPCResponse {
	return RPCResponse{JSONRPC: "2.0", ID: id, Result: result}
}

// Error with default code -32600 (Invalid Request)
func NewErrorResponse(id interface{}, message string) RPCResponse {
	return NewErrorResponseWithCode(id, -32600, message)
}

// Error with custom code
func NewErrorResponseWithCode(id interface{}, code int, message string) RPCResponse {
	return RPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error: &RPCError{
			Code:    code,
			Message: message,
		},
	}
}

type Tool struct {
	ID           string                                                `json:"id"`
	Title        string                                                `json:"title"`
	Description  string                                                `json:"description"`
	InputSchema  map[string]interface{}                                `json:"input_schema"`
	OutputSchema map[string]interface{}                                `json:"output_schema"`
	Annotations  map[string]interface{}                                `json:"annotations,omitempty"`
	Call         func(map[string]interface{}) (interface{}, error)     `json:"-"`
	Stream       func(map[string]interface{}, func(interface{})) error `json:"-"` // 👈 add this
}

func (t Tool) Metadata() ToolMetadata {
	return ToolMetadata{
		Name:        t.ID,
		Description: t.Description,
	}
}

type ToolMetadata struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Param struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Required bool   `json:"required"`
}

type ErrorDetail struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

func ValidationError(msg string) *ErrorDetail {
	return &ErrorDetail{
		Type:    "VALIDATION_ERROR",
		Message: msg,
	}
}
