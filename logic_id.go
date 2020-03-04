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

/* 设计规则
每个逻辑服务器都有一个ID,叫 logicId,对应一个逻辑服务器"实例"
logicId = funcId(1-999) + cataId(1-99) + subId(1-99) + instId(1-99)
其中 funcId 代表一个逻辑功能 比如斗地主游戏
cataId 代表一个分类,比如中级场
subId 子分类, 比如癞子
instId 实例, 比如癞子斗地主中级场启动了 2 个实例用于负载均衡,可能 instId 分别对应1,2
gameId 对应一个游戏,忽略 instId 由其他部分组成

实例:
1 斗地主, 1 初级场  1 普通斗地主  两个实例分别为 1, 2
则这两个进程的ID 分别为 1010101, 1010102
gameId 都为 1010100
*/
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

func (self *LogicID) GameId() uint64 {
	return self.id % LI_MAX_ID
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

// sugar functions
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

func GetGameIdByLogicId(logicId uint64) (uint64, error) {
	lgc, err := NewLogicIDByID(logicId)
	if err != nil {
		return 0, err
	}
	return lgc.GameId(), nil
}
