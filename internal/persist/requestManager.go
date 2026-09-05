package persist

import "sync"

type Request struct {
	FileId     *string
	ReplyMsgId *int64
	Language   *string
}

type RequestManager struct {
	chats map[int64]Request
	mux   sync.RWMutex
}

func NewRequestManager() *RequestManager {
	return &RequestManager{
		chats: make(map[int64]Request),
		mux:   sync.RWMutex{},
	}
}

func (m *RequestManager) GetRequest(chatId int64) *Request {
	m.mux.RLock()
	defer m.mux.RUnlock()
	chat, ok := m.chats[chatId]
	if !ok {
		return nil
	}
	return &chat
}

func (m *RequestManager) PutRequest(
	chatId int64,
	data Request,
) {
	m.mux.Lock()
	m.chats[chatId] = data
	m.mux.Unlock()
}

func (m *RequestManager) RemoveRequest(chatId int64) {
	m.mux.Lock()
	delete(m.chats, chatId)
	m.mux.Unlock()
}
