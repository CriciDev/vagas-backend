## Context

O pacote `internal/opportunities` já implementa CRUD com leitura pública de oportunidades `published`. O método `List` do repositório monta filtros dinâmicos (`type`, `work_mode`, `location`) e retorna todas as linhas ordenadas por `created_at DESC, id DESC`. Não há paginação, e o handler devolve o array direto. A issue #31 pede paginação com metadados, deixando abertas as decisões de estilo e limites.

## Goals / Non-Goals

**Goals:**

- Paginar a listagem pública de oportunidades.
- Definir valores padrão e limite máximo previsíveis.
- Retornar metadados de paginação junto dos dados.
- Tratar entradas inválidas de forma previsível, sem erro.

**Non-Goals:**

- Implementar cursor/keyset pagination ou infinite scroll.
- Adicionar ordenação configurável.
- Paginar outros recursos (developers têm issue própria, #11).

## Decisions

- Usar `page` e `page_size` como query params. Alternativa considerada: `limit`/`offset`. `page`/`page_size` é mais amigável para o frontend e permite expor `total_pages` diretamente.
- `page_size` padrão de 20 e máximo de 100. Alternativa considerada: 10/50. 20/100 é um equilíbrio comum para APIs REST.
- Clampar valores inválidos em vez de rejeitar: `page` < 1 vira 1, `page_size` < 1 vira o padrão, `page_size` > 100 vira 100, valores não numéricos caem no padrão. Alternativa considerada: responder 400. Clampar atende ao critério de aceite "valores inválidos tratados de forma previsível" e simplifica o cliente.
- Responder com envelope `{ "data": [...], "meta": { "page", "page_size", "total", "total_pages" } }`. Alternativa considerada: manter o array e mover metadados para headers. O envelope é mais explícito para o frontend.
- Calcular `total` com uma query `count(*)` usando os mesmos filtros, e aplicar `LIMIT`/`OFFSET` na query de dados. `total_pages` é `ceil(total / page_size)`.

## Risks / Trade-offs

- Mudar o corpo da listagem de array para objeto é breaking change -> Mitigado porque o frontend ainda não consome o endpoint e a mudança está registrada nesta proposta.
- Uma query extra de `count` por listagem adiciona custo -> Aceitável no MVP; pode evoluir para contagem estimada se necessário.

## Migration Plan

- Ajustar o método `List` do repositório para receber paginação e retornar `([]Opportunity, total, error)`.
- Ajustar o service para montar o envelope com metadados.
- Ajustar o handler para ler e clampar `page`/`page_size`.
- Atualizar os testes de listagem para o novo formato.
- Rollback restaura a listagem em array e remove os campos de paginação.

## Open Questions

- Nenhuma. As decisões abertas da issue (`page`/`page_size` vs `limit`/`offset`, tamanho máximo) foram resolvidas nesta proposta.
