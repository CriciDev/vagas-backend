## Context

O backend está em um scaffold inicial Go com pacotes `internal/devs`, `internal/auth`, `internal/jobs`, `internal/companies` e `cmd/server/main.go` vazio. O repositório já possui `infra/docker-compose.yml` com PostgreSQL, mas ainda não possui API executável, Dockerfile da aplicação ou instruções completas para contribuição local.

Esta mudança deve transformar o scaffold em uma API mínima de comunidade, priorizando simplicidade, baixo atrito para novos contribuidores e uma base clara para evoluções futuras.

## Goals / Non-Goals

**Goals:**

- Inicializar um servidor HTTP com Gin em `cmd/server/main.go`.
- Implementar CRUD simplificado de developers com persistência em PostgreSQL.
- Expor leitura pública de developers e proteger escrita com autenticação bearer e role `admin`.
- Criar uma base reaproveitável de autenticação/autorização em `internal/auth`.
- Permitir execução local via Docker Compose com API e PostgreSQL.
- Atualizar o README com comandos mínimos para subir e validar a aplicação.

**Non-Goals:**

- Implementar cadastro público completo de usuários.
- Implementar recuperação de senha, confirmação de e-mail ou OAuth/social login.
- Implementar ownership entre usuário autenticado e perfil de developer.
- Implementar paginação avançada, filtros complexos ou upload de avatar.
- Criar uma arquitetura completa de migrations além do necessário para bootstrap local.

## Decisions

- Usar Gin como roteador HTTP principal. Alternativas consideradas: `net/http` puro ou outro framework. Gin atende ao pedido explícito, reduz boilerplate para rotas/middlewares e é familiar para novos contribuidores Go.
- Manter organização por feature em `internal/<feature>`. Alternativas consideradas: criar camadas globais `controllers`, `services` e `repositories`. A estrutura atual já aponta para pacotes por domínio e evita uma reorganização prematura.
- Persistir developers em PostgreSQL desde o primeiro CRUD. Alternativas consideradas: armazenamento em memória para acelerar o primeiro passo. O projeto já possui PostgreSQL no Compose e a persistência real evita retrabalho logo na primeira contribuição útil.
- Proteger escrita de developers apenas com role `admin` neste primeiro passo. Alternativas consideradas: permitir que cada developer gerencie o próprio perfil. Ownership exige modelagem adicional de usuários e vínculo com perfis, então fica fora do baby step.
- Usar token bearer assinado para autenticação. Alternativas consideradas: basic auth em todas as requisições ou sessão server-side. Token bearer se integra melhor com frontend Angular e permite middleware simples por role.
- Usar variáveis de ambiente para configuração. Alternativas consideradas: arquivo de configuração versionado. Ambiente funciona melhor com Docker Compose, CI e deploy futuro.

## Risks / Trade-offs

- Autenticação inicial simples pode não cobrir fluxos reais de usuários → manter escopo documentado e isolar auth em `internal/auth` para evolução posterior.
- Proteger escrita somente por `admin` limita colaboração direta de developers → aceitar como restrição do primeiro passo e planejar ownership em mudança futura.
- Bootstrap local de banco pode divergir de migrations formais futuras → manter SQL inicial pequeno e fácil de migrar para ferramenta dedicada depois.
- Gin introduz dependência externa no scaffold → manter uso direto e idiomático, sem wrappers genéricos prematuros.

## Migration Plan

- Adicionar dependências Go necessárias.
- Criar schema inicial para `developers` e usuário admin local.
- Atualizar Docker Compose para subir PostgreSQL e API.
- Documentar comandos no README.
- Validar localmente com `go test ./...` e execução via Docker Compose.

Rollback consiste em remover as novas dependências, Dockerfile/configurações da API e os arquivos preenchidos em `internal/devs` e `internal/auth`.

## Open Questions

- Qual será o fluxo futuro para developers criarem e manterem seus próprios perfis sem intervenção de admin?
- Quais campos serão obrigatórios no perfil público de developer após validação com a comunidade?
