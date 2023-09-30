# concorrente-projeto

## Proposta

Realizar um experimento para avaliar o desempenho de um programa concorrente implementado em duas linguagens de programação diferentes.

## Objetivos

- Implementar um programa concorrente em duas linguagens de programação diferentes (JavaScript e GO).
- Realizar testes de carga e estresse.
- Comparar o desempenho das duas implementações.
- Gerar uma estatística de desempenho com um intervalo de confiança sobre as linguagens.

Ao final do experimento, espera-se possuir um artigo com definição de como se construir um programa concorrente em cada uma das linguagens, bem como uma análise comparativa de desempenho entre as duas implementações.

## Metodologia

Inicialmente iremos construir um programa concorrente em JavaScript e em GO com os seguintes requisitos:

- O programa deve possuir um Web Server para receber requisições HTTP.
  - Um endpoint para receber uma Request de solicitação de criação de recurso.
  - O endpoint deve receber um JSON com os dados do recurso a ser criado.
  - O endpoint deve inserir o recurso em uma Queue.
- O programa deve se conectar a um banco de dados para persistir os dados.
- O programa deve possuir um Background Job para processar a Queue e inserir em Batch no banco de dados.

A concorrência se dá pelo fato de que o programa deve ser capaz de receber múltiplas requisições HTTP e processá-las em paralelo.

Após a implementação, iremos realizar testes de carga e estresse para avaliar o desempenho das duas implementações.

- Com o teste de carga vamos avaliar o desempenho dos programas para uma determinada quantidade de requisições com objetivo de coletar dados como o tempo de resposta.
- Com o teste de estresse vamos avaliar o desempenho das duas implementações e avaliar até que ponto cada uma delas consegue suportar com objetivo de coletar dados como tempo de resposta, limiar de req/s que cada server suporta, etc.
- Para reprodutibilidade dos testes iremos limitar a quantidade de recursos de hardware disponíveis para os programas e o banco de dados, além de que estes irão executar virtualizados.
