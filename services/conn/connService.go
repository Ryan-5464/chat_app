package connservice

import (
	i "server/interfaces"
	typ "server/types"
	"server/util"
)

func NewConnectionService() *ConnectionService {
	return &ConnectionService{
		pool: make(map[typ.UserId]i.Socket),
	}
}

type ConnectionService struct {
	pool map[typ.UserId]i.Socket
}

func (c *ConnectionService) StoreConnection(conn i.Socket, userId typ.UserId) {
	util.Log.FunctionInfo()
	util.Log.Dbugf("New User Connection: Id = %v", userId.String())
	c.pool[userId] = conn
}

func (c *ConnectionService) GetConnection(userId typ.UserId) i.Socket {
	util.Log.FunctionInfo()
	return c.pool[userId]
}

func (c *ConnectionService) DisconnectUser(userId typ.UserId) {
	util.Log.FunctionInfo()
	delete(c.pool, userId)
	util.Log.Dbugf("User disconnected: Id = %v", userId.String())
}

func (c *ConnectionService) GetActiveConnections() map[typ.UserId]i.Socket {
	return c.pool
}
