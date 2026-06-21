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

---

### Estrutura do Projeto

```text
.
├── cmd/server/main.go      # Entrada da aplicação
├── go.mod                  # Módulo Go
├── Makefile                # Atalhos de desenvolvimento
└── README.md
```

O projeto segue uma estrutura idiomática para Go:

- O scaffold começa seco: `cmd/server/main.go` é o único código da aplicação.
- Crie diretórios como `internal`, `sql` ou `infra` apenas quando houver comportamento/configuração real.
- Crie pacotes de domínio em `internal/<dominio>` apenas quando houver regras/tipos reais; evite arquivos vazios para antecipar camadas.

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
