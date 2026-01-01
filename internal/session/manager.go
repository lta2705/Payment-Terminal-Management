package session

import (
	"fmt"
	"github.com/bytedance/gopkg/util/logger"
	"net"
	"sync"
)

type SessionManager struct {
	sessions sync.Map
}

func NewSessionManager() *SessionManager {
	return &SessionManager{}
}

func (sm *SessionManager) Add(conn net.Conn, trmId string) {

	sm.sessions.Store(trmId, conn)
	logger.Info("Session added", trmId)
}

func (sm *SessionManager) Remove(addr string) {
	sm.sessions.Delete(addr)
}

func (sm *SessionManager) Count() int {
	count := 0
	sm.sessions.Range(func(_, _ interface{}) bool {
		count++
		return true
	})
	return count
}

func (sm *SessionManager) Check(trmId string) bool {
	_, exists := sm.sessions.Load(trmId)
	return exists
}

func (sm *SessionManager) Send(trmId string, message string) bool {
	value, exists := sm.sessions.Load(trmId)
	if !exists {
		logger.Info("Terminal ID not found:")
		return false
	}

	conn, ok := value.(net.Conn)
	if !ok {
		logger.Info("Session is not a net conn:")
		return false
	}
	_, err := conn.Write([]byte(message))
	if err != nil {
		logger.Error("Error in sending message to terminal", err)
		return false
	}

	return true
}

// Sends to all clients
func (sm *SessionManager) Broadcast(message string) {
	sm.sessions.Range(func(key, value interface{}) bool {
		if conn, ok := value.(net.Conn); ok {
			conn.Write([]byte(fmt.Sprintf("[BROADCAST]: %s\n", message)))
		}
		return true
	})
}

// Close all connections
func (sm *SessionManager) CloseAll() {
	sm.sessions.Range(func(key, value interface{}) bool {
		if conn, ok := value.(net.Conn); ok {
			conn.Close()
		}
		return true
	})
}
