package token_metadata

import (
	token_metadata "github.com/gagliardetto/metaplex-go/clients/token-metadata"
	"github.com/gagliardetto/solana-go"
)

type Protocol uint

const (
	PROTOCOL_LEGACY Protocol = iota
	PROTOCOL_TOKEN22_MINT
	PROTOCOL_LIBREPLEX
	PROTOCOL_METAPLEX_CORE
)

type Metadata struct {
	Key             token_metadata.Key
	UpdateAuthority solana.PublicKey
	Mint            solana.PublicKey
	Data            Data

	// Immutable, once flipped, all sales of this metadata are considered secondary.
	PrimarySaleHappened bool

	// Whether or not the data struct is mutable, default is not
	IsMutable bool

	// Collection
	Collection *token_metadata.Collection `bin:"optional"`

	Protocol Protocol `bin:"-"`
}

type Data struct {
	// The name of the asset
	Name string

	// The symbol for the asset
	Symbol string

	// URI pointing to JSON representing the asset
	Uri string

	// Royalty basis points that goes to creators in secondary sales (0-10000)
	SellerFeeBasisPoints uint16

	// Array of creators, optional
	Creators *[]token_metadata.Creator `bin:"optional"`
}
