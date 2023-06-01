package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCreateTwin{}, "twin/CreateTwin", nil)
	cdc.RegisterConcrete(&MsgUpdateTwin{}, "twin/UpdateTwin", nil)
	cdc.RegisterConcrete(&MsgDeleteTwin{}, "twin/DeleteTwin", nil)
	cdc.RegisterConcrete(&MsgTrain{}, "twin/Train", nil)
	cdc.RegisterConcrete(&MsgConfirmTrainPhaseEnded{}, "twin/ConfirmTrainPhaseEnded", nil)
	cdc.RegisterConcrete(&MsgConfirmBestTrainResultIs{}, "twin/ConfirmBestTrainResultIs", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateTwin{},
		&MsgUpdateTwin{},
		&MsgDeleteTwin{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgTrain{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgConfirmTrainPhaseEnded{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgConfirmBestTrainResultIs{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
