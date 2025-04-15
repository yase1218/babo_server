package uuid

import (
	"errors"
	"sync"
	"time"
)

const (
	epoch              int64 = 1577808000000 // 起始时间戳
	workerIdBits       uint8 = 5             // 机器id所占的位数
	datacenterIdBits   uint8 = 5             // 数据中心id所占的位数
	maxWorkerId        int64 = -1 ^ (-1 << workerIdBits)
	maxDatacenterId    int64 = -1 ^ (-1 << datacenterIdBits)
	sequenceBits       uint8 = 12 // 序列号所占的位数
	workerIdShift      uint8 = sequenceBits
	datacenterIdShift  uint8 = sequenceBits + workerIdBits
	timestampLeftShift uint8 = sequenceBits + workerIdBits + datacenterIdBits
	sequenceMask       int64 = -1 ^ (-1 << sequenceBits)
)

type Generator struct {
	sync.Mutex
	init          bool
	lastTimestamp int64
	workerId      int64
	datacenterId  int64
	sequence      int64
}

var genner *Generator

func Init(workerId, datacenterId int64) error {
	var err error
	genner, err = new_gener(workerId, datacenterId)
	if err != nil {
		return err
	}
	return nil
}

func new_gener(workerId, datacenterId int64) (*Generator, error) {
	if workerId < 0 || workerId > maxWorkerId {
		return nil, errors.New("worker Id can't be greater than %d or less than 0")
	}
	if datacenterId < 0 || datacenterId > maxDatacenterId {
		return nil, errors.New("datacenter Id can't be greater than %d or less than 0")
	}
	return &Generator{
		init:          true,
		lastTimestamp: 0,
		workerId:      workerId,
		datacenterId:  datacenterId,
		sequence:      0,
	}, nil
}

func Generate() (int64, error) {
	genner.Lock()
	defer genner.Unlock()

	if !genner.init {
		return 0, errors.New("uuid generator not initialized")
	}

	now := time.Now().UnixNano() / int64(time.Millisecond)

	if genner.lastTimestamp == now {
		genner.sequence = (genner.sequence + 1) & sequenceMask
		if genner.sequence == 0 {
			// 等待下一个毫秒
			for now <= genner.lastTimestamp {
				now = time.Now().UnixNano() / int64(time.Millisecond)
			}
		}
	} else {
		genner.sequence = 0
	}

	if now < genner.lastTimestamp {
		return 0, errors.New("clock moved backwards")
	}

	genner.lastTimestamp = now

	return ((now - epoch) <<
		timestampLeftShift) |
		(genner.datacenterId << datacenterIdShift) |
		(genner.workerId << workerIdShift) |
		genner.sequence, nil
}
