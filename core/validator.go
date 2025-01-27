package core

import "fmt"

type Validator interface {
	ValidateBlock(*Block) error
}

type BlockValidator struct {
	Bc *BlockChain
}

func NewBlockValidator(bc *BlockChain) *BlockValidator {
	return &BlockValidator{
		Bc: bc,
	}
}

func (v *BlockValidator) ValidateBlock(b *Block) error {
	if v.Bc.HasBlock(b.Height) {
		return fmt.Errorf("chain alredy contains block (%d) with hash (%s)", b.Height, b.Hash(BlockHasher{}))
	}

	if b.Height != v.Bc.Height()+1 {
		return fmt.Errorf("block (%s) too high", b.Hash(BlockHasher{}))
	}

	prevHeader, err := v.Bc.GetHeader(b.Height - 1)
	if err != nil {
		return err
	}

	hash := BlockHasher{}.Hash(prevHeader)
	if hash != b.PrevBlockHash {
		return fmt.Errorf("the hash of the previous block (%s) is invalid", b.PrevBlockHash)
	}

	if err := b.Verify(); err != nil {
		return err
	}

	return nil
}
