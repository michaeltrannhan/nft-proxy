package metaplex_core

import (
	"encoding/binary"
	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
)

type Asset struct {
	Key             uint8
	Owner           solana.PublicKey
	UpdateAuthority *solana.PublicKey `bin:"optional"`
	Name            string
	Uri             string
	Seq             []uint64 `bin:"optional"`
}

func (asset *Asset) UnmarshalWithDecoder(dec *bin.Decoder) (err error) {
	asset.Key, err = dec.ReadUint8()
	if err != nil {
		return err
	}

	_o, err := dec.ReadBytes(32)
	if err != nil {
		return err
	}
	asset.Owner = solana.PublicKeyFromBytes(_o)

	if v, _ := dec.ReadBool(); v {
		_ua, err := dec.ReadBytes(32)
		if err != nil {
			return err
		}
		pk := solana.PublicKeyFromBytes(_ua)
		asset.UpdateAuthority = &pk
	}

	size, err := dec.ReadUint32(binary.LittleEndian)
	_n, err := dec.ReadBytes(int(size))
	if err != nil {
		return err
	}
	asset.Name = string(_n)

	size, err = dec.ReadUint32(binary.LittleEndian)
	_uri, err := dec.ReadBytes(int(size))
	if err != nil {
		return err
	}
	asset.Uri = string(_uri)
	return nil
}
