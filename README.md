# Teste Aprendendo Conceitos de Blockchain

Made from scratch, generic propose, production ready blockchain for uses like office files, criptocurrency
Build using TDD


**Test Coverage**

![Coverage](coverage/badge.svg)

## Features

- [v] Server
	- [X] Creating blocks
	- [x] Connecting transports
	- [x] Broadcasting transactions
	- [ ] Broadcasting block
- [X] Block
    - [X] Block's hash
    - [x] Test
- [X] Transaction
    - [x] Transaction list Hash
    - [x] Test
- [x] Key
- [x] Transport => tcp, udp, 
    - [X] Local transport layer
- [X] Crypto Keypairs and signature
- [X] Block Signing
- [X] Blockchain struct
- [X] Storage (memory storage)
- [X] Transaction Encoding/Decoding
- [X] Block Encoding/Decoding



## Todos
Improvements and fixes that can be implemented

- [ ] Add a database or a better storage method to store transactions and block data

## Types 

- Hash

```go
type Hash [32]uint8
type Address [20]uint8
```

### Mistakes to remember 

On the struct Transaction on Signing the transaction the object was missing the value of the transaction's `Signature`, returnning a null value 

```go
func TestSignTransaction(t *testing.T) {
	privateKey := crypto.GeneratePrivateKey()

	tx := &Transaction{
		Data: []byte("foo bar baz"),
	}

	assert.Nil(t, tx.Sign(privateKey))
	assert.NotNil(t, tx.Signature) // FAIL

}
```
Beacuse of the method signatue (tx T) insted of (tx *T), the object was missing the reference

**`old`** :
```go
func (tx Transaction) Sign(privKey crypto.PrivateKey) error {
```

**`fixed`**:
 ```go
func (tx *Transaction) Sign(privKey crypto.PrivateKey) error {
```

## Instalação

Para instalar as dependências do projeto, execute:
```bash
go mod tidy
```
Para executar o projeto 
```
make run
```

## Como Contribuir

1. Faça um fork do projeto
2. Crie uma nova branch (`git checkout -b feature/nome-da-feature`)
3. Faça commit das suas alterações (`git commit -am 'Adiciona nova feature'`)
4. Faça push para a branch (`git push origin feature/nome-da-feature`)
5. Abra um Pull Request

## Licença

Este projeto está licenciado sob a Licença MIT - veja o arquivo [LICENSE](LICENSE) para mais detalhes.
