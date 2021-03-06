package helpers

import (
	"bytes"
	"encoding/binary"
	"github.com/idena-network/idena-go/common"
	"github.com/pkg/errors"
	"math/big"
)

var indexOufOfRange = errors.New("index out of range")
var noValue = errors.New("no value")

func assertLen(index int, args ...[]byte) error {
	if index >= len(args) {
		return indexOufOfRange
	}
	return nil
}

func ExtractAddr(index int, args ...[]byte) (common.Address, error) {
	if err := assertLen(index, args...); err != nil {
		return common.Address{}, err
	}
	if args[index] == nil {
		return common.Address{}, noValue
	}
	addr := common.Address{}
	addr.SetBytes(args[index])
	return addr, nil
}

func ExtractUInt64(index int, args ...[]byte) (uint64, error) {
	if err := assertLen(index, args...); err != nil {
		return 0, err
	}
	if args[index] == nil {
		return 0, noValue
	}
	var ret uint64
	buf := bytes.NewBuffer(args[index])
	if err := binary.Read(buf, binary.LittleEndian, &ret); err != nil {
		return 0, errors.Wrapf(err, "index=%v", index)
	}
	return ret, nil
}

func ExtractUInt16(index int, args ...[]byte) (uint16, error) {
	if err := assertLen(index, args...); err != nil {
		return 0, err
	}
	if args[index] == nil {
		return 0, noValue
	}
	var ret uint16
	buf := bytes.NewBuffer(args[index])
	if err := binary.Read(buf, binary.LittleEndian, &ret); err != nil {
		return 0, errors.Wrapf(err, "index=%v", index)
	}
	return ret, nil
}

func ExtractByte(index int, args ...[]byte) (byte, error) {
	if err := assertLen(index, args...); err != nil {
		return 0, err
	}
	if len(args[index]) == 0 {
		return 0, errors.Wrapf(noValue, "index=%v", index)
	}
	return args[index][0], nil
}

func ExtractBigInt(index int, args ...[]byte) (*big.Int, error) {
	if err := assertLen(index, args...); err != nil {
		return nil, err
	}
	if args[index] == nil {
		return nil, errors.Wrapf(noValue, "index=%v", index)
	}
	ret := new(big.Int)
	ret.SetBytes(args[index])
	return ret, nil
}

func ExtractArray(index int, args ...[]byte) ([]byte, error) {
	if err := assertLen(index, args...); err != nil {
		return nil, err
	}
	if args[index] == nil {
		return nil, errors.Wrapf(noValue, "index=%v", index)
	}
	return args[index], nil
}
