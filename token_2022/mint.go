package token_2022

import (
	"encoding/binary"
	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"log"
)

type Mint2022 struct {
	// Optional authority used to mint new tokens. The mint authority may only be provided during
	// mint creation. If no mint authority is present then the mint has a fixed supply and no
	// further tokens may be minted.
	MintAuthority *solana.PublicKey `bin:"optional"`

	// Total supply of tokens.
	Supply uint64

	// Number of base 10 digits to the right of the decimal place.
	Decimals uint8

	// Is `true` if this structure has been initialized
	IsInitialized bool

	// Optional authority to freeze token accounts.
	FreezeAuthority *solana.PublicKey `bin:"optional"`

	TvlData []byte

	extensions *MintExtensions
}

func (mint *Mint2022) Extensions() (*MintExtensions, error) {
	if mint.extensions != nil {
		return mint.extensions, nil
	}

	mintExt := MintExtensions{}
	//log.Printf("TVL: %v", mint.TvlData)
	dec := bin.NewBinDecoder(mint.TvlData)

	isInit, _ := dec.ReadBool()
	if !isInit {
		return nil, nil
	}

extLooop:
	for {
		if int(dec.Position()) >= dec.Len()-1 {
			break extLooop
		}

		ext, _ := dec.ReadUint16(binary.LittleEndian)
		//log.Printf("EXT:  %v - Pos: %v - %v", ext, dec.Position(), mint.TvlData[int(dec.Position()):])

		switch ext {
		case ExtUninitialized:
			break extLooop
		case ExtTransferFeeConfig:
			var cfg TransferFeeConfig
			err := cfg.UnmarshalWithDecoder(dec)
			if err != nil {
				return &mintExt, err
			}
			mintExt.TransferFeeConfig = &cfg

			break
		case ExtDefaultAccountState:
			var cfg DefaultAccountState
			err := cfg.UnmarshalWithDecoder(dec)
			if err != nil {
				return &mintExt, err
			}
			mintExt.DefaultAccountState = &cfg
			break
		case ExtNonTransferable:
			mintExt.NonTransferable = true

			_, err := dec.ReadUint16(binary.LittleEndian)
			if err != nil {
				return &mintExt, err
			}

			break
		case ExtPermanentDelegate:
			var cfg PermanentDelegate

			err := cfg.UnmarshalWithDecoder(dec)
			if err != nil {
				return &mintExt, err
			}
			mintExt.PermanentDelegate = &cfg
			break
		case ExtMetadataPointer:
			var cfg MetadataPointer

			err := cfg.UnmarshalWithDecoder(dec)
			if err != nil {
				return &mintExt, err
			}
			mintExt.MetadataPointer = &cfg
			break
		case ExtTokenMetadata:
			var cfg TokenMetadata

			err := cfg.UnmarshalWithDecoder(dec)
			if err != nil {
				return &mintExt, err
			}
			mintExt.TokenMetadata = &cfg
			log.Printf("%+v\n", cfg)
			break
		case ExtGroupPointer:
			var cfg GroupPointer

			err := cfg.UnmarshalWithDecoder(dec)
			if err != nil {
				return &mintExt, err
			}
			mintExt.GroupPointer = &cfg
			break
		case ExtGroupMemberPointer:
			var cfg GroupMemberPointer

			err := cfg.UnmarshalWithDecoder(dec)
			if err != nil {
				return &mintExt, err
			}
			mintExt.GroupMemberPointer = &cfg
			break
		//case ExtTokenGroup:
		//	log.Println("ExtTokenGroup")
		//case ExtGroupMemberPointer:
		//	log.Println("ExtGroupMemberPointer")
		//case ExtTokenGroupMember:
		//	log.Println("ExtTokenGroupMember")
		default:
			iLen, err := dec.ReadUint16(binary.LittleEndian)
			if err != nil {
				return &mintExt, err
			}
			//log.Printf("Mint Missing Decoder: %v - LEN: %v", ext, iLen)
			//log.Printf(" - SKIPPING: %v", mint.TvlData[int(dec.Position()):int(dec.Position())+int(iLen)])
			//log.Printf(" - NEXT: %v", mint.TvlData[int(dec.Position())+int(iLen):int(dec.Position())+int(iLen)+6])
			_, _ = dec.ReadNBytes(int(iLen)) //Read out that amount of bytes
			break
		}

	}

	mint.extensions = &mintExt
	return &mintExt, nil
}

