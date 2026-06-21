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

Suba a API e o PostgreSQL:

```bash
docker compose -f infra/docker-compose.yml up --build
```

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
