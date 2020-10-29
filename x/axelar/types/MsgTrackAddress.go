package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/axelarnetwork/axelar-core/x/axelar/exported"
)

// Ensure MsgTrackAddress implements sdk.Msg interface
var _ sdk.Msg = &MsgTrackAddress{}

type MsgTrackAddress struct {
	Sender  sdk.AccAddress
	Address exported.ExternalChainAddress
}

func NewMsgTrackAddress(sender sdk.AccAddress, address exported.ExternalChainAddress) MsgTrackAddress {
	return MsgTrackAddress{
		Sender:  sender,
		Address: address,
	}
}

func (msg MsgTrackAddress) Route() string {
	return RouterKey
}

func (msg MsgTrackAddress) Type() string {
	return "TrackAddress"
}

func (msg MsgTrackAddress) ValidateBasic() error {
	if msg.Sender.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing sender")
	}
	if msg.Address.IsInvalid() {
		return sdkerrors.Wrap(ErrInvalidExternalAddress, msg.Address.String())
	}

	return nil
}

func (msg MsgTrackAddress) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgTrackAddress) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}