func (mint *Mint2022) TransferFee() (uint16, error) {
	if len(mint.TvlData) == 0 {
		return 0, nil
	}

	exts, err := mint.Extensions()
	if err != nil {
		return 0, err
	}
	if exts == nil || exts.TransferFeeConfig == nil {
		return 0, nil
	}

	return exts.TransferFeeConfig.NewerTransferFee.TransferFeeBasisPoints, nil
}

func (mint *Mint2022) UnmarshalWithDecoder(dec *bin.Decoder) (err error) {
	{
		v, err := dec.ReadUint32(binary.LittleEndian)
		if err != nil {
			return err
		}
		if v == 1 {
			v, err := dec.ReadNBytes(32)
			if err != nil {
				return err
			}
			mint.MintAuthority = solana.PublicKeyFromBytes(v).ToPointer()
		} else {
			// discard:
			_, err := dec.ReadNBytes(32)
			if err != nil {
				return err
			}
		}
	}
	{
		v, err := dec.ReadUint64(binary.LittleEndian)
		if err != nil {
			return err
		}
		mint.Supply = v
	}
	{
		v, err := dec.ReadUint8()
		if err != nil {
			return err
		}
		mint.Decimals = v
	}
	{
		v, err := dec.ReadBool()
		if err != nil {
			return err
		}
		mint.IsInitialized = v
	}
	{
		v, err := dec.ReadUint32(binary.LittleEndian)
		if err != nil {
			return err
		}
		if v == 1 {
			v, err := dec.ReadNBytes(32)
			if err != nil {
				return err
			}
			mint.FreezeAuthority = solana.PublicKeyFromBytes(v).ToPointer()
		} else {
			// discard:
			_, err := dec.ReadNBytes(32)
			if err != nil {
				return err
			}
		}
	}
	{
		//Read the remaining buffer up to 165
		_, _ = dec.ReadNBytes(165 - int(dec.Position()))

		mint.TvlData, err = dec.ReadNBytes(dec.Len() - int(dec.Position()))
		if err != nil {
			return err
		}
	}
	return nil
}

func (mint Mint2022) MarshalWithEncoder(encoder *bin.Encoder) (err error) {
	{
		if mint.MintAuthority == nil {
			err = encoder.WriteUint32(0, binary.LittleEndian)
			if err != nil {
				return err
			}
			empty := solana.PublicKey{}
			err = encoder.WriteBytes(empty[:], false)
			if err != nil {
				return err
			}
		} else {
			err = encoder.WriteUint32(1, binary.LittleEndian)
			if err != nil {
				return err
			}
			err = encoder.WriteBytes(mint.MintAuthority[:], false)
			if err != nil {
				return err
			}
		}
	}
	err = encoder.WriteUint64(mint.Supply, binary.LittleEndian)
	if err != nil {
		return err
	}
	err = encoder.WriteUint8(mint.Decimals)
	if err != nil {
		return err
	}
	err = encoder.WriteBool(mint.IsInitialized)
	if err != nil {
		return err
	}
	{
		if mint.FreezeAuthority == nil {
			err = encoder.WriteUint32(0, binary.LittleEndian)
			if err != nil {
				return err
			}
			empty := solana.PublicKey{}
			err = encoder.WriteBytes(empty[:], false)
			if err != nil {
				return err
			}
		} else {
			err = encoder.WriteUint32(1, binary.LittleEndian)
			if err != nil {
				return err
			}
			err = encoder.WriteBytes(mint.FreezeAuthority[:], false)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
