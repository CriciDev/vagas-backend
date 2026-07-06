## Why

A listagem pública de oportunidades (`GET /api/opportunities`) retorna hoje todas as oportunidades publicadas de uma vez, em um array JSON simples. Conforme a base cresce, a resposta fica grande e cara, e o frontend Angular não tem como pedir páginas. A issue #31 pede paginação com metadados para preparar a integração com o front.

## What Changes

- Adicionar paginação por `page` e `page_size` na listagem pública de oportunidades.
- Aplicar `page_size` padrão de 20 e máximo de 100.
- Clampar valores inválidos (`page` < 1, `page_size` fora do intervalo, valores não numéricos) para o intervalo válido, sem retornar erro.
- Alterar a resposta da listagem de um array puro para um envelope `{ "data": [...], "meta": {...} }`, com `meta` contendo `page`, `page_size`, `total` e `total_pages`.
- Fazer o repositório retornar o total de registros (via `count`) além da página atual.
- Adicionar testes para paginação padrão, customizada, valores inválidos e metadados.

## Capabilities

### New Capabilities

- None.

### Modified Capabilities

- `opportunity-management`: a listagem pública passa a ser paginada e a devolver dados e metadados de paginação.

## Impact

- Mudança de contrato na listagem: o corpo de `GET /api/opportunities` passa de array para objeto `{ data, meta }`. Como o frontend ainda não consome esse endpoint, o impacto é aceitável no MVP.
- Ajusta `internal/opportunities` (model, usecase, repository, controller) e os testes existentes de listagem.
- Sem novas dependências.
- Atende a issue https://github.com/CriciDev/vagas-backend/issues/31.
