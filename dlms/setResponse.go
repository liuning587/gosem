package cosem

import (
	"bytes"
	"fmt"
	. "gosem/axdr"
)

type setResponseTag uint8

const (
	TagSetResponseNormal                setResponseTag = 0x1
	TagSetResponseDataBlock             setResponseTag = 0x2
	TagSetResponseLastDataBlock         setResponseTag = 0x3
	TagSetResponseLastDataBlockWithList setResponseTag = 0x4
	TagSetResponseWithList              setResponseTag = 0x5
)

// Value will return primitive value of the target.
// This is used for comparing with non custom typed object
func (s setResponseTag) Value() uint8 {
	return uint8(s)
}

// SetResponse implement CosemI
type SetResponse struct{}

func (gr *SetResponse) New(tag setResponseTag) CosemPDU {
	switch tag {
	case TagSetResponseNormal:
		return &SetResponseNormal{}
	case TagSetResponseDataBlock:
		return &SetResponseDataBlock{}
	case TagSetResponseLastDataBlock:
		return &SetResponseLastDataBlock{}
	case TagSetResponseLastDataBlockWithList:
		return &SetResponseLastDataBlockWithList{}
	case TagSetResponseWithList:
		return &SetResponseWithList{}
	default:
		panic("Tag not recognized!")
	}
}

func (gr *SetResponse) Decode(src *[]byte) (out CosemPDU, err error) {
	if (*src)[0] != TagSetResponse.Value() {
		err = ErrWrongTag(0, (*src)[0], byte(TagSetResponse))
		return
	}

	switch (*src)[1] {
	case TagSetResponseNormal.Value():
		out, err = DecodeSetResponseNormal(src)
	case TagSetResponseDataBlock.Value():
		out, err = DecodeSetResponseDataBlock(src)
	case TagSetResponseLastDataBlock.Value():
		out, err = DecodeSetResponseLastDataBlock(src)
	case TagSetResponseLastDataBlockWithList.Value():
		out, err = DecodeSetResponseLastDataBlockWithList(src)
	case TagSetResponseWithList.Value():
		out, err = DecodeSetResponseWithList(src)
	default:
		err = fmt.Errorf("byte tag not recognized (%v)", (*src)[1])
	}

	return
}

// SetResponseNormal implement CosemPDU
type SetResponseNormal struct {
	InvokePriority uint8
	Result         AccessResultTag
}

func CreateSetResponseNormal(invokeId uint8, result AccessResultTag) *SetResponseNormal {
	return &SetResponseNormal{
		InvokePriority: invokeId,
		Result:         result,
	}
}

func (sr SetResponseNormal) Encode() []byte {
	var buf bytes.Buffer
	buf.WriteByte(TagSetResponse.Value())
	buf.WriteByte(TagSetResponseNormal.Value())
	buf.WriteByte(sr.InvokePriority)
	buf.WriteByte(sr.Result.Value())

	return buf.Bytes()
}

func DecodeSetResponseNormal(src *[]byte) (out SetResponseNormal, err error) {
	if (*src)[0] != TagSetResponse.Value() {
		err = ErrWrongTag(0, (*src)[0], byte(TagSetResponse))
		return
	}
	if (*src)[1] != TagSetResponseNormal.Value() {
		err = ErrWrongTag(0, (*src)[1], byte(TagSetResponseNormal))
		return
	}
	out.InvokePriority = (*src)[2]
	out.Result, err = GetAccessTag(uint8((*src)[3]))
	if err != nil {
		return
	}
	(*src) = (*src)[4:]

	return
}

// SetResponseDataBlock implement CosemPDU
type SetResponseDataBlock struct {
	InvokePriority uint8
	BlockNum       uint32
}

func CreateSetResponseDataBlock(invokeId uint8, blockNum uint32) *SetResponseDataBlock {
	return &SetResponseDataBlock{
		InvokePriority: invokeId,
		BlockNum:       blockNum,
	}
}

func (sr SetResponseDataBlock) Encode() []byte {
	var buf bytes.Buffer
	buf.WriteByte(TagSetResponse.Value())
	buf.WriteByte(TagSetResponseDataBlock.Value())
	buf.WriteByte(sr.InvokePriority)
	blockNum, _ := EncodeDoubleLongUnsigned(sr.BlockNum)
	buf.Write(blockNum)

	return buf.Bytes()
}

func DecodeSetResponseDataBlock(src *[]byte) (out SetResponseDataBlock, err error) {
	if (*src)[0] != TagSetResponse.Value() {
		err = ErrWrongTag(0, (*src)[0], byte(TagSetResponse))
		return
	}
	if (*src)[1] != TagSetResponseDataBlock.Value() {
		err = ErrWrongTag(0, (*src)[1], byte(TagSetResponseDataBlock))
		return
	}
	out.InvokePriority = (*src)[2]
	(*src) = (*src)[3:]

	_, out.BlockNum, err = DecodeDoubleLongUnsigned(src)

	return
}

// SetResponseLastDataBlock implement CosemPDU
type SetResponseLastDataBlock struct {
	InvokePriority uint8
	Result         AccessResultTag
	BlockNum       uint32
}

func CreateSetResponseLastDataBlock(invokeId uint8, result AccessResultTag, blockNum uint32) *SetResponseLastDataBlock {
	return &SetResponseLastDataBlock{
		InvokePriority: invokeId,
		Result:         result,
		BlockNum:       blockNum,
	}
}

