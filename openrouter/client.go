package openrouter

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "github.com/lukeware-digital/mcp-code-review/types"
)

type Client struct {
    APIKey     string
    BaseURL    string
    HTTPClient *http.Client
    Model      string
}

func NewClient(apiKey, model string) *Client {
    return &Client{
        APIKey:     apiKey,
        BaseURL:    "https://openrouter.ai/api/v1",
        HTTPClient: &http.Client{},
        Model:      model,
    }
}

func (c *Client) PerformCodeReview(code string, language string) (string, error) {
    prompt := c.buildCodeReviewPrompt(code, language)
    
    requestBody := types.OpenRouterRequest{
        Model: c.Model,
        Messages: []types.OpenRouterMessage{
            {
                Role:    "user",
                Content: prompt,
            },
        },
        MaxTokens: 4000,
    }

    jsonData, err := json.Marshal(requestBody)
    if err != nil {
        return "", fmt.Errorf("erro ao serializar requisição: %v", err)
    }

    req, err := http.NewRequest("POST", c.BaseURL+"/chat/completions", bytes.NewBuffer(jsonData))
    if err != nil {
        return "", fmt.Errorf("erro ao criar requisição: %v", err)
    }

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+c.APIKey)
    req.Header.Set("HTTP-Referer", "https://github.com/mcp-code-review")
    req.Header.Set("X-Title", "MCP-Code-Review")

    resp, err := c.HTTPClient.Do(req)
    if err != nil {
        return "", fmt.Errorf("erro na requisição HTTP: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return "", fmt.Errorf("erro da API OpenRouter: status %d", resp.StatusCode)
    }

    var openRouterResp types.OpenRouterResponse
    if err := json.NewDecoder(resp.Body).Decode(&openRouterResp); err != nil {
        return "", fmt.Errorf("erro ao decodificar resposta: %v", err)
    }

    if len(openRouterResp.Choices) == 0 {
        return "", fmt.Errorf("nenhuma resposta da API")
    }

    return openRouterResp.Choices[0].Message.Content, nil
}

func (c *Client) buildCodeReviewPrompt(code, language string) string {
    return fmt.Sprintf(`Por favor, faça um code review detalhado do seguinte código em %s:

%s

Forneça um análise abrangente cobrindo:

1. **QUALIDADE DO CÓDIGO**
   - Legibilidade e clareza
   - Estrutura e organização
   - Nomenclatura de variáveis/funções
   - Complexidade ciclomática

2. **BOAS PRÁTICAS**
   - Princípios SOLID (se aplicável)
   - Padrões de design
   - Reutilização de código
   - Coesão e acoplamento

3. **PERFORMANCE**
   - Complexidade algorítmica
   - Uso eficiente de recursos
   - Possíveis gargalos

4. **SEGURANÇA**
   - Vulnerabilidades potenciais
   - Validação de entrada
   - Tratamento de erros

5. **MELHORIAS SUGERIDAS**
   - Refatorações específicas
   - Sugestões de otimização
   - Alternativas mais eficientes

Forneça exemplos concretos de melhorias quando aplicável. Seja direto e construtivo.`, language, code)
}