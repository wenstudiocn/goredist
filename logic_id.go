package dist

import (
"errors"
)

// 设计规则参看 kingame 的 README
var (
	ErrInvalidLogicId = errors.New("Invalid logic id")
)

type LogicID struct {
	funcId uint64 //999
	cataId uint64 //99
	subId  uint64 //99
	instId uint64 //99
	id uint64
}

func NewLogicIDByID(id uint64) (*LogicID, error) {
	sid := &LogicID{}

	sid.instId = id % 100
	sid.subId = id / 100 % 100
	sid.cataId = id / 10000 % 100
	sid.funcId = id / 1000000
	sid.id = id

	if sid.funcId >= 1000 || sid.funcId <= 0 {
		return nil, ErrInvalidLogicId
	}
	return sid, nil
}

func NewLogicIDByParts(funcId, cataId, subId, instId uint64) (*LogicID, error) {
	sid := &LogicID{
		funcId: funcId,
		cataId:cataId,
		subId:subId,
		instId:instId,
		id: funcId * 1000000 + cataId * 10000 + subId * 100 + instId,
	}

	if sid.funcId <= 0 || sid.funcId >= 1000 ||
		sid.cataId > 99 ||
		sid.subId > 99 ||
		sid.instId > 99 {
		return nil, ErrInvalidLogicId
	}
	return sid, nil
}

func (self *LogicID)FuncId() uint64 {
	return self.funcId
}

func (self *LogicID)CataId() uint64 {
	return self.cataId
}

func (self *LogicID)SubId() uint64 {
	return self.subId
}

func (self *LogicID)InstId() uint64 {
	return self.instId
}

func (self *LogicID)Id() uint64 {
	return self.id
}