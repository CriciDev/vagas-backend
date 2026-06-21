## Why

O projeto precisa sair do scaffold inicial para uma API mínima executável que permita listar e manter perfis de developers da comunidade. Este é o primeiro passo pequeno para validar a stack Go com Gin, facilitar contribuição local via Docker e preparar a base de autenticação com permissões e roles.

## What Changes

- Adicionar Gin como roteador HTTP principal do backend.
- Implementar um CRUD simplificado de developers sobre o scaffold existente em `internal/devs`.
- Criar contratos HTTP mínimos para criar, listar, buscar, atualizar e remover developers.
- Adicionar uma base simples de autenticação e autorização por roles para proteger operações de escrita.
- Configurar ambiente Docker local para que contribuidores consigam subir API e banco de dados com comandos simples.
- Documentar como executar o backend localmente com Docker.

## Capabilities

### New Capabilities

- `developer-profiles`: CRUD simplificado de perfis de developers da comunidade.
- `basic-auth-roles`: autenticação básica e autorização por roles para proteger rotas conforme o perfil do usuário.
- `local-docker-environment`: ambiente local com Docker para subir API e dependências do backend.

### Modified Capabilities

- None.

## Impact

- `go.mod` e `go.sum` receberão dependências do Gin e bibliotecas necessárias para a base HTTP/autenticação.
- `cmd/server/main.go` passará a inicializar o servidor HTTP e registrar rotas.
- `internal/devs` será preenchido com modelo, controller, usecase e repository mínimos.
- `internal/auth` será preenchido com contratos e middleware mínimos de autenticação/autorização.
- `infra/docker-compose.yml` e arquivos Docker relacionados serão ajustados para execução local da API e PostgreSQL.
- `README.md` receberá instruções mínimas de execução local.
