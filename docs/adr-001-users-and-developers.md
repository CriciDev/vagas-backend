# ADR 001: Separar users e developers

## Status

Aceito

## Contexto

O primeiro CRUD do projeto cria perfis públicos de developers para a comunidade.

Ao mesmo tempo, a API precisa de autenticação básica para permitir que apenas usuários com permissão façam operações de escrita.

Essas duas responsabilidades parecem parecidas no começo, mas representam conceitos diferentes.

## Decisão

O banco terá duas tabelas separadas neste primeiro momento.

`developers` guarda os perfis públicos dos devs da comunidade.

`users` guarda as credenciais e permissões de quem pode autenticar na API.

## Motivo

Um developer é um perfil exibido publicamente, com dados como nome, e-mail, skills, disponibilidade e bio.

Um user é uma identidade autenticável, com e-mail, senha criptografada e role.

Separar os dois modelos evita misturar dados públicos de perfil com dados de autenticação e autorização.

Também deixa o primeiro passo mais simples. Neste momento, só um admin autenticado pode criar, editar ou remover developers.

## Consequências

A API consegue validar login e permissões sem acoplar isso diretamente ao perfil público de developer.

O fluxo atual ainda não permite que cada developer faça login e edite o próprio perfil.

Quando esse fluxo for necessário, uma mudança futura pode criar o vínculo entre `users` e `developers`, por exemplo com um campo `developer.user_id` ou outro modelo de ownership.
