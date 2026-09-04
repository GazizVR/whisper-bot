package persist

import "sync"

type Chat struct {
	MediaPath *string
	Language  *string
}

type ChatManager struct {
	chats map[int64]Chat
	mux   sync.RWMutex
}

func NewChatManager() *ChatManager {
	return &ChatManager{
		chats: make(map[int64]Chat),
		mux:   sync.RWMutex{},
	}
}

func (m *ChatManager) GetChat(chatId int64) *Chat {
	m.mux.RLock()
	defer m.mux.RUnlock()
	chat, ok := m.chats[chatId]
	if !ok {
		return nil
	}
	return &chat
}

func (m *ChatManager) PutChat(
	chatId int64,
	data Chat,
) {
	m.mux.Lock()
	m.chats[chatId] = data
	m.mux.Unlock()
}

func (m *ChatManager) RemoveChat(chatId int64) {
	m.mux.Lock()
	delete(m.chats, chatId)
	m.mux.Unlock()
}