func (sr SetResponseLastDataBlock) Encode() []byte {
	var buf bytes.Buffer
	buf.WriteByte(TagSetResponse.Value())
	buf.WriteByte(TagSetResponseLastDataBlock.Value())
	buf.WriteByte(sr.InvokePriority)
	buf.WriteByte(sr.Result.Value())
	blockNum, _ := EncodeDoubleLongUnsigned(sr.BlockNum)
	buf.Write(blockNum)

	return buf.Bytes()
}

func DecodeSetResponseLastDataBlock(src *[]byte) (out SetResponseLastDataBlock, err error) {
	if (*src)[0] != TagSetResponse.Value() {
		err = ErrWrongTag(0, (*src)[0], byte(TagSetResponse))
		return
	}
	if (*src)[1] != TagSetResponseLastDataBlock.Value() {
		err = ErrWrongTag(0, (*src)[1], byte(TagSetResponseDataBlock))
		return
	}
	out.InvokePriority = (*src)[2]
	out.Result, err = GetAccessTag(uint8((*src)[3]))
	if err != nil {
		return
	}
	(*src) = (*src)[4:]

	_, out.BlockNum, err = DecodeDoubleLongUnsigned(src)

	return
}

// SetResponseLastDataBlockWithList implement CosemPDU
type SetResponseLastDataBlockWithList struct {
	InvokePriority uint8
	ResultCount    uint8
	ResultList     []AccessResultTag
	BlockNum       uint32
}

func CreateSetResponseLastDataBlockWithList(invokeId uint8, resList []AccessResultTag, blockNum uint32) *SetResponseLastDataBlockWithList {
	if len(resList) < 1 || len(resList) > 255 {
		panic("ResultList cannot have zero or >255 member")
	}
	return &SetResponseLastDataBlockWithList{
		InvokePriority: invokeId,
		ResultCount:    uint8(len(resList)),
		ResultList:     resList,
		BlockNum:       blockNum,
	}
}

func (sr SetResponseLastDataBlockWithList) Encode() []byte {
	var buf bytes.Buffer
	buf.WriteByte(TagSetResponse.Value())
	buf.WriteByte(TagSetResponseLastDataBlockWithList.Value())
	buf.WriteByte(sr.InvokePriority)
	buf.WriteByte(sr.ResultCount)
	for _, acc := range sr.ResultList {
		buf.WriteByte(acc.Value())
	}
	blockNum, _ := EncodeDoubleLongUnsigned(sr.BlockNum)
	buf.Write(blockNum)

	return buf.Bytes()
}

func DecodeSetResponseLastDataBlockWithList(src *[]byte) (out SetResponseLastDataBlockWithList, err error) {
	if (*src)[0] != TagSetResponse.Value() {
		err = ErrWrongTag(0, (*src)[0], byte(TagSetResponse))
		return
	}
	if (*src)[1] != TagSetResponseLastDataBlockWithList.Value() {
		err = ErrWrongTag(0, (*src)[1], byte(TagSetResponseLastDataBlockWithList))
		return
	}
	out.InvokePriority = (*src)[2]

	out.ResultCount = uint8((*src)[3])
	(*src) = (*src)[4:]
	for i := 0; i < int(out.ResultCount); i++ {
		v, e := GetAccessTag(uint8((*src)[0]))
		if e != nil {
			err = e
			return
		}
		out.ResultList = append(out.ResultList, v)
		(*src) = (*src)[1:]
	}

	_, out.BlockNum, err = DecodeDoubleLongUnsigned(src)

	return
}

// SetResponseWithList implement CosemPDU
type SetResponseWithList struct {
	InvokePriority uint8
	ResultCount    uint8
	ResultList     []AccessResultTag
}

func CreateSetResponseWithList(invokeId uint8, resList []AccessResultTag) *SetResponseWithList {
	if len(resList) < 1 || len(resList) > 255 {
		panic("ResultList cannot have zero or >255 member")
	}
	return &SetResponseWithList{
		InvokePriority: invokeId,
		ResultCount:    uint8(len(resList)),
		ResultList:     resList,
	}
}

func (sr SetResponseWithList) Encode() []byte {
	var buf bytes.Buffer
	buf.WriteByte(TagSetResponse.Value())
	buf.WriteByte(TagSetResponseWithList.Value())
	buf.WriteByte(sr.InvokePriority)
	buf.WriteByte(sr.ResultCount)
	for _, acc := range sr.ResultList {
		buf.WriteByte(acc.Value())
	}

	return buf.Bytes()
}

func DecodeSetResponseWithList(src *[]byte) (out SetResponseWithList, err error) {
	if (*src)[0] != TagSetResponse.Value() {
		err = ErrWrongTag(0, (*src)[0], byte(TagSetResponse))
		return
	}
	if (*src)[1] != TagSetResponseWithList.Value() {
		err = ErrWrongTag(0, (*src)[1], byte(TagSetResponseWithList))
		return
	}
	out.InvokePriority = (*src)[2]

	out.ResultCount = uint8((*src)[3])
	(*src) = (*src)[4:]
	for i := 0; i < int(out.ResultCount); i++ {
		v, e := GetAccessTag(uint8((*src)[0]))
		if e != nil {
			err = e
			return
		}
		out.ResultList = append(out.ResultList, v)
		(*src) = (*src)[1:]
	}

	return
}
