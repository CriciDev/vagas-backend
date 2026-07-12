## Why

O README define oportunidades como parte central do MVP, mas o backend atual possui contrato real apenas para developers e autenticação básica. Precisamos registrar a primeira versão de `opportunities` antes de implementar endpoints, banco e regras de publicação, para permitir debate de produto e reduzir retrabalho.

## What Changes

- Definir o modelo mínimo de uma `opportunity` para o MVP.
- Sugerir campos obrigatórios e opcionais usando termos de domínio em inglês.
- Definir quais operações de leitura serão públicas e quais operações de escrita serão protegidas.
- Registrar cenários esperados para criação, listagem, busca por id, edição e remoção.
- Registrar decisões abertas para validação com a comunidade antes ou durante a implementação.
- Não implementar código, banco de dados ou endpoints nesta mudança.

## Capabilities

### New Capabilities

- `opportunity-management`: contrato inicial para criar, listar, buscar, atualizar, remover e validar opportunities do MVP.

### Modified Capabilities

- None.

## Impact

- Cria artefatos OpenSpec para orientar a futura implementação de `opportunities`.
- Propõe futuro impacto em API HTTP sob `/api/opportunities`, persistência PostgreSQL e pacote Go dedicado em `internal/opportunities`.
- Mantém compatibilidade com as specs existentes de `developer-profiles` e `basic-auth-roles`, usando leitura pública e escrita protegida por role `admin` no primeiro momento.
- Atende a issue https://github.com/CriciDev/vagas-backend/issues/12 como proposta inicial, deixando decisões de produto explicitamente abertas para discussão.
