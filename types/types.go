package types

import (
    "encoding/json"
    "fmt"
)

// JSON-RPC version type and constant
// Padroniza índice JSONRPC
//
type JSONRPC string

const JSONRPCVersion JSONRPC = "2.0"

// MCPServerConfig não possui mais a chave da API, evitando vazamento
// A chave deve ser lida de env nos pontos de uso sensíveis.
type MCPServerConfig struct {
    Model string `json:"model"`
}

// Estruturas para o protocolo MCP
type MCPRequest struct {
    JSONRPC JSONRPC        `json:"jsonrpc"`
    ID      interface{}    `json:"id,omitempty"`
    Method  string         `json:"method"`
    Params  json.RawMessage `json:"params,omitempty"`
}

type MCPResponse struct {
    JSONRPC JSONRPC        `json:"jsonrpc"`
    ID      interface{}    `json:"id,omitempty"`
    Result  interface{}    `json:"result,omitempty"`
    Error   *MCPError      `json:"error,omitempty"`
}

type MCPError struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}

func (e MCPError) Error() string {
    return fmt.Sprintf("MCP error %d: %s", e.Code, e.Message)
}

func NewMCPInvalidParams(msg string, data interface{}) MCPError {
    return MCPError{Code: -32602, Message: msg, Data: data}
}

func NewMCPInvalidRequest(msg string) MCPError {
    return MCPError{Code: -32600, Message: msg}
}

// Generics replaced: prefer RawMessage for schemas e maior robustez
type Tool struct {
    Name        string          `json:"name"`
    Description string          `json:"description"`
    InputSchema json.RawMessage `json:"inputSchema"`
}

type InitializeParams struct {
    ProtocolVersion string          `json:"protocolVersion"`
    Capabilities    json.RawMessage `json:"capabilities"`
    ClientInfo      json.RawMessage `json:"clientInfo,omitempty"`
}

type ToolCallParams struct {
    Name      string          `json:"name"`
    Arguments json.RawMessage `json:"arguments"`
}

type InitializeResult struct {
    ProtocolVersion string              `json:"protocolVersion"`
    Capabilities    json.RawMessage     `json:"capabilities"`
    ServerInfo      ServerInfo          `json:"serverInfo"`
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
    Model       string               `json:"model"`
    Messages    []OpenRouterMessage  `json:"messages"`
    MaxTokens   *int                 `json:"max_tokens,omitempty"`
    TopP        *float64             `json:"top_p,omitempty"`
    Temperature *float64             `json:"temperature,omitempty"`
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
	