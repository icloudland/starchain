package smartcontract

import (
	"github.com/icloudland/starchain/common"
	"github.com/icloudland/starchain/core/contract"
	"github.com/icloudland/starchain/smartcontract/types"
	"math/big"
	"github.com/icloudland/starchain/smartcontract/storage"
	"github.com/icloudland/starchain/vm/avm/interfaces"
	"github.com/icloudland/starchain/smartcontract/service"
	sig"github.com/icloudland/starchain/core/signature"
	"bytes"
	"strconv"
	"github.com/icloudland/starchain/vm/avm"
	."github.com/icloudland/starchain/vm/avm/types"
	"github.com/icloudland/starchain/errors"
	"github.com/icloudland/starchain/core/ledger"
	"github.com/icloudland/starchain/core/transaction"
	"github.com/icloudland/starchain/smartcontract/states"
	"github.com/icloudland/starchain/core/asset"
	"github.com/icloudland/starchain/common/serialization"
)

type Engine interface {
	Create(caller common.Uint160,code []byte) ([]byte,error)
	Call(caller common.Uint160,codeHash common.Uint160,input []byte) ([]byte,error)
}

type SmartContract struct {
	Engine 	Engine
	Code 	[]byte
	Input 	[]byte
	ParameterTypes []contract.ContractParameterType
	Caller	common.Uint160
	CodeHash common.Uint160
	VMType	types.VmType
	ReturnType contract.ContractParameterType
}

type Context struct {
	Language       types.LangType
	Caller         common.Uint160
	StateMachine   *service.StateMachine
	DBCache        storage.DBCache
	Code           []byte
	Input          []byte
	CodeHash       common.Uint160
	Time           *big.Int
	BlockNumber    *big.Int
	CacheCodeTable interfaces.ICodeTable
	SignableData   sig.SignableData
	Gas            common.Fixed64
	ReturnType     contract.ContractParameterType
	ParameterTypes []contract.ContractParameterType
}


func NewSmartContract(context *Context) (*SmartContract, error) {
	if vmType, ok := types.LangVm[context.Language]; ok {
		var e Engine
		switch vmType {
		case types.AVM:
			e = avm.NewExecutionEngine(
				context.SignableData,
				new(avm.ECDsaCrypto),
				context.CacheCodeTable,
				context.StateMachine,
				context.Gas,
			)
		// case types.EVM:
		// 	e = evm.NewExecutionEngine(context.DBCache, context.Time, context.BlockNumber, context.Gas)
		}

		return &SmartContract{
			Engine:         e,
			Code:           context.Code,
			CodeHash:       context.CodeHash,
			Input:          context.Input,
			Caller:         context.Caller,
			VMType:         vmType,
			ReturnType:     context.ReturnType,
			ParameterTypes: context.ParameterTypes,
		}, nil
	} else {
		return nil, errors.NewDetailErr(errors.NewErr("Not Support Language Type!"), errors.ErrNoCode, "")
	}

}

func (sc *SmartContract) DeployContract() ([]byte, error) {
	return sc.Engine.Create(sc.Caller, sc.Code)
}

func (sc *SmartContract) InvokeContract() (interface{}, error) {
	_, err := sc.Engine.Call(sc.Caller, sc.CodeHash, sc.Input)
	if err != nil {
		return nil, err
	}
	return sc.InvokeResult()
}

func (sc *SmartContract) InvokeResult() (interface{}, error) {
	switch sc.VMType {
	case types.AVM:
		engine := sc.Engine.(*avm.ExecutionEngine)
		if engine.GetEvaluationStackCount() > 0 && avm.Peek(engine).GetStackItem() != nil {
			switch sc.ReturnType {
			case contract.Boolean:
				return avm.PopBoolean(engine), nil
			case contract.Integer:
				return avm.PopBigInt(engine).String(), nil
			case contract.ByteArray:
				bs := avm.PopByteArray(engine)
				return common.BytesToInt(bs), nil
			case contract.String:
				return string(avm.PopByteArray(engine)), nil
			case contract.Hash160, contract.Hash256:
				return common.BytesToHexString(common.ToArrayReverse(avm.PopByteArray(engine))), nil
			case contract.PublicKey:
				return common.BytesToHexString(avm.PopByteArray(engine)), nil
			case contract.Object:
				data := avm.PeekStackItem(engine)
				switch data.(type) {
				case *Boolean:
					return data.GetBoolean(), nil
				case *Integer:
					return data.GetBigInteger(), nil
				case *ByteArray:
					return common.BytesToInt(data.GetByteArray()), nil
				case *InteropInterface:
					interop := data.GetInterface()
					switch interop.(type) {
					case *ledger.Header:
						return service.GetHeaderInfo(interop.(*ledger.Header)), nil
					case *ledger.Block:
						return service.GetBlockInfo(interop.(*ledger.Block)), nil
					case *transaction.Transaction:
						return service.GetTransactionInfo(interop.(*transaction.Transaction)), nil
					case *states.AccountState:
						return service.GetAccountInfo(interop.(*states.AccountState)), nil
					case *asset.Asset:
						return service.GetAssetInfo(interop.(*asset.Asset)), nil
					}
				}
			}
		}
	case types.EVM:
	}
	return nil, nil
}

func (sc *SmartContract) InvokeParamsTransform() ([]byte, error) {
	switch sc.VMType {
	case types.AVM:
		builder := avm.NewParamsBuilder(new(bytes.Buffer))
		b := bytes.NewBuffer(sc.Input)
		for _, k := range sc.ParameterTypes {
			switch k {
			case contract.Boolean:
				p, err := serialization.ReadBool(b)
				if err != nil {
					return nil, err
				}
				builder.EmitPushBool(p)
			case contract.Integer:
				p, err := serialization.ReadVarBytes(b)
				if err != nil {
					return nil, err
				}
				i, err := strconv.ParseInt(string(p), 10, 64)
				if err != nil {
					return nil, err
				}
				builder.EmitPushInteger(int64(i))
			case contract.Hash160, contract.Hash256:
				p, err := serialization.ReadVarBytes(b)
				if err != nil {
					return nil, err
				}
				builder.EmitPushByteArray(common.ToArrayReverse(p))
			case contract.ByteArray, contract.String:
				p, err := serialization.ReadVarBytes(b)
				if err != nil {
					return nil, err
				}
				builder.EmitPushByteArray(p)
			case contract.Array:
			//val, err := serialization.ReadVarUint(b, 0)
			//if err != nil {
			//	return nil, err
			//}

			}
		}
		builder.EmitPushCall(sc.CodeHash.ToArray())
		return builder.ToArray(), nil
	case types.EVM:
	}
	return nil, nil
}

