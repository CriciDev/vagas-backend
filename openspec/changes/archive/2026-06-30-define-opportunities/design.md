## Context

O backend já possui Gin, PostgreSQL, autenticação bearer com role `admin` e CRUD de `developers`. O README coloca oportunidades como elemento central do MVP, mas ainda não existe contrato para o domínio de `opportunities`.

Esta mudança é deliberadamente documental: ela define o primeiro contrato de produto e API para oportunidades, sem implementar banco, rotas ou código Go. A implementação futura deve seguir o padrão simples já usado em `internal/devs`, evitando camadas genéricas ou diretórios vazios.

## Goals / Non-Goals

**Goals:**

- Definir o modelo inicial de `opportunity` com campos obrigatórios e opcionais.
- Definir operações mínimas de CRUD para `opportunities`.
- Definir leitura pública e escrita protegida para o primeiro MVP.
- Definir status inicial para controlar visibilidade pública.
- Deixar decisões de produto abertas de forma explícita para validação com a comunidade.

**Non-Goals:**

- Implementar endpoints, schema SQL, migrations ou código Go nesta mudança.
- Modelar contas de organizações ou vínculo obrigatório com uma empresa cadastrada.
- Implementar auto-expiration, workflow de aprovação ou moderação avançada.
- Implementar busca full-text, paginação avançada ou filtros complexos.
- Definir monetização, candidatura dentro da plataforma ou tracking de applicants.

## Decisions

- Usar `opportunity` como termo de domínio principal. Alternativas consideradas: `job`, `position` e `opening`. `Opportunity` cobre vagas, projetos, freelance, voluntariado e mentoria sem limitar o MVP a emprego formal.
- Propor campos obrigatórios mínimos: `title`, `description`, `organization_name`, `type`, `work_mode` e pelo menos um canal de contato (`contact_email` ou `contact_url`). Alternativas consideradas: exigir localização, senioridade e remuneração desde o início. Esses campos são úteis, mas podem bloquear publicações legítimas no MVP.
- Representar organização inicialmente por campos textuais (`organization_name` e `organization_url`) em vez de uma entidade relacional. Alternativas consideradas: criar `organizations` agora. A entidade própria deve esperar validação de ownership e fluxo de publicação.
- Usar `status` para visibilidade, com valores iniciais `draft`, `published`, `closed` e `archived`. Alternativas consideradas: usar apenas delete físico ou um booleano `active`. Status deixa o contrato mais claro para listagem pública e evolução futura.
- Expor leitura pública apenas de `published opportunities`. Alternativas consideradas: listar tudo publicamente ou exigir autenticação. O produto precisa facilitar descoberta pública, mas não deve expor rascunhos ou itens arquivados.
- Restringir criação, atualização e remoção a `admin` no primeiro momento. Alternativas consideradas: organizações autenticadas publicarem diretamente. O backend ainda não possui contas de organizações nem ownership, então `admin` reduz risco operacional.
- Planejar endpoints futuros em `/api/opportunities`, espelhando o padrão de `/api/developers`. Alternativas consideradas: `/api/jobs`. `opportunities` mantém o vocabulário mais amplo escolhido para o domínio.

## Risks / Trade-offs

- Admin-only publication limits self-service by companies -> Mitigar deixando ownership de organizações como decisão aberta e próxima evolução natural.
- Textual organization fields can create duplicate organization names -> Mitigar aceitando duplicidade no MVP e registrando futura capability de organizations se necessário.
- Public listing only for `published` items may hide useful drafts during testing -> Mitigar com dados seed ou endpoints admin futuros se a implementação precisar.
- No auto-expiration can leave stale opportunities online -> Mitigar com `status = closed` manual no MVP e discutir `expires_at` antes da implementação.
- Broad `type` values can become inconsistent -> Mitigar documentando enum inicial e validando payloads na futura implementação.

## Migration Plan

- A implementação futura deve criar pacote `internal/opportunities` somente quando houver comportamento real.
- A implementação futura deve adicionar tabela `opportunities` ao bootstrap PostgreSQL local.
- A implementação futura deve registrar rotas públicas de leitura e rotas protegidas de escrita em `cmd/server/main.go`.
- Rollback da implementação futura deve remover endpoints, tabela local e pacote `internal/opportunities` sem afetar `developers` ou `auth`.

## Open Questions

- `opportunities` devem ser sempre vinculadas a uma entidade `organization` cadastrada?
- O MVP deve aceitar todos os tipos propostos (`full_time`, `part_time`, `contract`, `freelance`, `volunteer`, `project`, `mentorship`) ou começar com uma lista menor?
- Quem poderá publicar depois do MVP: somente `admin`, organizações autenticadas ou qualquer usuário validado?
- Deve existir `expires_at` obrigatório ou expiração automática por padrão?
- Campos como `salary_range`, `seniority` e `skills` devem ser obrigatórios, opcionais ou deixados para uma segunda versão?
