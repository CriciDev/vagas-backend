## Context

O backend usa Gin com `gin.Default()`, que já inclui os middlewares de Logger e Recovery. Nenhum identificador por requisição é gerado hoje, então logs de acesso e respostas não podem ser correlacionados. A issue #43 pede um `request id` simples, sem tracing distribuído nem serviço externo.

## Goals / Non-Goals

**Goals:**

- Gerar um `request id` para toda requisição que chega à API.
- Reusar um id válido enviado pelo cliente no header `X-Request-ID`.
- Devolver o id no header da resposta.
- Tornar o id acessível no contexto da requisição e nas linhas de log de acesso.

**Non-Goals:**

- Integrar OpenTelemetry, spans ou tracing distribuído.
- Usar serviço externo ou armazenamento do id.
- Propagar o id para chamadas de saída (a API não faz chamadas externas hoje).

## Decisions

- Usar o header `X-Request-ID` para entrada e saída. Alternativa considerada: `X-Correlation-ID`. `X-Request-ID` é o padrão de fato mais comum em proxies e bibliotecas.
- Aceitar o id do cliente quando válido e gerar um novo quando ausente ou inválido. Alternativa considerada: sempre gerar um id novo. Reusar o id do cliente facilita correlação ponta a ponta quando um gateway já define um.
- Validar o id recebido antes de confiar nele: aceitar somente uma string não vazia, com no máximo 128 caracteres e caracteres seguros (`A-Z a-z 0-9 . _ -`). Isso evita log/response injection a partir de um header arbitrário. Ids inválidos são descartados e um novo é gerado.
- Gerar o id com `crypto/rand` (16 bytes em hexadecimal) em vez de adicionar uma dependência de UUID, mantendo o módulo seco conforme `AGENTS.md`.
- Colocar o middleware em `internal/middleware` por ser um comportamento HTTP transversal, e não específico de `auth` ou `devs`.
- Substituir `gin.Default()` por `gin.New()` com a ordem explícita `RequestID` → `Logger` → `Recovery`, para que o logger de acesso já enxergue o id pelo `Keys` do contexto.

## Risks / Trade-offs

- Confiar no header do cliente pode permitir ids forjados -> Mitigado com validação de tamanho e charset; um id inválido é substituído por um gerado.
- Trocar `gin.Default()` por `gin.New()` muda a montagem do router -> Mitigado mantendo os mesmos middlewares (Logger e Recovery), apenas com ordem explícita.

## Migration Plan

- Adicionar o pacote `internal/middleware` apenas com o middleware de request id.
- Registrar o middleware em `cmd/server/main.go` antes das rotas.
- Rollback remove o pacote `internal/middleware` e restaura `gin.Default()` sem afetar `auth` ou `developers`.

## Open Questions

- O header deve ser `X-Request-ID` ou `X-Correlation-ID`? (Proposto: `X-Request-ID`.)
- A API deve aceitar o id enviado pelo cliente ou sempre gerar um novo? (Proposto: aceitar se válido, senão gerar.)
