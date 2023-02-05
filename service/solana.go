package services

import (
	ctx "context"
	"errors"
	"github.com/babilu-online/common/context"
	token_metadata "github.com/gagliardetto/metaplex-go/clients/token-metadata"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"log"
	"os"
)

type SolanaService struct {
	context.DefaultService
	client *rpc.Client
}

const SOLANA_SVC = "solana_svc"

func (svc SolanaService) Id() string {
	return SOLANA_SVC
}

func (svc *SolanaService) Start() error {
	svc.client = rpc.New(os.Getenv("RPC_URL"))

	return nil
}

func (svc *SolanaService) Client() *rpc.Client {
	return svc.client
}

func (svc *SolanaService) RecentBlockhash() (solana.Hash, error) {
	bhash, err := svc.Client().GetRecentBlockhash(ctx.Background(), rpc.CommitmentFinalized)
	if err != nil {
		return solana.Hash{}, err
	}

	return bhash.Value.Blockhash, nil
}

func (svc *SolanaService) TokenData(key solana.PublicKey) (*token_metadata.Metadata, error) {
	var meta token_metadata.Metadata

	ata, _, err := solana.FindTokenMetadataAddress(key)
	if err != nil {
		return nil, err
	}

	err = svc.Client().GetAccountDataBorshInto(ctx.TODO(), ata, &meta)
	if err != nil {
		return nil, err
	}

	return &meta, nil
}

func (svc *SolanaService) CreatorKeys(tokenMint solana.PublicKey) ([]solana.PublicKey, error) {
	metadata, err := svc.TokenData(tokenMint)
	if err != nil {
		log.Printf("%s creatorKeys err: %s", tokenMint, err)
		return nil, err
	}

	if metadata.Data.Creators == nil {
		return nil, errors.New("unable to find creators")
	}

	creatorKeys := make([]solana.PublicKey, len(*metadata.Data.Creators))
	for i, c := range *metadata.Data.Creators {
		creatorKeys[i] = c.Address
	}
	return creatorKeys, nil
}
