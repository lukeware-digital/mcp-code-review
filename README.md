# MCP config

```json
{
  "mcpServers": {
    "mcp-code-review": {
      "command": "/caminho/para/mcp-code-review",
      "env": {
        "OPENROUTER_API_KEY": "sua-chave-api-openrouter"
      }
    }
  }
}
```

# Model User
```
minimax/minimax-m2:free
```

# Test

```shell
echo '{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"code_review","arguments":{"code":"function test() { return true; }","language":"javascript"}}}' | OPENROUTER_API_KEY="sk-or-v1-831db936471b8bec9198bc05f9b8954055a8e6888ef8c9da4a5ae48247dec28..." ./mcp-code-review
```