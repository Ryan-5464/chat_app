package connservice

import (
	i "server/interfaces"
	typ "server/types"
)

func NewConnectionService(lgr i.Logger) *ConnectionService {
	return &ConnectionService{
		lgr:  lgr,
		pool: make(map[typ.UserId]i.Socket),
	}
}

type ConnectionService struct {
	lgr  i.Logger
	pool map[typ.UserId]i.Socket
}

func (c *ConnectionService) StoreConnection(conn i.Socket, userId typ.UserId) {
	c.lgr.LogFunctionInfo()
	msgNewUserConnection := "New User Connection: Id = " + userId.String()
	c.lgr.DLog(msgNewUserConnection)
	c.pool[userId] = conn
}

func (c *ConnectionService) GetConnection(userId typ.UserId) i.Socket {
	c.lgr.LogFunctionInfo()
	return c.pool[userId]
}

func (c *ConnectionService) DisconnectUser(userId typ.UserId) {
	c.lgr.LogFunctionInfo()
	delete(c.pool, userId)
	msgUserDisconnected := "User disconnected: Id = " + userId.String()
	c.lgr.Log(msgUserDisconnected)
}
