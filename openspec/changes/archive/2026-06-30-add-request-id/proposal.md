## Why

Quando algo falha em produção ou em ambiente compartilhado, é difícil correlacionar uma resposta de erro com as linhas de log que a originaram. A API atual usa o logger padrão do Gin, que não expõe nenhum identificador por requisição. Um `request id` simples permite que suporte e debug liguem um cliente, uma resposta e o log correspondente sem tracing distribuído.

## What Changes

- Adicionar um middleware de `request id` aplicado a todas as rotas da API.
- Reusar o id enviado pelo cliente no header `X-Request-ID` quando ele for válido, e gerar um novo id quando estiver ausente ou inválido.
- Retornar o id no header `X-Request-ID` de toda resposta.
- Disponibilizar o id no contexto da requisição para que handlers e logs usem o mesmo identificador.
- Incluir o id nas linhas de log de acesso da API.
- Adicionar teste cobrindo geração automática e reuso de id válido do cliente.

## Capabilities

### New Capabilities

- `request-correlation`: identificação por requisição via header `X-Request-ID`, propagada para resposta, contexto e logs.

### Modified Capabilities

- None.

## Impact

- Adiciona o pacote `internal/middleware` com comportamento real (sem placeholders).
- Ajusta `cmd/server/main.go` para registrar o middleware antes do logger e do recovery.
- Não altera contratos existentes de `developers` ou `auth`; apenas acrescenta um header de resposta.
- Sem novas dependências: o id é gerado com `crypto/rand` da biblioteca padrão.
- Atende a issue https://github.com/CriciDev/vagas-backend/issues/43.
