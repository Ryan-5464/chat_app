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
	c.lgr.Log("User connected; UserId -- " + string(userId))
	c.pool[userId] = conn
}

func (c *ConnectionService) GetConnection(userId typ.UserId) i.Socket {
	c.lgr.LogFunctionInfo()
	c.lgr.Log("User disconnected; UserId -- " + string(userId))
	return c.pool[userId]
}
