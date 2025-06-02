# Multithreading
Desafio fullcycle - Multithreading

Neste desafio você terá que usar o que aprendemos com Multithreading e APIs para buscar o resultado mais rápido entre duas APIs distintas.

As duas requisições serão feitas simultaneamente para as seguintes APIs:

`https://brasilapi.com.br/api/cep/v1/ + cep`

`http://viacep.com.br/ws/" + cep + "/json/`

Os requisitos para este desafio são:

- Acatar a API que entregar a resposta mais rápida e descartar a resposta mais lenta.

- O resultado da request deverá ser exibido no command line com os dados do endereço, bem como qual API a enviou.

- Limitar o tempo de resposta em 1 segundo. Caso contrário, o erro de timeout deve ser exibido.

## Como executar
Realizar o clone do projeto e dentro da pasta, abrir o terminal e executar o seguinte comando:

`go run main.go {numero do cep}`

### Exemplo

`go run main.go 01001-001`

### Resultado:

Endereço: Praça da Sé, Sé - São Paulo/SP - CEP: 01001001 Fonte: Brasil API` 

