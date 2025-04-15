package pkg

// import (
// 	"fmt"

// 	"gitlab.wowstudio.com/server/rambo/pbgo"
// )

// type Condition struct {
// 	Type pbgo.ConditionType

// 	i64        int64
// 	i64s       []int64
// 	i64map     map[int64]int64
// 	strI64map  map[string]int64
// 	itemCounts map[int64]int64
// 	i64Pair    *I64Pair
// }

// func (c *Condition) I64() int64 {
// 	return c.i64
// }

// func (c *Condition) I64s() []int64 {
// 	return c.i64s
// }

// func (c *Condition) I64Map() map[int64]int64 {
// 	return c.i64map
// }

// func (c *Condition) StrI64Map() map[string]int64 {
// 	return c.strI64map
// }

// func (c *Condition) ItemAmounts() map[int64]int64 {
// 	return c.itemCounts
// }

// // I64MapKey 获取map中的i64的 value
// func (c *Condition) I64MapKey() int64 {
// 	for k := range c.i64map {
// 		return k
// 	}
// 	return 0
// }

// // I64MapValue 获取map中的i64的 value
// func (c *Condition) I64MapValue() int64 {
// 	for _, v := range c.i64map {
// 		return v
// 	}
// 	return 0
// }

// func BuildCondition(typeInt int64, stringss [][]string) (*Condition, error) {
// 	c := &Condition{
// 		Type: pbgo.ConditionType(typeInt),
// 	}
// 	var err error
// 	switch pbgo.ConditionType(typeInt) {
// 	case pbgo.ConditionType_CTRoleLevelAmount, pbgo.ConditionType_CTChapterZoneCompleted, pbgo.ConditionType_CTRoleLevel, pbgo.ConditionType_CTAcqLevelEquip,
// 		pbgo.ConditionType_CTRuneLevel, pbgo.ConditionType_CTContinueCarnivalTemp, pbgo.ConditionType_CTChooseAwakenSkillChapterTemp, pbgo.ConditionType_CTContinueKillAllChapterTemp,
// 		pbgo.ConditionType_CTNotHurtChapterTemp, pbgo.ConditionType_CTFullHpZoneTemp, pbgo.ConditionType_CTCustomerLevel, pbgo.ConditionType_CTAreaShopStepClear,
// 		pbgo.ConditionType_CTAreaShopConsume, pbgo.ConditionType_CTKillAllChapterTemp, pbgo.ConditionType_CTChapterFail, pbgo.ConditionType_CTFacilityLv, pbgo.ConditionType_CTPassChapterNumTemp,
// 		pbgo.ConditionType_CTSkillLvAmount, pbgo.ConditionType_CTStartByType, pbgo.ConditionType_CTQualityRole:
// 		c.i64map, err = func(strss [][]string) (map[int64]int64, error) {
// 			m := make(map[int64]int64)
// 			for _, ss := range strss {
// 				if len(ss) != 2 {
// 					err := fmt.Errorf("BuildCondition i64map len %d, values %v", len(ss), strss)
// 					return nil, err
// 				}
// 				k, err := StrToInt64(ss[0])
// 				if err != nil {
// 					err := fmt.Errorf("BuildCondition i64map StrToInt64[0] %v", strss)
// 					return nil, err
// 				}
// 				v, err := StrToInt64(ss[1])
// 				if err != nil {
// 					err := fmt.Errorf("BuildCondition i64map StrToInt64[1] %v", strss)
// 					return nil, err
// 				}
// 				m[k] = v
// 			}
// 			return m, nil
// 		}(stringss)
// 		if err != nil {
// 			return nil, err
// 		}
// 	case pbgo.ConditionType_CTEffectParam:
// 		c.strI64map, err = func(strss [][]string) (map[string]int64, error) {
// 			m := make(map[string]int64)
// 			for _, ss := range strss {
// 				if len(ss) != 2 {
// 					err := fmt.Errorf("BuildCondition strI64map len %d, values %v", len(ss), strss)
// 					return nil, err
// 				}
// 				k := ss[0]
// 				v, err := StrToInt64(ss[1])
// 				if err != nil {
// 					err := fmt.Errorf("BuildCondition strI64map StrToInt64 %v", strss)
// 					return nil, err
// 				}
// 				m[k] = v
// 			}
// 			return m, nil
// 		}(stringss)
// 		if err != nil {
// 			return nil, err
// 		}
// 	case pbgo.ConditionType_CTItemCost:
// 		c.itemCounts, err = func(strss [][]string) (map[int64]int64, error) {
// 			m := make(map[int64]int64)
// 			for _, ss := range strss {
// 				if len(ss) != 2 {
// 					err := fmt.Errorf("BuildCondition itemAmounts len %d, values %v", len(ss), strss)
// 					return nil, err
// 				}
// 				k, err := StrToInt64(ss[0])
// 				if err != nil {
// 					err := fmt.Errorf("BuildCondition itemAmounts StrToInt64[0] %v", strss)
// 					return nil, err
// 				}
// 				v, err := StrToInt64(ss[1])
// 				if err != nil {
// 					err := fmt.Errorf("BuildCondition itemAmounts StrToInt64[1] %v", strss)
// 					return nil, err
// 				}
// 				m[k] = v
// 			}
// 			return m, nil
// 		}(stringss)
// 	case pbgo.ConditionType_CTStoryLineCompleted, pbgo.ConditionType_CTHaveCustomerId, pbgo.ConditionType_CTHaveRole:
// 		c.i64s, err = func(strss [][]string) ([]int64, error) {
// 			m := make([]int64, 0, len(strss))
// 			if len(strss) == 0 || len(strss[0]) == 0 {
// 				err := fmt.Errorf("BuildCondition i64s len %d, values %v", len(strss), strss)
// 				return nil, err
// 			}
// 			for _, v := range strss[0] {
// 				i, err := StrToInt64(v)
// 				if err != nil {
// 					err := fmt.Errorf("BuildCondition i64s StrToInt64 %v", strss)
// 					return nil, err
// 				}
// 				m = append(m, i)
// 			}
// 			return m, nil
// 		}(stringss)
// 	default:
// 		c.i64, err = func(strss [][]string) (int64, error) {
// 			if len(strss) == 0 || len(strss[0]) != 1 {
// 				err := fmt.Errorf("BuildCondition i64s len %d, values %v", len(strss), strss)
// 				return 0, err
// 			}
// 			i, err := StrToInt64(strss[0][0])
// 			if err != nil {
// 				err := fmt.Errorf("BuildCondition i64 StrToInt64 %v", strss)
// 				return 0, err
// 			}
// 			return i, nil
// 		}(stringss)
// 	}
// 	// 做一个条件验证 防止配置表有非法的数据
// 	// switch ConditionType(typeInt) {
// 	// case ConditionTypeChapterZoneCompleted:
// 	// 	chapterId := o.I64MapKey()
// 	// 	zoneId := o.I64MapValue()
// 	// 	zoneData := GetChapterZonePropDataByDataId(chapterId, zoneId)
// 	// 	if zoneData == nil {
// 	// 		log.Panic(validator.SprintErrFieldMsg(d.Table(), d.Key(), valueFieldName, values, "章节布局表没有找到对应关卡"))
// 	// 	}
// 	// }

// 	return c, nil
// }
