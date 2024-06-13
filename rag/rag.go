package rag

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"woyteck.pl/ragnarok/db"
	"woyteck.pl/ragnarok/openai"
	"woyteck.pl/ragnarok/types"
)

type Rag struct {
	llm           *openai.Client
	messagesStore db.MessagesStore
}

func New(llm *openai.Client, messagesStore db.MessagesStore) *Rag {
	return &Rag{
		llm:           llm,
		messagesStore: messagesStore,
	}
}

func (r *Rag) Ask(conversation *types.Conversation) error {
	messages := []*openai.Message{}
	for _, message := range conversation.Messages {
		messages = append(messages, &openai.Message{
			Role:    message.Role,
			Content: message.Content,
		})
	}
	req := openai.CompletionRequest{
		Messages: messages,
	}
	req.Model = "gpt-3.5-turbo"

	resp, err := r.llm.GetCompletion(req)
	if err != nil {
		return err
	}

	if len(resp.Choices) == 0 {
		return fmt.Errorf("LLM returned 0 completions")
	}

	answer := resp.Choices[0].Message.Content
	now := time.Now()
	newMessage := &types.Message{
		ID:             uuid.New(),
		ConversationId: conversation.ID,
		Role:           "assistant",
		Content:        answer,
		CreatedAt:      &now,
	}
	conversation.Messages = append(conversation.Messages, newMessage)
	_, err = r.messagesStore.InsertMessage(context.Background(), newMessage)
	if err != nil {
		return err
	}

	return nil
}
