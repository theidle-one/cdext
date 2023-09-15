package systemLog

import (
	"time"
)

type LicenceSystemLog struct {
	LicenceID string `bson:"licenceID" json:"licenceID"`
	SystemLog `bson:",inline"`
}

const (
	ACTION_UPDATE   = "update"
	ACTION_DELETE   = "delete"
	ACTION_ADD      = "add"
	ACTION_DOWNLOAD = "download"
	ACTION_UPLOAD   = "upload"
	ACTION_LOGIN    = "login"
	ACTION_GET      = "get"
)

type SystemLog struct {
	Username  string      `bson:"username" json:"username"`
	Service   string      `bson:"service" json:"service"`
	Action    string      `bson:"action" json:"action"`
	Resource  string      `bson:"resource" json:"resource"`
	CreatedAt time.Time   `bson:"createdAt" json:"createdAt"`
	Body      interface{} `bson:"body" json:"body"`
}
