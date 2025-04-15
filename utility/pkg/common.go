package pkg

import (
	"babo/utility/zlog"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"

	"sort"

	"go.uber.org/zap"
)

const MaxUint = (0xFF_FF_FF_FF)

func I64Pairs(l []*I64Pair) {
	sort.Sort(I64PairSlice(l))
}

type I64Pair struct {
	K int64
	V int64
}

func NewI64Pair(k, v int64) *I64Pair {
	return &I64Pair{
		K: k,
		V: v,
	}
}

type I64PairSlice []*I64Pair

func (p I64PairSlice) Len() int           { return len(p) }
func (p I64PairSlice) Less(i, j int) bool { return p[i].K < p[j].K }
func (p I64PairSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func I64SlicePairs(l []*I64SlicePair) {
	sort.Sort(I64SlicePairSlice(l))
}

type I64SlicePair struct {
	K int64
	V []int64
}

func NewI64SlicePair(k int64, v []int64) *I64SlicePair {
	return &I64SlicePair{
		K: k,
		V: v,
	}
}

type I64SlicePairSlice []*I64SlicePair

func (p I64SlicePairSlice) Len() int           { return len(p) }
func (p I64SlicePairSlice) Less(i, j int) bool { return p[i].K < p[j].K }
func (p I64SlicePairSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func ProtectError() {
	if err := recover(); err != nil {
		buff := make([]byte, 1024*4)
		n := runtime.Stack(buff, false)
		zlog.Error("panic", zap.Any("err", err), zap.String("stack", string(buff[:n])))
	}
}

func WaitForTerminate() {
	exitChan := make(chan struct{})
	signalChan := make(chan os.Signal, 1)
	go func() {
		defer ProtectError()
		<-signalChan
		close(exitChan)
	}()
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-exitChan
}

func StrToInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

type ItemBrief struct {
	Id    int64
	Count int64
}

func BuildItemBrief(id, count int64) *ItemBrief {
	return &ItemBrief{
		Id:    id,
		Count: count,
	}
}

func BuildItemBriefs(source [][]int64) []*ItemBrief {
	simpleList := make([]*ItemBrief, 0, len(source))
	for _, v := range source {
		if len(v) != 2 {
			continue
		}
		simple := &ItemBrief{
			Id:    v[0],
			Count: v[1],
		}
		simpleList = append(simpleList, simple)
	}
	return simpleList
}

func Max(x, y int64) int64 {
	if x > y {
		return x
	}
	return y
}

func Min(x, y int64) int64 {
	if x < y {
		return x
	}
	return y
}

func ToInt64(i interface{}) (int64, error) {
	if i == nil {
		return 0, fmt.Errorf("ToInt64 nil")
	}

	switch v := i.(type) {
	case int64:
		return v, nil
	default:
		return 0, fmt.Errorf("ToInt64 type error")
	}
}

func InSlice(s []int64, v int64) bool {
	for _, i := range s {
		if i == v {
			return true
		}
	}
	return false
}
