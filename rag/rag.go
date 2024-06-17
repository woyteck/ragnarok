package rag

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"woyteck.pl/ragnarok/db"
	"woyteck.pl/ragnarok/openai"
	"woyteck.pl/ragnarok/prompter"
	"woyteck.pl/ragnarok/types"
	"woyteck.pl/ragnarok/vectordb"
)

type Rag struct {
	llm *openai.Client
	// messagesStore db.MessagesStore
	store    *db.Store
	pr       *prompter.Prompter
	vectorDB *vectordb.QdrantClient
}

func New(llm *openai.Client, store *db.Store, pr *prompter.Prompter, vectorDB *vectordb.QdrantClient) *Rag {
	return &Rag{
		llm:      llm,
		store:    store,
		pr:       pr,
		vectorDB: vectorDB,
	}
}

func (r *Rag) Ask(ctx context.Context, conversationId uuid.UUID, userPrompt string) (*types.Conversation, error) {
	conversation, err := r.store.Conversation.GetConversationByUUID(ctx, conversationId)
	if err != nil {
		return nil, err
	}

	messages, err := r.store.Message.GetMessagesByConversationUUID(ctx, conversationId)
	if err != nil {
		return nil, err
	}
	conversation.Messages = messages

	//add basic conversation context
	if len(conversation.Messages) == 0 {
		convContext, err := r.pr.Get("main_context")
		if err != nil {
			return nil, err
		}

		err = r.addMessage(ctx, conversation, "system", convContext)
		if err != nil {
			return nil, err
		}
	}

	//search the vector database and add additional context
	additionalContext, err := r.findMoreContext(ctx, userPrompt)
	if err != nil {
		return nil, err
	}
	if additionalContext != "" {
		err = r.addMessage(ctx, conversation, "system", additionalContext)
		if err != nil {
			return nil, err
		}
	}

	//add user's prompt
	err = r.addMessage(ctx, conversation, "user", userPrompt)
	if err != nil {
		return nil, err
	}

	//build request
	req := r.buildRequest(conversation)

	//send request to llm
	resp, err := r.llm.GetCompletion(req)
	if err != nil {
		return nil, err
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("LLM returned 0 completions")
	}

	answer := resp.Choices[0].Message.Content

	err = r.addMessage(ctx, conversation, "assistant", answer)
	if err != nil {
		return nil, err
	}

	return conversation, nil
}

func (r *Rag) findMoreContext(ctx context.Context, text string) (string, error) {
	vector, err := r.llm.GetEmbedding(text, "text-embedding-ada-002")
	if err != nil {
		return "", err
	}

	searchResults, err := r.vectorDB.Search("memory", vector, 5)
	if err != nil {
		return "", err
	}

	contexts := []string{}
	for _, searchResult := range searchResults {
		if searchResult.Score > 0.9 {
			fragment, err := r.store.MemoryFragment.GetMemoryFragmentByUUID(ctx, searchResult.ID)
			if err != nil {
				return "", err
			}
			contexts = append(contexts, fragment.ContentRefined)
		}
	}

	return strings.Join(contexts, "\n"), nil
}

func (r *Rag) buildRequest(conversation *types.Conversation) openai.CompletionRequest {
	reqMessages := []*openai.Message{}
	for _, message := range conversation.Messages {
		reqMessages = append(reqMessages, &openai.Message{
			Role:    message.Role,
			Content: message.Content,
		})
	}
	req := openai.CompletionRequest{
		Messages: reqMessages,
	}
	req.Model = "gpt-4-turbo"

	return req
}

func (r *Rag) addMessage(ctx context.Context, conversation *types.Conversation, role string, text string) error {
	contextMessage := types.NewMessage(conversation.ID, role, text)
	conversation.Messages = append(conversation.Messages, contextMessage)
	err := r.store.Message.InsertMessage(ctx, contextMessage)
	if err != nil {
		return err
	}

	return nil
}
