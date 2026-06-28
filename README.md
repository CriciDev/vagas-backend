### Criciúma Devs Jobs — Backend

Interface web open source construída em Angular e Go para conectar profissionais de tecnologia e entidades da região de Criciúma. Aqui, empresas ou instituições podem publicar oportunidades e devs podem criar perfis, demonstrar habilidades e informar disponibilidade.

---

### Objetivo

- **Missão do projeto**: gerar a oportunidade de contribuição em um projeto real para novos desenvolvedores, e assim ter colaborado em um projeto real e open source no seu portifólio.
- **Para organizações, empresas**: evidenciar oportunidades de trabalho ou negócio.
- **Para desenvolvedores**: manter um perfil, demonstrar disponibilidade em projetos.
- **Para a comunidade**: facilitar conexões locais de qualidade com uma experiência simples e direta.

---

### Tecnologias Utilizadas

- **Angular**
- **Golang**
- **Gin**
- **PostgreSQL**
- **Docker**

---

### Estrutura do Projeto

```text
.
├── cmd/server/main.go
├── infra/docker-compose.yml
├── infra/postgres/init/001_init.sql
├── internal/auth
├── internal/config
├── internal/database
├── internal/devs
├── internal/health
├── Dockerfile
├── Makefile
├── go.mod
└── README.md
```

O projeto segue uma estrutura idiomática para Go:

- `cmd/server/main.go` é a entrada da aplicação.
- `internal/auth` concentra autenticação, tokens e autorização por role.
- `internal/devs` concentra o CRUD de developers.
- `internal/database` concentra conexão, schema inicial e seed local.
- `infra` concentra o ambiente local com Docker e PostgreSQL.

---

### Funcionalidades (MVP)

- **Lista de oportunidades**: listagem com busca/filtros simples (ex.: tipo de vaga).
- **Página de devs**: listagem de perfis com principais skills e disponibilidade.
- **Formulários de criação**:
  - empresas/instituições: publicar necessidade/oportunidade
  - devs: criar/editar perfil e sinalizar disponibilidade (aceitando trabalhos)
- **Autenticação básica**: login/logout, proteção de rotas específicas (guards).

---

### Rodando o Projeto

Pré-requisitos:

- Docker
- Docker Compose
- Node.js/npm, apenas se for usar OpenCode + OpenSpec no fluxo de contribuição

Setup de variáveis de ambiente:

```bash
make setup-env
```

Isso cria:

- `/.env` para `make run`
- `/infra/.env` para Docker Compose

Os arquivos reais não devem ser commitados.

Suba a API e o PostgreSQL:

```bash
docker compose --env-file infra/.env -f infra/docker-compose.yml up --build
```

Para rodar sem Docker:

```bash
make run
```

`make run` carrega `/.env` apenas durante a execução, sem alterar o ambiente do terminal.

Verifique a API:

```bash
curl http://localhost:8080/health
```

Credenciais locais do admin:

- E-mail: `admin@criciumadevs.local`
- Senha: `admin123`

Login local:

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@criciumadevs.local","password":"admin123"}'
```

Listar devs:

```bash
curl http://localhost:8080/api/developers
```

Criar dev com token admin:

```bash
curl -X POST http://localhost:8080/api/developers \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"name":"Ada Lovelace","email":"ada@example.com","skills":["Go","Angular"],"available":true,"bio":"Dev da comunidade"}'
```

Atualizar dev:

```bash
curl -X PUT http://localhost:8080/api/developers/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"name":"Ada Lovelace","email":"ada@example.com","skills":["Go","Angular","PostgreSQL"],"available":false,"bio":"Dev da comunidade"}'
```

Remover dev:

```bash
curl -X DELETE http://localhost:8080/api/developers/1 \
  -H "Authorization: Bearer <token>"
