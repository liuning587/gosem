package cosem

import (
	"bytes"
	"encoding/hex"
	. "gosem/axdr"
)

type GetDataResult struct {
	IsData bool
	Value  interface{}
}

func CreateGetDataResultAsData(value DlmsData) *GetDataResult {
	return &GetDataResult{true, value}
}

func CreateGetDataResultAsResult(value accessResultTag) *GetDataResult {
	return &GetDataResult{false, value}
}

func CreateGetDataResult(value interface{}) *GetDataResult {
	switch val := value.(type) {
	case DlmsData:
		return CreateGetDataResultAsData(val)
	case accessResultTag:
		return CreateGetDataResultAsResult(val)
	default:
		panic("Value must be either DlmsData or accessResultTag")
	}
}

func (dt *GetDataResult) Encode() []byte {
	var output bytes.Buffer
	if dt.IsData == true {
		output.WriteByte(0x1)
		value := dt.Value.(DlmsData)
		output.Write(value.Encode())
	} else {
		output.WriteByte(0x0)
		value := dt.Value.(accessResultTag)
		output.WriteByte(byte(value))
	}

	return output.Bytes()
}

func (dt *GetDataResult) ValueAsData() DlmsData {
	if !dt.IsData {
		panic("Value is DataAccessResult!")
	}

	return dt.Value.(DlmsData)
}

func (dt *GetDataResult) ValueAsAccess() accessResultTag {
	if dt.IsData {
		panic("Value is DlmsData!")
	}

	return dt.Value.(accessResultTag)
}

func DecodeGetDataResult(src *[]byte) (out GetDataResult, err error) {
	if (*src)[0] == 0x0 {
		out.IsData = false
		out.Value, err = GetAccessTag(uint8((*src)[1]))
		if err == nil {
			(*src) = (*src)[2:]
		}
	} else {
		out.IsData = true
		(*src) = (*src)[1:]
		decoder := NewDataDecoder(src)
		out.Value, err = decoder.Decode(src)
	}

	return
}

// DataBlockG is DataBlock for the GET-response. Result must be either byte slice
// or accessResultTag after creation, or else it will fail on Encode()
type DataBlockG struct {
	LastBlock   bool
	BlockNumber uint32
	IsResult    bool
	Result      interface{}
}

func CreateDataBlockGAsData(lastBlock bool, blockNum uint32, result interface{}) *DataBlockG {
	switch res := result.(type) {
	case string:
		bt, e := hex.DecodeString(res)
		if e != nil {
			panic(e)
		}
		return &DataBlockG{lastBlock, blockNum, false, bt}

	case []byte:
		return &DataBlockG{lastBlock, blockNum, false, res}

	default:
		panic("CreateDataBlockGAsData result only accept hexstring or byte slice.")
	}
}

func CreateDataBlockGAsResult(lastBlock bool, blockNum uint32, result accessResultTag) *DataBlockG {
	return &DataBlockG{lastBlock, blockNum, true, result}
}

func CreateDataBlockG(lastBlock bool, blockNum uint32, result interface{}) *DataBlockG {
	switch res := result.(type) {
	case string:
		bt, e := hex.DecodeString(res)
		if e != nil {
			panic(e)
		}
		return CreateDataBlockGAsData(lastBlock, blockNum, bt)

	case []byte:
		return CreateDataBlockGAsData(lastBlock, blockNum, res)

	case accessResultTag:
		return CreateDataBlockGAsResult(lastBlock, blockNum, res)

	default:
		panic("DataBlockG result only accept hexstring, byte slice, or DataAccessResult.")
	}
}

func (dt *DataBlockG) Encode() []byte {
	var output bytes.Buffer

	if dt.LastBlock {
		output.WriteByte(0x1)
	} else {
		output.WriteByte(0x0)
	}

	blk, e := EncodeDoubleLongUnsigned(dt.BlockNumber)
	if e != nil {
		panic(e)
	}
	output.Write(blk)

	if dt.IsResult == true {
		output.WriteByte(0x1)
		value := dt.Result.(accessResultTag)
		output.WriteByte(byte(value))
	} else {
		output.WriteByte(0x0)
		value := dt.Result.([]byte)
		output.WriteByte(byte(len(value)))
		output.Write(value)
	}

	return output.Bytes()
}

func (dt *DataBlockG) ResultAsBytes() []byte {
	if dt.IsResult {
		panic("Value is DataAccessResult!")
	}

	return dt.Result.([]byte)
}

func (dt *DataBlockG) ResultAsAccess() accessResultTag {
	if !dt.IsResult {
		panic("Value is byte slice!")
	}

	return dt.Result.(accessResultTag)
}

func DecodeDataBlockG(src *[]byte) (out DataBlockG, err error) {
	if (*src)[0] == 0x0 {
		out.LastBlock = false
	} else {
		out.LastBlock = true
	}
	(*src) = (*src)[1:]

	_, out.BlockNumber, err = DecodeDoubleLongUnsigned(src)

	if (*src)[0] == 0x0 {
		out.IsResult = false
	} else {
		out.IsResult = true
	}
	(*src) = (*src)[1:]

	if out.IsResult {
		out.Result, err = GetAccessTag(uint8((*src)[0]))
		(*src) = (*src)[0:]
	} else {
		_, val, e := DecodeLength(src)
		if e != nil {
			err = e
			return
		}
		out.Result = (*src)[:val]
		(*src) = (*src)[val:]
	}

	return
}

// DataBlockSA is DataBlock for the SET-request, ACTION-request and ACTION-response
type DataBlockSA struct {
	LastBlock   bool
	BlockNumber uint32
	Raw         []byte
}

func CreateDataBlockSA(lastBlock bool, blockNum uint32, result interface{}) *DataBlockSA {
	switch res := result.(type) {
	case string:
		bt, e := hex.DecodeString(res)
		if e != nil {
			panic(e)
		}
		return &DataBlockSA{lastBlock, blockNum, bt}

	case []byte:
		return &DataBlockSA{lastBlock, blockNum, res}

	default:
		panic("DataBlockSA result only accept hexstring or byte slice.")
	}
}

func (dt *DataBlockSA) Encode() []byte {
	var output bytes.Buffer

	if dt.LastBlock {
		output.WriteByte(0x1)
	} else {
		output.WriteByte(0x0)
	}

	blk, e := EncodeDoubleLongUnsigned(dt.BlockNumber)
	if e != nil {
		panic(e)
	}
	output.Write(blk)
	output.Write(dt.Raw)

	return output.Bytes()
}

// Response of ActionRequest. ReturnParam is optional parameter therefore pointer
type ActionResponseWithOptData struct {
	Result      actionResultTag
	ReturnParam *GetDataResult
}

func CreateActionResponseWithOptData(result actionResultTag, returnParam *GetDataResult) *ActionResponseWithOptData {

	return &ActionResponseWithOptData{Result: result, ReturnParam: returnParam}
}

func (dt *ActionResponseWithOptData) Encode() []byte {
	var output bytes.Buffer

	output.WriteByte(byte(dt.Result))

	if dt.ReturnParam == nil {
		output.WriteByte(0x0)
	} else {
		output.WriteByte(0x1)
		output.Write(dt.ReturnParam.Encode())
	}

	return output.Bytes()
}
