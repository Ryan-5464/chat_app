package connservice

import (
	i "server/interfaces"
	typ "server/types"
	"server/util"
)

func NewConnectionService(u i.UserService) *ConnectionService {
	return &ConnectionService{
		pool:   make(map[typ.UserId]i.Socket),
		status: make(map[typ.UserId]string),
		userS:  u,
	}
}

type ConnectionService struct {
	pool   map[typ.UserId]i.Socket
	status map[typ.UserId]string
	userS  i.UserService
}

func (c *ConnectionService) SetUserService(userS i.UserService) {
	c.userS = userS
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

func (c *ConnectionService) ChangeOnlineStatus(status string, userId typ.UserId) error {
	c.status[userId] = status

	contacts, err := c.userS.GetContacts(userId)
	if err != nil {
		return err
	}

	for _, contact := range contacts {
		conn := c.GetConnection(typ.UserId(contact.Id))

		if conn == nil {
			continue
		}

		payload := struct {
			Type         int
			OnlineStatus string
			UserId       typ.UserId
		}{
			Type:         11,
			OnlineStatus: status,
			UserId:       userId,
		}

		if err := conn.WriteJSON(payload); err != nil {
			util.Log.Errorf("failed to write to websocket connection: %v", err)
			return err
		}

	}
	return nil
}

func (c *ConnectionService) GetOnlineStatus(userId typ.UserId) string {
	return c.status[userId]
}
