package store

import (
	"sync"
)

type Request struct {
	FileId     *string
	ReplyMsgId *int64
	Language   *string
}

type RequestStore struct {
	chats map[int64]Request
	mux   sync.RWMutex
}

func NewRequestStore() *RequestStore {
	return &RequestStore{
		chats: make(map[int64]Request),
		mux:   sync.RWMutex{},
	}
}

func (m *RequestStore) Get(chatId int64) *Request {
	m.mux.RLock()
	defer m.mux.RUnlock()
	chat, ok := m.chats[chatId]
	if !ok {
		return nil
	}
	return &chat
}

func (m *RequestStore) Set(
	chatId int64,
	data Request,
) {
	m.mux.Lock()
	m.chats[chatId] = data
	m.mux.Unlock()
}

func (m *RequestStore) Delete(chatId int64) {
	m.mux.Lock()
	delete(m.chats, chatId)
	m.mux.Unlock()
}
