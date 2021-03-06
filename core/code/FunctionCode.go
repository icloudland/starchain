package code

import (."github.com/icloudland/starchain/core/contract"
	."github.com/icloudland/starchain/common"
	"io"
	"github.com/icloudland/starchain/common/serialization"
	"fmt"
	"github.com/icloudland/starchain/common/log"
)

type FunctionCode struct {
	Code []byte
	ParameterTypes []ContractParameterType
	ReturnType ContractParameterType

	codeHash Uint160
}


// method of SerializableData
func (fc *FunctionCode) Serialize(w io.Writer) error {
	err := serialization.WriteByte(w, byte(fc.ReturnType))
	if err != nil {
		return err
	}
	err = serialization.WriteVarBytes(w, ContractParameterTypeToByte(fc.ParameterTypes))
	if err != nil {
		return err
	}
	err = serialization.WriteVarBytes(w, fc.Code)
	if err != nil {
		return err
	}
	return nil
}

// method of SerializableData
func (fc *FunctionCode) Deserialize(r io.Reader) error {
	returnType, err := serialization.ReadByte(r)
	if err != nil {
		return err
	}
	fc.ReturnType = ContractParameterType(returnType)

	parameterTypes, err := serialization.ReadVarBytes(r)
	if err != nil {
		return err
	}
	fc.ParameterTypes = ByteToContractParameterType(parameterTypes)

	fc.Code, err = serialization.ReadVarBytes(r)
	if err != nil {
		return err
	}

	return nil
}


// method of ICode
// Get code
func (fc *FunctionCode) GetCode() []byte {
	return fc.Code
}

// method of ICode
// Get the list of parameter value
func (fc *FunctionCode) GetParameterTypes() []ContractParameterType {
	return fc.ParameterTypes
}

// method of ICode
// Get the list of return value
func (fc *FunctionCode) GetReturnType() ContractParameterType {
	return fc.ReturnType
}

// method of ICode
// Get the hash of the smart contract
func (fc *FunctionCode) CodeHash() Uint160 {
	var log = log.NewLog()
	u160 := Uint160{}
	if fc.codeHash == u160 {
		u160, err := ToCodeHash(fc.Code)
		if err != nil {
			log.Debug(fmt.Sprintf("[FunctionCode] ToCodeHash err=%s", err))
			return u160
		}
		fc.codeHash = u160
	}
	return fc.codeHash
}