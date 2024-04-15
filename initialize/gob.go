package initialize

import (
	"encoding/gob"
	g "odisk/global"
	m "odisk/model"
)
func InitGob()  {
	gob.Register(m.UserInfo{})
	gob.Register(g.EmailData{})
	gob.Register(g.Message{})
	gob.Register(g.SendEmailData{})
	gob.Register(g.SimpleData{})
}