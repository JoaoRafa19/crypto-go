/***************************************************************
 * Arquivo: transport.go
 * Descrição: Definição da interface de transporte de rede.
 * Autor: JoaoRafa19
 * Data de criação: 2024-2025
 * Versão: 0.0.1
 * Licença: MIT License
 * Observações: 
 ***************************************************************/

package network

type NetAddr string

type Transport interface {
	Consume() <-chan RPC
	Connect(Transport) error
	SendMessage(NetAddr, []byte) error
	Addr() NetAddr
	Broadcast([]byte) error
}
