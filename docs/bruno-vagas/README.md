# Bruno collection local

Esta pasta contém a collection local da API para uso com Bruno.

## Importar

1. Abra o Bruno.
2. Importe a workspace/collection a partir da pasta `docs/bruno-vagas`.
3. Selecione o ambiente `local`.

## Variáveis

- `baseUrl`: URL da API local, por padrão `http://localhost:8080`.
- `token`: token JWT retornado pelo login.
- `developerId`: id retornado na criação do developer.
- `opportunityId`: id retornado na criação da oportunidade.

## Uso

1. Execute a API local.
2. Rode o request `Login`.
3. O request `Login` preenche automaticamente `token`.
4. Rode o request `Create developer`.
5. O request `Create developer` preenche automaticamente `developerId`.
6. Rode o request `Get developer by id`.
7. Rode `Create opportunity` para criar uma oportunidade.
8. Rode `Get opportunity by id` para consultar a oportunidade criada.

Os requests `List developers` e `Get developer by id` ficam disponíveis para validar leitura da API.

Os requests de oportunidades cobrem listagem, filtros, criação, consulta, atualização e exclusão.
