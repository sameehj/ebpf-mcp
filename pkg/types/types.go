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

func NewErrorResponse(id interface{}, message string) RPCResponse {
	return RPCResponse{
		JSONRPC: "2.0", ID: id,
		Error: &RPCError{Code: -32600, Message: message},
	}
}

type Tool struct {
	ID          string                                                  `json:"id"`
	Title       string                                                  `json:"title"`
	Description string                                                  `json:"description"`
	Parameters  []Param                                                 `json:"parameters"`
	Call        func(input map[string]interface{}) (interface{}, error) `json:"-"`
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
