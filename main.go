package main

import (
    "log"
    "github.com/lukeware-digital/mcp-code-review/server"
    "os"
)

func main() {
    // Obter API Key do OpenRouter
    apiKey := os.Getenv("OPENROUTER_API_KEY")
    if apiKey == "" {
        log.Fatal("OPENROUTER_API_KEY não encontrada nas variáveis de ambiente")
    }

    // Modelo padrão - pode ser customizado
    model := "minimax/minimax-m2:free"
    if customModel := os.Getenv("OPENROUTER_MODEL"); customModel != "" {
        model = customModel
    }

    // Criar e iniciar servidor
    mcpServer := server.NewServer(apiKey, model)
    mcpServer.Start()
}