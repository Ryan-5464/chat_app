package connservice

import (
	i "server/interfaces"
	typ "server/types"
	"server/util"
)

func NewConnectionService() *ConnectionService {
	return &ConnectionService{
		pool:   make(map[typ.UserId]i.Socket),
		status: make(map[typ.UserId]string),
	}
}

type ConnectionService struct {
	pool   map[typ.UserId]i.Socket
	status map[typ.UserId]string
}

func (c *ConnectionService) StoreConnection(conn i.Socket, userId typ.UserId) {
	util.Log.FunctionInfo()
	util.Log.Dbugf("New User Connection: Id = %v", userId.String())
	c.pool[userId] = conn
	if c.status[userId] == "" {
		c.status[userId] = "Online"
	}
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

func (c *ConnectionService) ChangeOnlineStatus(status string, userId typ.UserId) {
	c.status[userId] = status
}

func (c *ConnectionService) GetOnlineStatus(userId typ.UserId) string {
	return c.status[userId]
}
