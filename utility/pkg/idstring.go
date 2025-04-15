package pkg

// import (
// 	"fmt"
// 	"strconv"
// 	"strings"

// 	"github.com/speps/go-hashids"
// )

// const (
// 	idStrLen = 32
// 	salt     = "oidutSWoW"
// )

// var (
// 	h *hashids.HashID
// )

// func init() {
// 	hd := hashids.NewData()
// 	hd.Salt = salt
// 	hd.MinLength = idStrLen

// 	var err error
// 	if h, err = hashids.NewWithData(hd); err != nil {
// 		panic(fmt.Errorf("ID codec init failed. %v", err))
// 	}
// }

// func EncodeId(id int64) (string, error) {
// 	if id < 0 {
// 		return strconv.FormatInt(id, 10), nil
// 	}

// 	str, err := h.EncodeInt64([]int64{id})
// 	if err != nil {
// 		return "", fmt.Errorf("ID codec string failed. id:%d", id)
// 	}
// 	return str, nil
// }

// func DecodeId(str string) (int64, error) {
// 	if strings.IndexRune(str, '-') == 0 {
// 		return strconv.ParseInt(str, 10, 64)
// 	}

// 	ids, err := h.DecodeInt64WithError(str)
// 	if err != nil {
// 		return 0, fmt.Errorf("IDStr decode int64 failed. str:%s", str)
// 	}
// 	if len(ids) == 0 {
// 		return 0, fmt.Errorf("IDStr decode int64 failed,ids lenth err。str:%s", str)
// 	}
// 	return ids[0], nil
// }
