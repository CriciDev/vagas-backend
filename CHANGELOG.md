# Changelog

Todas as mudanças relevantes deste projeto serão documentadas neste arquivo.

O formato segue [Keep a Changelog](https://keepachangelog.com/pt-BR/1.1.0/) e este projeto usa [Versionamento Semântico](https://semver.org/lang/pt-BR/).

## [0.0.1] - 2026-07-05

### Changed

- `GET /api/opportunities` agora aceita `page` e `page_size` e retorna um envelope com `data` e `meta` (`page`, `page_size`, `total`, `total_pages`).

### Breaking

- `GET /api/opportunities` deixou de retornar um array JSON puro e passou a retornar `{ "data": [...], "meta": { ... } }`.
