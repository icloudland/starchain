package code

import (
	."github.com/icloudland/starchain/core/contract"
	."github.com/icloudland/starchain/common"
)

type ICode interface {
	GetCode() []byte
	GetParameterTypes() []ContractParameterType
	GetReturnTypes() []ContractParameterType

	CodeHash() Uint160
}
