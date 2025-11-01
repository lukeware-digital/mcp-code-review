package types

// Estruturas para o protocolo MCP
type MCPServerConfig struct {
    OpenRouterAPIKey string `json:"openrouter_api_key"`
    Model            string `json:"model"`
}

type MCPRequest struct {
    JSONRPC string      `json:"jsonrpc"`
    ID      interface{} `json:"id"`
    Method  string      `json:"method"`
    Params  MCPParams   `json:"params,omitempty"`
}

type MCPParams struct {
    Arguments interface{} `json:"arguments,omitempty"`
    Name      string      `json:"name,omitempty"`
}

type MCPResponse struct {
    JSONRPC string      `json:"jsonrpc"`
    ID      interface{} `json:"id"`
    Result  interface{} `json:"result,omitempty"`
    Error   *MCPError   `json:"error,omitempty"`
}

type MCPError struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}

type Tool struct {
    Name        string                 `json:"name"`
    Description string                 `json:"description"`
    InputSchema map[string]interface{} `json:"inputSchema"`
}

type InitializeResult struct {
    ProtocolVersion string                 `json:"protocolVersion"`
    Capabilities    map[string]interface{} `json:"capabilities"`
    ServerInfo      ServerInfo             `json:"serverInfo"`
}

type ServerInfo struct {
    Name    string `json:"name"`
    Version string `json:"version"`
}

type ToolsCallResult struct {
    Content []ContentItem `json:"content"`
}

type ContentItem struct {
    Type string `json:"type"`
    Text string `json:"text"`
}

// Estruturas para OpenRouter
type OpenRouterRequest struct {
    Model     string               `json:"model"`
    Messages  []OpenRouterMessage  `json:"messages"`
    MaxTokens int                  `json:"max_tokens,omitempty"`
}

type OpenRouterMessage struct {
    Role    string `json:"role"`
    Content string `json:"content"`
}

type OpenRouterResponse struct {
    Choices []OpenRouterChoice `json:"choices"`
    Error   *OpenRouterError   `json:"error,omitempty"`
}

type OpenRouterChoice struct {
    Message OpenRouterMessage `json:"message"`
}

type OpenRouterError struct {
    Message string `json:"message"`
}