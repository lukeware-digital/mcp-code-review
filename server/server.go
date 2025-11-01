package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/lukeware-digital/mcp-code-review/openrouter"
	"github.com/lukeware-digital/mcp-code-review/types"
)

type MCPServer struct {
	OpenRouterClient *openrouter.Client
	Tools            []types.Tool
}

func NewServer(apiKey, model string) *MCPServer {
	openRouterClient := openrouter.NewClient(apiKey, model)

	tools := []types.Tool{
		{
			Name:        "code_review",
			Description: "Realiza code review detalhado de qualquer linguagem de programação",
			InputSchema: mustMarshalJSON(map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"code": map[string]interface{}{
						"type":        "string",
						"description": "O código a ser revisado",
					},
					"language": map[string]interface{}{
						"type":        "string",
						"description": "Linguagem de programação (ex: go, python, javascript, etc.)",
					},
				},
				"required": []string{"code", "language"},
			}),
		},
	}

	return &MCPServer{
		OpenRouterClient: openRouterClient,
		Tools:            tools,
	}
}

func (s *MCPServer) HandleMessage(message []byte) ([]byte, error) {
	var request types.MCPRequest
	if err := json.Unmarshal(message, &request); err != nil {
		return nil, fmt.Errorf("erro ao decodificar JSON: %v", err)
	}

	log.Printf("Método recebido: %s, ID: %v", request.Method, request.ID)

	// Se é uma notificação (sem ID), não envie resposta
	if request.ID == nil {
		log.Printf("Processando notificação: %s", request.Method)
		return s.handleNotification(request)
	}

	switch request.Method {
	case "initialize":
		return s.handleInitialize(request)
	case "tools/list":
		return s.handleToolsList(request)
	case "tools/call":
		return s.handleToolsCall(request)
	case "shutdown":
		return s.handleShutdown(request)
	case "ping":
		return s.handlePing(request)
	default:
		return s.createErrorResponse(request.ID, -32601, "Method not found", "Método não encontrado: "+request.Method)
	}
}

func (s *MCPServer) handleNotification(request types.MCPRequest) ([]byte, error) {
	switch request.Method {
	case "notifications/initialized":
		log.Println("Notificação de inicialização recebida - nenhuma resposta necessária")
		return nil, nil
	case "notifications/cancelled":
		log.Println("Notificação de cancelamento recebida - nenhuma resposta necessária")
		return nil, nil
	default:
		log.Printf("Notificação não reconhecida: %s", request.Method)
		return nil, nil
	}
}

func (s *MCPServer) handleInitialize(request types.MCPRequest) ([]byte, error) {
	var initParams types.InitializeParams

	if len(request.Params) > 0 {
		if err := json.Unmarshal(request.Params, &initParams); err != nil {
			return s.createErrorResponse(request.ID, -32700, "Parse error", err.Error())
		}
	}

	protocolVersion := "2024-11-07"
	if initParams.ProtocolVersion != "" {
		protocolVersion = initParams.ProtocolVersion
	}

	result := types.InitializeResult{
		ProtocolVersion: protocolVersion,
		Capabilities: mustMarshalJSON(map[string]interface{}{
			"tools": map[string]interface{}{},
			// Removemos resources e prompts já que não os implementamos
		}),
		ServerInfo: types.ServerInfo{
			Name:    "mcp-code-review",
			Version: "1.0.0",
		},
	}

	return s.createSuccessResponse(request.ID, result)
}

func (s *MCPServer) handleToolsList(request types.MCPRequest) ([]byte, error) {
	result := map[string]interface{}{
		"tools": s.Tools,
	}
	return s.createSuccessResponse(request.ID, result)
}

func (s *MCPServer) handleToolsCall(request types.MCPRequest) ([]byte, error) {
	var toolCallParams types.ToolCallParams

	if err := json.Unmarshal(request.Params, &toolCallParams); err != nil {
		return s.createErrorResponse(request.ID, -32700, "Parse error", err.Error())
	}

	if toolCallParams.Name != "code_review" {
		return s.createErrorResponse(request.ID, -32602, "Invalid params", "Ferramenta não encontrada: "+toolCallParams.Name)
	}

	var args struct {
		Code     string `json:"code"`
		Language string `json:"language"`
	}
	if err := json.Unmarshal(toolCallParams.Arguments, &args); err != nil {
		return s.createErrorResponse(request.ID, -32602, "Invalid params", "Falha ao decodificar argumentos: "+err.Error())
	}

	if args.Code == "" || args.Language == "" {
		return s.createErrorResponse(request.ID, -32602, "Invalid params", "Parâmetros 'code' e 'language' são obrigatórios")
	}

	log.Printf("Iniciando code review para linguagem: %s", args.Language)
	// Realiza o code review via OpenRouter
	review, err := s.OpenRouterClient.PerformCodeReview(args.Code, args.Language)
	if err != nil {
		log.Printf("Erro no code review: %v", err)
		return s.createErrorResponse(request.ID, -32000, "Internal error",
			fmt.Sprintf("Erro ao realizar code review: %v", err))
	}

	result := types.ToolsCallResult{
		Content: []types.ContentItem{{Type: "text", Text: review}},
	}

	log.Printf("Code review concluído com sucesso")
	return s.createSuccessResponse(request.ID, result)
}

func (s *MCPServer) handleShutdown(request types.MCPRequest) ([]byte, error) {
	log.Println("Recebido comando de shutdown")
	os.Exit(0)
	return nil, nil
}

func (s *MCPServer) handlePing(request types.MCPRequest) ([]byte, error) {
	return s.createSuccessResponse(request.ID, "pong")
}

func (s *MCPServer) createSuccessResponse(id interface{}, result interface{}) ([]byte, error) {
	response := types.MCPResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	}
	return json.Marshal(response)
}

func (s *MCPServer) createErrorResponse(id interface{}, code int, message, data string) ([]byte, error) {
	response := types.MCPResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error: &types.MCPError{
			Code:    code,
			Message: message,
			Data:    data,
		},
	}
	return json.Marshal(response)
}

func (s *MCPServer) Start() {
	log.SetOutput(os.Stderr)
	log.Println("MCP Code Review Server iniciado (stderr)")
	log.Println("Aguardando mensagens JSON-RPC...")

	scanner := bufio.NewScanner(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		log.Printf("Recebido: %s", string(line))

		response, err := s.HandleMessage(line)
		if err != nil {
			log.Printf("Erro ao processar mensagem: %v", err)
			continue
		}

		// Só envia resposta se não for nil (notificações retornam nil)
		if response != nil {
			writer.Write(response)
			writer.WriteByte('\n')
			if err := writer.Flush(); err != nil {
				log.Printf("Erro ao enviar resposta: %v", err)
			} else {
				log.Printf("Resposta enviada com sucesso")
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Erro no scanner: %v", err)
	}
}

func mustMarshalJSON(v interface{}) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}
