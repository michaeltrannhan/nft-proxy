package token_2022

import "math"

type ExtensionType uint8

const (
	/// Used as padding if the account size would otherwise be 355, same as a multisig
	ExtUninitialized = iota
	/// Includes transfer fee rate info and accompanying authorities to withdraw and set the fee
	ExtTransferFeeConfig
	/// Includes withheld transfer fees
	ExtTransferFeeAmount
	/// Includes an optional mint close authority
	ExtMintCloseAuthority
	/// Auditor configuration for confidential transfers
	ExtConfidentialTransferMint
	/// State for confidential transfers
	ExtConfidentialTransferAccount
	/// Specifies the default Account::state for new Accounts
	ExtDefaultAccountState
	/// Indicates that the Account owner authority cannot be changed
	ExtImmutableOwner
	/// Require inbound transfers to have memo
	ExtMemoTransfer
	/// Indicates that the tokens from this mint can't be transfered
	ExtNonTransferable
	/// Tokens accrue interest over time,
	ExtInterestBearingConfig
	/// Locks privileged token operations from happening via CPI
	ExtCpiGuard
	/// Includes an optional permanent delegate
	ExtPermanentDelegate
	/// Indicates that the tokens in this account belong to a non-transferable
	/// mint
	ExtNonTransferableAccount
	/// Mint requires a CPI to a program implementing the "transfer hook"
	/// interface
	ExtTransferHook
	/// Indicates that the tokens in this account belong to a mint with a
	/// transfer hook
	ExtTransferHookAccount
	/// Includes encrypted withheld fees and the encryption public that they are
	/// encrypted under
	ExtConfidentialTransferFeeConfig
	/// Includes confidential withheld transfer fees
	ExtConfidentialTransferFeeAmount
	/// Mint contains a pointer to another account (or the same account) that
	/// holds metadata
	ExtMetadataPointer
	/// Mint contains token-metadata
	ExtTokenMetadata
	/// Mint contains a pointer to another account (or the same account) that
	/// holds group configurations
	ExtGroupPointer
	/// Mint contains token group configurations
	ExtTokenGroup
	/// Mint contains a pointer to another account (or the same account) that
	/// holds group member configurations
	ExtGroupMemberPointer
	/// Mint contains token group member configurations
	ExtTokenGroupMember
	ExtVariableLenMintTest = math.MaxUint16 - 2
	ExtAccountPaddingTest  = iota //TODO Check?
	ExtMintPaddingTest
)
