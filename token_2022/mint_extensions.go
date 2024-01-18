package token_2022

import (
	"encoding/binary"
	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
)

type MintExtensions struct {
	NonTransferable     bool
	TransferFeeConfig   *TransferFeeConfig
	DefaultAccountState *DefaultAccountState
	PermanentDelegate   *PermanentDelegate
	MetadataPointer     *MetadataPointer
	GroupPointer        *GroupPointer
	GroupMemberPointer  *GroupMemberPointer
	TokenMetadata       *TokenMetadata
}

type TransferFeeConfig struct {
	TransferFeeConfigAuthority *solana.PublicKey `bin:"optional"`
	WithdrawWithheldAuthority  *solana.PublicKey `bin:"optional"`
	WithheldAmount             uint64
	OlderTransferFee           TransferFee
	NewerTransferFee           TransferFee
}

func (mint *TransferFeeConfig) UnmarshalWithDecoder(dec *bin.Decoder) (err error) {
	{
		_, err := dec.ReadUint16(binary.LittleEndian)
		if err != nil {
			return err
		}
		v, err := dec.ReadNBytes(32)
		if err != nil {
			return err
		}
		mint.TransferFeeConfigAuthority = solana.PublicKeyFromBytes(v).ToPointer()
	}
	{
		v, err := dec.ReadNBytes(32)
		if err != nil {
			return err
		}
		mint.WithdrawWithheldAuthority = solana.PublicKeyFromBytes(v).ToPointer()
	}
	{
		v, err := dec.ReadUint64(binary.LittleEndian)
		if err != nil {
			return err
		}
		mint.WithheldAmount = v
	}
	{
		err := mint.OlderTransferFee.UnmarshalWithDecoder(dec)
		if err != nil {
			return err
		}
	}
	{
		err := mint.NewerTransferFee.UnmarshalWithDecoder(dec)
		if err != nil {
			return err
		}
	}

	return nil
}

type TransferFee struct {
	Epoch                  uint64
	MaximumFee             uint64
	TransferFeeBasisPoints uint16
}

func (mint *TransferFee) UnmarshalWithDecoder(dec *bin.Decoder) (err error) {
	{
		v, err := dec.ReadUint64(binary.LittleEndian)
		if err != nil {
			return err
		}
		mint.Epoch = v
	}
	{
		v, err := dec.ReadUint64(binary.LittleEndian)
		if err != nil {
			return err
		}
		mint.MaximumFee = v
	}
	{
		v, err := dec.ReadUint16(binary.LittleEndian)
		if err != nil {
			return err
		}
		mint.TransferFeeBasisPoints = v
	}

	return nil
}

// DefaultAccountState - LEN:1
type DefaultAccountState struct {
	State uint8
}

func (s *DefaultAccountState) String() string {
	switch s.State {
	case 0:
		return "Uninitialized"
	case 1:
		return "Initialized"
	case 2:
		return "Frozen"
	}
	return ""
}

func (s *DefaultAccountState) UnmarshalWithDecoder(dec *bin.Decoder) (err error) {
	{
		_, err := dec.ReadUint16(binary.LittleEndian)
		if err != nil {
			return err
		}

		s.State, err = dec.ReadUint8()
		if err != nil {
			return err
		}
	}
	return nil
}

// PermanentDelegate - LEN:33
type PermanentDelegate struct {
	Delegate *solana.PublicKey `bin:"optional"` // 33
}

func (s *PermanentDelegate) UnmarshalWithDecoder(dec *bin.Decoder) (err error) {
	{
		_, err := dec.ReadUint8()
		if err != nil {
			return err
		}

		_, err = dec.ReadBool()
		if err != nil {
			return err
		}

		v2, err := dec.ReadNBytes(32)
		if err != nil {
			return err
		}
		s.Delegate = solana.PublicKeyFromBytes(v2).ToPointer()
	}
	return nil
}

