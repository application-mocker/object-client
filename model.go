package object_client

type BaseNode struct {
	Id       string `json:"id"`
	CreateAt int64  `json:"create_at"` // the nano-second of create time
	UpdateAt int64  `json:"update_at"` // the nano-second of last update time
	DeleteAt int64  `json:"delete_at"` // the nano-second of delete time, if delete_at less than zero, the data not delete
}

type CommonNode struct {
	BaseNode

	Value interface{} `json:"data"`
}

type SimpleNode struct {
	BaseNode

	// value here
	DataValue map[string]interface{} `json:"data"` // the custom datas
}