```

Rodar testes localmente:

```bash
go test ./...
```

Rode a aplicação:

```bash
make run
```

Rode os testes:

```bash
make test
```

Rode a verificação principal:

```bash
make check
```

Compile o binário:

```bash
make build
```

---

### Fluxo com OpenCode + OpenSpec

Este projeto usa **OpenCode** como agente de desenvolvimento no terminal e **OpenSpec** para registrar intenção, decisões e tarefas antes da implementação.

A ideia é simples: mudanças maiores devem começar com uma proposta revisável, não direto no código. Isso ajuda a comunidade a discutir o que será feito, por que será feito e quais tarefas precisam ser executadas.

Instale o OpenCode:

```bash
curl -fsSL https://opencode.ai/install | bash
```

Ou, se preferir npm:

```bash
npm install -g opencode-ai
```

Instale o OpenSpec:

```bash
npm install -g @fission-ai/openspec@latest
```

Instale as dependências locais dos comandos do OpenCode:

```bash
cd .opencode
npm install
cd ..
```

Abra o OpenCode na raiz do projeto:

```bash
opencode
```

Comandos principais usados neste repositório:

- `/opsx-explore`: explorar uma ideia, problema ou decisão sem alterar código.
- `/opsx-propose <nome-da-mudanca>`: criar uma proposta OpenSpec com intenção, design e tarefas.
- `/opsx-apply <nome-da-mudanca>`: implementar as tarefas de uma mudança já proposta.
- `/opsx-sync <nome-da-mudanca>`: sincronizar specs da mudança com as specs principais.
- `/opsx-archive <nome-da-mudanca>`: arquivar uma mudança concluída.

Fluxo recomendado para contribuições maiores:

1. Escolha uma issue no GitHub.
2. Use `/opsx-explore` se a solução ainda estiver pouco clara.
3. Use `/opsx-propose nome-da-mudanca` para criar a proposta.
4. Discuta a proposta na issue ou no PR, se necessário.
5. Use `/opsx-apply nome-da-mudanca` para implementar.
6. Rode `make check` antes de abrir o PR.
7. Depois da mudança revisada e concluída, use `/opsx-archive nome-da-mudanca`.

Mudanças pequenas, como correções simples de texto, não precisam obrigatoriamente passar por OpenSpec. Use bom senso: se existe decisão de produto, contrato de API, banco de dados ou mudança de comportamento, prefira registrar a intenção antes.

---

### Como Contribuir

Fluxo sugerido:

1. Faça um fork do repositório.
2. Crie uma branch descritiva:
   ```bash
   git checkout -b feat/minha-contribuicao
   ```
3. Implemente sua mudança seguindo o estilo recomendado.
4. Execute/valide localmente.
5. Envie um PR (pull request) com descrição clara do que foi alterado e por quê.

Estilo de código recomendado:

- Siga o **Angular Style Guide** (nomenclatura, separação por responsabilidade).
- Utilize **tipagem explícita** em APIs públicas e services.
- Componentes pequenos e focados; prefira composição a herança.
- CSS: prefira componentes/estilos reutilizáveis em `shared`.
- Go: mantenha o scaffold seco, evite camadas vazias e rode `make check` antes de abrir PR.
- Commits: **Conventional Commits** (ex.: `feat:`, `fix:`, `docs:`).

Abrindo issues:

- Descreva claramente o problema/ideia.
- Inclua passos para reproduzir (se bug) e prints quando relevante.
- Marque com labels apropriadas (ex.: `bug`, `feat`, `good first issue`).

---

### Roadmap

MVP:

- Listar oportunidades e devs
- Criar/editar oportunidade
- Criar/editar perfil de dev
- Autenticação básica e proteção de rotas

Futuro (redesign e incrementos):

- Melhorias de UX/UI e identidade visual
- Filtros avançados (stack, senioridade, modalidade)
- Paginação/infinit scroll
- Upload de avatar/portfólio
- Integrações sociais e SEO básico
- Acessibilidade (WCAG) e i18n

---

### Contato / Comunidade

- Entre no [discord da comunidade](https://discord.gg/XaQhWEucMA)
- Dúvidas e sugestões: abra uma issue ou participe das discussões.

Junte-se a Criciuma Devs para construir junto soluções para a comunidade tech de Criciúma! 🚀
