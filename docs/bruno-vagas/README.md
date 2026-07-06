# Bruno collection local

Esta pasta contém a collection local da API para uso com Bruno.

## Importar

1. Abra o Bruno.
2. Importe a workspace/collection a partir da pasta `docs/bruno-vagas`.
3. Selecione o ambiente `local`.

## Variáveis

- `baseUrl`: URL da API local, por padrão `http://localhost:8080`.
- `adminEmail`: e-mail do admin local.
- `adminPassword`: senha do admin local.
- `token`: token JWT retornado pelo login.

## Uso

1. Execute a API local.
2. Rode o request `Login`.
3. Copie o valor do campo `token` da resposta para a variável `token` do ambiente `local`.
4. Rode o request `Create developer`.

Os requests `List developers` e `Get developer by id` ficam disponíveis para validar leitura da API.
