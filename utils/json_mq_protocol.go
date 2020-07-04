package utils

// msg header in message queue
type JsonMq struct {
	Ts       int64       `json:"ts"`       // when
	EventId  uint32      `json:"eventId"`  // what
	Guid     uint64      `json:"guid"`     // who
	Source   uint64      `json:"source"`   // where
	Progress int32       `json:"progress"` // progress see enums -> ActionProgress
	Tip      string      `json:"tip"`      // advice how to handle this
	Details  interface{} `json:"details"`  // details
}

