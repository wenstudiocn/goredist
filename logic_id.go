package dist

import (
	"errors"
)

const (
	LI_FUNC_MULTI = 1000000
	LI_CATA_MULTI = 10000
	LI_SUB_MULTI  = 100
	LI_INST_MULTI = 1
	LI_MAX_ID     = 100
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
	id     uint64
}

func NewLogicIDByID(id uint64) (*LogicID, error) {
	sid := &LogicID{}

	sid.instId = id % LI_SUB_MULTI
	sid.subId = id / LI_SUB_MULTI % LI_MAX_ID
	sid.cataId = id / LI_CATA_MULTI % LI_MAX_ID
	sid.funcId = id / LI_FUNC_MULTI
	sid.id = id

	if sid.funcId >= 1000 || sid.funcId <= 0 {
		return nil, ErrInvalidLogicId
	}
	return sid, nil
}

func NewLogicIDByParts(funcId, cataId, subId, instId uint64) (*LogicID, error) {
	sid := &LogicID{
		funcId: funcId,
		cataId: cataId,
		subId:  subId,
		instId: instId,
		id:     funcId*1000000 + cataId*10000 + subId*100 + instId,
	}

	if sid.funcId <= 0 || sid.funcId >= 1000 ||
		sid.cataId > 99 ||
		sid.subId > 99 ||
		sid.instId > 99 {
		return nil, ErrInvalidLogicId
	}
	return sid, nil
}

func (self *LogicID) FuncId() uint64 {
	return self.funcId
}

func (self *LogicID) CataId() uint64 {
	return self.cataId
}

func (self *LogicID) SubId() uint64 {
	return self.subId
}

func (self *LogicID) InstId() uint64 {
	return self.instId
}

func (self *LogicID) HeaderCata() uint64 {
	return self.funcId*LI_FUNC_MULTI + self.cataId*LI_CATA_MULTI
}

func (self *LogicID) HeaderSub() uint64 {
	return self.funcId*LI_FUNC_MULTI + self.cataId*LI_CATA_MULTI + self.subId*LI_SUB_MULTI
}

func (self *LogicID) Id() uint64 {
	return self.id
}

func GetFuncIdByLogicId(logicId uint64) (uint64, error) {
	lgc, err := NewLogicIDByID(logicId)
	if err != nil {
		return 0, err
	}
	return lgc.FuncId(), nil
}

func GetCataIdByLogicId(logicId uint64) (uint64, error) {
	lgc, err := NewLogicIDByID(logicId)
	if err != nil {
		return 0, err
	}
	return lgc.CataId(), nil
}

func GetSubIdByLogicId(logicId uint64) (uint64, error) {
	lgc, err := NewLogicIDByID(logicId)
	if err != nil {
		return 0, err
	}
	return lgc.FuncId(), nil
}
