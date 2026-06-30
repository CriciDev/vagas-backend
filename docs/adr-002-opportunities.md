# ADR 002: Criar opportunities como recurso próprio

## Status

Aceito

## Contexto

O MVP do projeto precisa divulgar oportunidades para a comunidade de tecnologia de Criciúma.

O backend já possui autenticação básica, controle de role `admin` e CRUD de `developers`.

As oportunidades têm regras próprias de publicação, visibilidade pública, contato e categorização.

## Decisão

Criar `opportunities` como recurso próprio da API.

A API terá rotas sob `/api/opportunities`.

Leituras públicas listarão e buscarão apenas oportunidades com status `published`.

Criação, atualização e remoção serão permitidas apenas para usuários autenticados com role `admin` neste primeiro momento.

A tabela `opportunities` guardará os dados da oportunidade sem exigir vínculo com uma futura entidade de organização.

## Motivo

Uma opportunity representa uma vaga, projeto, freelance, voluntariado ou mentoria publicada para a comunidade.

Esse conceito é diferente de `developers`, que representa perfis públicos de pessoas desenvolvedoras.

Também é diferente de `users`, que representa identidades autenticáveis e permissões de acesso.

Usar `opportunities` como entidade própria evita forçar o domínio para apenas vagas formais e permite evoluir o produto sem misturar responsabilidades.

Restringir escrita a `admin` reduz risco operacional enquanto o projeto ainda não possui contas de organizações, ownership ou fluxo de aprovação.

Permitir leitura pública de oportunidades publicadas favorece descoberta sem exigir login.

## Consequências

A API passa a ter um CRUD real de oportunidades com validação de campos obrigatórios, enums e canal de contato.

O banco passa a ter uma tabela `opportunities` com status para controlar a visibilidade pública.

O MVP ainda não permite que organizações publiquem ou gerenciem oportunidades diretamente.

O vínculo entre uma oportunidade e uma organização cadastrada fica para uma mudança futura, quando houver modelo de organizações e ownership.

O status `published` se torna a regra de exposição pública. Oportunidades em `draft`, `closed` ou `archived` não aparecem nas leituras públicas.