// MetadataPointer - LEN:66
type MetadataPointer struct {
	Authority       *solana.PublicKey `bin:"optional"` // 33
	MetadataAddress *solana.PublicKey `bin:"optional"` // 33
}

func (s *MetadataPointer) UnmarshalWithDecoder(dec *bin.Decoder) (err error) {
	{
		_, err := dec.ReadUint8()
		if err != nil {
			return err
		}

		_, err = dec.ReadBool()
		if err != nil {
			return err
		}

		v2, err := dec.ReadNBytes(32)
		if err != nil {
			return err
		}
		s.Authority = solana.PublicKeyFromBytes(v2).ToPointer()

		v3, err := dec.ReadNBytes(32)
		if err != nil {
			return err
		}
		s.MetadataAddress = solana.PublicKeyFromBytes(v3).ToPointer()
	}
	return nil
}

// GroupPointer - LEN:66
type GroupPointer struct {
	Authority    *solana.PublicKey `bin:"optional"` // 33
	GroupAddress *solana.PublicKey `bin:"optional"` // 33
}

func (s *GroupPointer) UnmarshalWithDecoder(dec *bin.Decoder) (err error) {
	{
		_, err := dec.ReadUint8()
		if err != nil {
			return err
		}

		_, err = dec.ReadBool()
		if err != nil {
			return err
		}

		v2, err := dec.ReadNBytes(32)
		if err != nil {
			return err
		}
		s.Authority = solana.PublicKeyFromBytes(v2).ToPointer()

		v3, err := dec.ReadNBytes(32)
		if err != nil {
			return err
		}
		s.GroupAddress = solana.PublicKeyFromBytes(v3).ToPointer()
	}
	return nil
}

// GroupMemberPointer - LEN:66
type GroupMemberPointer struct {
	Authority     *solana.PublicKey `bin:"optional"` // 33
	MemberAddress *solana.PublicKey `bin:"optional"` // 33
}

func (s *GroupMemberPointer) UnmarshalWithDecoder(dec *bin.Decoder) (err error) {
	{
		_, err := dec.ReadUint8()
		if err != nil {
			return err
		}

		_, err = dec.ReadBool()
		if err != nil {
			return err
		}

		v2, err := dec.ReadNBytes(32)
		if err != nil {
			return err
		}
		s.Authority = solana.PublicKeyFromBytes(v2).ToPointer()

		v3, err := dec.ReadNBytes(32)
		if err != nil {
			return err
		}
		s.MemberAddress = solana.PublicKeyFromBytes(v3).ToPointer()
	}
	return nil
}

// TokenMetadata - LEN:66
type TokenMetadata struct {
	Authority          *solana.PublicKey `bin:"optional"` // 33
	Mint               solana.PublicKey  // 32
	Name               string
	Symbol             string
	Uri                string
	AdditionalMetadata []string
}

func (s *TokenMetadata) UnmarshalWithDecoder(dec *bin.Decoder) (err error) {
	{
		_, err := dec.ReadUint16(binary.LittleEndian)
		if err != nil {
			return err
		}

		v, err := dec.ReadNBytes(32)
		if err != nil {
			return err
		}
		s.Authority = solana.PublicKeyFromBytes(v).ToPointer()

		v, err = dec.ReadNBytes(32)
		if err != nil {
			return err
		}
		s.Mint = solana.PublicKeyFromBytes(v)

		x, _ := dec.ReadUint32(binary.LittleEndian)
		v2, err := dec.ReadNBytes(int(x))
		if err != nil {
			return err
		}
		s.Name = string(v2)

		x, _ = dec.ReadUint32(binary.LittleEndian)
		v2, err = dec.ReadNBytes(int(x))
		if err != nil {
			return err
		}
		s.Symbol = string(v2)

		x, _ = dec.ReadUint32(binary.LittleEndian)

		v2, err = dec.ReadNBytes(int(x))
		if err != nil {
			return err
		}
		s.Uri = string(v2)
	}
	return nil
}
