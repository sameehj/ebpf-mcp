package resources

import (
	"context"
	"errors"

	"github.com/mark3labs/mcp-go/mcp"
)

func Register(mcpServer *mcp.Server) {
	mcpServer.Register("resources.list", listHandler)
	mcpServer.Register("resources.get", getHandler)
	mcpServer.Register("resources.create", createHandler)
}

func listHandler(ctx context.Context, req mcp.Request) (any, error) {
	return mcp.ListResourcesResult{Resources: List()}, nil
}

func getHandler(ctx context.Context, req mcp.Request) (any, error) {
	var params mcp.GetResourceParams
	if err := req.UnmarshalParams(&params); err != nil {
		return nil, err
	}

	content, ok := Get(params.ResourceID)
	if !ok {
		return nil, errors.New("resource not found")
	}
	return mcp.GetResourceResult{Content: []mcp.Content{content}}, nil
}

func createHandler(ctx context.Context, req mcp.Request) (any, error) {
	var params mcp.CreateResourceParams
	if err := req.UnmarshalParams(&params); err != nil {
		return nil, err
	}

	if params.ResourceID == "" || len(params.Content) == 0 {
		return nil, errors.New("missing resource_id or content")
	}

	Put(params.ResourceID, params.Type, params.Content[0])
	return mcp.CreateResourceResult{OK: true}, nil
}
