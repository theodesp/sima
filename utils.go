package sima

import (
	"context"
	"github.com/spaolacci/murmur3"
	"io"
	"reflect"
	"runtime"
)

type ReceiverType = func(context context.Context, sender *Topic) interface{}

func GetFunctionName(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

func HashValue(value interface{}) uint64 {
	hash := murmur3.New64()
	switch value.(type) {
	case *Topic:
		io.WriteString(hash, value.(*Topic).name)
		break
	case ReceiverType:
		io.WriteString(hash, GetFunctionName(value.(ReceiverType)))
		break
	default:
		panic("unknown value")
	}

	return hash.Sum64()
}
