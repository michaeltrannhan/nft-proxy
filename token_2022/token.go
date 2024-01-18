package token_2022

import (
	"encoding/binary"
	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/token"
)

type Account2022 struct {
	// The mint associated with this account
	Mint solana.PublicKey

	// The owner of this account.
	Owner solana.PublicKey

	// The amount of tokens this account holds.
	Amount uint64

	// If `delegate` is `Some` then `delegated_amount` represents
	// the amount authorized by the delegate
	Delegate *solana.PublicKey `bin:"optional"`

	// The account's state
	State token.AccountState

	// If is_some, this is a native token, and the value logs the rent-exempt reserve. An Account
	// is required to be rent-exempt, so the value is used by the Processor to ensure that wrapped
	// SOL accounts do not drop below this threshold.
	IsNative *uint64 `bin:"optional"`

	// The amount delegated
	DelegatedAmount uint64

	// Optional authority to close the account.
	CloseAuthority *solana.PublicKey `bin:"optional"`

	ExtensionCount uint8
	TvlData        []byte
}

func (t *Account2022) UnmarshalWithDecoder(dec *bin.Decoder) (err error) {
	{
		v, err := dec.ReadNBytes(32)
		if err != nil {
			return err
		}
		t.Mint = solana.PublicKeyFromBytes(v)
	}
	{
		v, err := dec.ReadNBytes(32)
		if err != nil {
			return err
		}
		t.Owner = solana.PublicKeyFromBytes(v)
	}
	{
		v, err := dec.ReadUint64(binary.LittleEndian)
		if err != nil {
			return err
		}
		t.Amount = v
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
			t.Delegate = solana.PublicKeyFromBytes(v).ToPointer()
		} else {
			// discard:
			_, err := dec.ReadNBytes(32)
			if err != nil {
				return err
			}
		}
	}
	{
		v, err := dec.ReadUint8()
		if err != nil {
			return err
		}
		t.State = token.AccountState(v)
	}
	{
		v, err := dec.ReadUint32(binary.LittleEndian)
		if err != nil {
			return err
		}
		if v == 1 {
			v, err := dec.ReadUint64(bin.LE)
			if err != nil {
				return err
			}
			t.IsNative = &v
		} else {
			// discard:
			_, err := dec.ReadUint64(bin.LE)
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
		t.DelegatedAmount = v
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
			t.CloseAuthority = solana.PublicKeyFromBytes(v).ToPointer()
		} else {
			// discard:
			_, err := dec.ReadNBytes(32)
			if err != nil {
				return err
			}
		}
	}
	{
		t.ExtensionCount, err = dec.ReadUint8()
		if err != nil {
			return err
		}
	}
	{
		t.TvlData, err = dec.ReadNBytes(dec.Len() - int(dec.Position()))
		if err != nil {
			return err
		}
	}
	return nil
}

func (t Account2022) MarshalWithEncoder(encoder *bin.Encoder) (err error) {
	{
		err = encoder.WriteBytes(t.Mint[:], false)
		if err != nil {
			return err
		}
	}
	{
		err = encoder.WriteBytes(t.Owner[:], false)
		if err != nil {
			return err
		}
	}
	{
		err = encoder.WriteUint64(t.Amount, bin.LE)
		if err != nil {
			return err
		}
	}
	{
		if t.Delegate == nil {
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
			err = encoder.WriteBytes(t.Delegate[:], false)
			if err != nil {
				return err
			}
		}
	}
	err = encoder.WriteUint8(uint8(t.State))
	if err != nil {
		return err
	}
	{
		if t.IsNative == nil {
			err = encoder.WriteUint32(0, binary.LittleEndian)
			if err != nil {
				return err
			}
			err = encoder.WriteUint64(0, bin.LE)
			if err != nil {
				return err
			}
		} else {
			err = encoder.WriteUint32(1, binary.LittleEndian)
			if err != nil {
				return err
			}
			err = encoder.WriteUint64(*t.IsNative, bin.LE)
			if err != nil {
				return err
			}
		}
	}
	{
		err = encoder.WriteUint64(t.DelegatedAmount, bin.LE)
		if err != nil {
			return err
		}
	}
	{
		if t.CloseAuthority == nil {
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
			err = encoder.WriteBytes(t.CloseAuthority[:], false)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (t *Account2022) WithheldAmount() (uint64, error) {
	dec := bin.NewBinDecoder(t.TvlData)

	dLen, err := dec.ReadUint8()
	if err != nil {
		return 0, err
	}
	_, err = dec.ReadNBytes(int(dLen))
	if err != nil {
		return 0, err
	}

	withheldAmount, err := dec.ReadUint64(binary.LittleEndian)
	if err != nil {
		return 0, err
	}

	return withheldAmount, nil
}
