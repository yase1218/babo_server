package weights

// import (
// 	"fmt"
// 	"strings"

// 	"gitlab.wowstudio.com/server/rambo/utility/i64"

// 	"github.com/pkg/errors"
// )

// // Weights 权重。arrays 左开右闭区间，[0，x),[x, y),...[y+n, max)
// type Weights struct {
// 	array []*weight
// }

// type weight struct {
// 	max   int64
// 	value int64
// }

// // New 创建一个权重计算器。raws=[]{weight，value}，遍历raws，累加weight。
// func New(raws [][]int64) (*Weights, error) {
// 	array := make([]*weight, 0, len(raws))

// 	var prev int64
// 	for _, kv := range raws {
// 		if len(kv) != 2 {
// 			return nil, fmt.Errorf("[weights.NewWeight] 每个kv的长度必须=2。%+v", kv)
// 		}
// 		if kv[0] < 0 {
// 			return nil, fmt.Errorf("[weights.NewWeight] 每个权重必须>=0。%+v", kv)
// 		}
// 		max := prev + kv[0]
// 		array = append(array, &weight{max: max, value: kv[1]})
// 		prev = max
// 	}

// 	return &Weights{array: array}, nil
// }

// // Default 创建一个平均权重的权重计算器。
// func Default(values []int64) *Weights {
// 	array := make([]*weight, 0, len(values))

// 	var prev int64
// 	for _, v := range values {
// 		max := prev + 100
// 		array = append(array, &weight{max: max, value: v})
// 		prev = max
// 	}

// 	return &Weights{array: array}
// }

// func (ws *Weights) RandByBase(base int64) (v int64, ok bool) {
// 	if ws.Len() <= 0 {
// 		return 0, false
// 	}
// 	return ws.Value(i64.Random(base))
// }

// func (ws *Weights) Rand() (v int64, ok bool) {
// 	if ws.Len() <= 0 {
// 		return 0, false
// 	}
// 	return ws.Value(i64.Random(ws.Max()))
// }

// func (ws *Weights) Len() int {
// 	return len(ws.array)
// }

// func (ws *Weights) Max() int64 {
// 	l := ws.Len()
// 	if l == 0 {
// 		return 0
// 	}
// 	return ws.array[l-1].max
// }

// func (ws *Weights) Value(r int64) (int64, bool) {
// 	if ws.Len() == 0 {
// 		return 0, false
// 	}
// 	if r == 0 {
// 		return ws.array[0].value, true
// 	}

// 	for _, w := range ws.array {
// 		if r < w.max {
// 			return w.value, true
// 		}
// 	}
// 	return 0, false
// }

// // RandExcept 根据权重随机生成一个结果，尽量保证结果不在排除列表中
// func (ws *Weights) RandExcept(exs map[int64]struct{}) (result int64, ok bool, err error) {
// 	return ws.RandExceptWithValueFunc(exs, func(value int64) int64 {
// 		return value
// 	})
// }

// // RandExceptWithValueFunc 根据权重随机生成一个结果，尽量保证结果不在排除列表中
// // execValueFunc 去重判断时，用 execValueFunc 把 value 处理成不同的值
// func (ws *Weights) RandExceptWithValueFunc(exs map[int64]struct{}, execValueFunc func(value int64) int64) (result int64, ok bool, err error) {
// 	size := ws.Len()
// 	if size == 0 {
// 		return
// 	}
// 	if size == 1 {
// 		r := ws.array[0].value
// 		if _, exist := exs[execValueFunc(r)]; exist {
// 			return
// 		}
// 		result = r
// 		ok = true
// 		return
// 	}

// 	nws, err := ws.rebuildWeightsWithValueFunc(exs, execValueFunc)
// 	if err != nil {
// 		err = errors.New("[Weights.RandExcept] rebuildWeights 失败。")
// 		result, ok = ws.Value(i64.Random(ws.Max()))
// 		return
// 	}
// 	if nws.Len() == 0 {
// 		result, ok = ws.Value(i64.Random(ws.Max()))
// 		return
// 	}
// 	result, ok = nws.Rand()
// 	return
// }

// func (ws *Weights) rebuildWeights(exs map[int64]struct{}) (*Weights, error) {
// 	return ws.rebuildWeightsWithValueFunc(exs, func(value int64) int64 {
// 		return value
// 	})
// }

// func (ws *Weights) rebuildWeightsWithValueFunc(exs map[int64]struct{}, execValueFunc func(value int64) int64) (*Weights, error) {
// 	bases := make([][]int64, 0, len(ws.array))

// 	var start, prev int64
// 	for _, w := range ws.array {
// 		if _, ok := exs[execValueFunc(w.value)]; ok {
// 			prev = w.max
// 			continue
// 		}
// 		bases = append(bases, []int64{start + (w.max - prev), w.value})
// 		start = w.max
// 		prev = w.max
// 	}

// 	return New(bases)
// }

// // RandList 根据权重随机生成一个列表，尽量保证结果不重复
// func (ws *Weights) RandList(count int) (result []int64) {
// 	if count <= 0 {
// 		return
// 	}

// 	array := ws.array
// 	size := ws.Len()

// 	if size == 0 {
// 		return
// 	}

// 	result = make([]int64, 0, count)

// 	if size == count {
// 		for _, w := range array {
// 			result = append(result, w.value)
// 		}
// 		return
// 	}

// 	if size < count {
// 		max := ws.Max()
// 		for _, r := range array {
// 			result = append(result, r.value)
// 		}
// 		for i := 0; i < count-size; i++ {
// 			id, _ := ws.Value(i64.Random(max))
// 			result = append(result, id)
// 		}
// 		return
// 	}

// 	partLen := size / count
// 	lastAdd := size % count

// 	var min int64
// 	for i := 0; i < count; i++ {
// 		l := partLen
// 		if i == count-1 {
// 			l = partLen + lastAdd
// 		}
// 		start := i * partLen
// 		end := start + l - 1

// 		max := array[end].max
// 		id, _ := ws.Value(i64.Random(max-min) + min)
// 		result = append(result, id)
// 		min = array[end].max
// 	}
// 	return
// }

// func (ws *Weights) String() string {
// 	b := &strings.Builder{}
// 	b.WriteString("Weight[")
// 	for i, w := range ws.array {
// 		if i > 0 {
// 			b.WriteString(", ")
// 		}
// 		b.WriteString(fmt.Sprintf("%d:%d", w.max, w.value))
// 	}
// 	b.WriteString("]")
// 	return b.String()
// }
