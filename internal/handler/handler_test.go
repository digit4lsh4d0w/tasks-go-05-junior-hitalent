package handler_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"task-5/internal/handler"
	"task-5/internal/model"

	"github.com/go-playground/validator"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockChatService struct {
	mock.Mock
}

func (m *MockChatService) CreateChat(ctx context.Context, chat *model.Chat) error {
	args := m.Called(ctx, chat)
	return args.Error(0)
}

func (m *MockChatService) DeleteChat(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockChatService) CreateMessage(ctx context.Context, msg *model.Message) error {
	args := m.Called(ctx, msg)
	return args.Error(0)
}

func (m *MockChatService) GetChatWithMessages(ctx context.Context, id uint, limit int) (*model.Chat, error) {
	args := m.Called(ctx, id, limit)

	var chat *model.Chat
	if args.Get(0) != nil {
		chat = args.Get(0).(*model.Chat)
	}

	return chat, args.Error(1)
}

type testEnv struct {
	service   *MockChatService
	logger    *slog.Logger
	validator *validator.Validate
}

func setupTestEnv(mockSetup func(*MockChatService)) testEnv {
	service := new(MockChatService)
	if mockSetup != nil {
		mockSetup(service)
	}

	return testEnv{
		service:   service,
		logger:    slog.New(slog.NewTextHandler(io.Discard, nil)),
		validator: validator.New(),
	}
}

func TestChatHandler_CreateChat(t *testing.T) {
	tests := []struct {
		name           string
		body           string
		mockSetup      func(m *MockChatService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Broken JSON",
			body:           `{"title": "Chat title"`,
			mockSetup:      func(m *MockChatService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Empty chat title",
			body:           `{"title": ""}`,
			mockSetup:      func(m *MockChatService) {},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:           "Chat title too long",
			body:           fmt.Sprintf(`{"title": "%s"}`, strings.Repeat("a", 201)),
			mockSetup:      func(m *MockChatService) {},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:           "Chat title consists entirely of spaces",
			body:           `{"title": "    "}`,
			mockSetup:      func(m *MockChatService) {},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "Chat duplicate",
			body: `{"title": "Chat title"}`,
			mockSetup: func(m *MockChatService) {
				m.On(
					"CreateChat",
					mock.Anything,
					mock.MatchedBy(func(chat *model.Chat) bool {
						return chat.Title == "Chat title"
					}),
				).Return(model.ErrAlreadyExists).Once()
			},
			expectedStatus: http.StatusConflict,
		},
		{
			name: "Success",
			body: `{"title": "Chat title"}`,
			mockSetup: func(m *MockChatService) {
				m.On(
					"CreateChat",
					mock.Anything,
					mock.MatchedBy(func(chat *model.Chat) bool {
						return chat.Title == "Chat title"
					}),
				).Run(func(args mock.Arguments) {
					chat := args.Get(1).(*model.Chat)
					chat.ID = 1
					chat.CreatedAt = time.Date(2026, 4, 20, 0, 0, 0, 0, time.UTC)
				}).Return(nil).Once()
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"id": 1, "title": "Chat title", "created_at": "2026-04-20T00:00:00Z"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			env := setupTestEnv(tc.mockSetup)
			h := handler.NewChatHandler(env.service, env.validator, env.logger)

			mux := http.NewServeMux()
			mux.HandleFunc("POST /chat", h.CreateChat)

			path := "/chat"
			request := httptest.NewRequest(http.MethodPost, path, strings.NewReader(tc.body))
			request.Header.Set("Content-Type", "application/json")
			recorder := httptest.NewRecorder()

			mux.ServeHTTP(recorder, request)

			require.Equal(t, tc.expectedStatus, recorder.Code)

			var response map[string]any
			err := json.Unmarshal(recorder.Body.Bytes(), &response)
			require.NoError(t, err, "response must be valid JSON")
			if tc.expectedStatus == http.StatusCreated {
				require.JSONEq(t, tc.expectedBody, recorder.Body.String())
			} else {
				require.Contains(t, response, "error", "response must contain \"error\" field")
				require.NotEmpty(t, response["error"])
			}

			env.service.AssertExpectations(t)
		})
	}
}

func TestChatHandler_DeleteChat(t *testing.T) {
	tests := []struct {
		name           string
		chatIDPath     string
		mockSetup      func(m *MockChatService)
		expectedStatus int
	}{
		{
			name:           "Invalid ID",
			chatIDPath:     "abc",
			mockSetup:      func(m *MockChatService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:       "Not found",
			chatIDPath: "2",
			mockSetup: func(m *MockChatService) {
				m.On(
					"DeleteChat",
					mock.Anything,
					uint(2),
				).Return(model.ErrNotFound).Once()
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:       "Success",
			chatIDPath: "1",
			mockSetup: func(m *MockChatService) {
				m.On(
					"DeleteChat",
					mock.Anything,
					uint(1),
				).Return(nil).Once()
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			env := setupTestEnv(tc.mockSetup)
			h := handler.NewChatHandler(env.service, env.validator, env.logger)

			mux := http.NewServeMux()
			mux.HandleFunc("DELETE /chat/{chat_id}", h.DeleteChat)

			path := "/chat/" + tc.chatIDPath
			request := httptest.NewRequest(http.MethodDelete, path, nil)
			recorder := httptest.NewRecorder()

			mux.ServeHTTP(recorder, request)

			require.Equal(t, tc.expectedStatus, recorder.Code)

			var response map[string]any
			err := json.Unmarshal(recorder.Body.Bytes(), &response)
			require.NoError(t, err, "response must be valid JSON")
			if tc.expectedStatus == http.StatusOK {
				require.Contains(t, response, "success", "response must contain \"success\" field")
				require.NotEmpty(t, response["success"])
			} else {
				require.Contains(t, response, "error", "response must contain \"error\" field")
				require.NotEmpty(t, response["error"])
			}

			env.service.AssertExpectations(t)
		})
	}
}

func TestChatHandler_CreateMessage(t *testing.T) {
	tests := []struct {
		name           string
		chatIDPath     string
		body           string
		mockSetup      func(m *MockChatService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Invalid ID",
			chatIDPath:     "abc",
			mockSetup:      func(m *MockChatService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Broken JSON",
			chatIDPath:     "1",
			body:           `{"text": ""`,
			mockSetup:      func(m *MockChatService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Empty message text",
			chatIDPath:     "1",
			body:           `{"text": ""}`,
			mockSetup:      func(m *MockChatService) {},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:           "Message text too long",
			chatIDPath:     "1",
			body:           fmt.Sprintf(`{"text": "%s"}`, strings.Repeat("a", 5001)),
			mockSetup:      func(m *MockChatService) {},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:           "Message text consists entirely of spaces",
			chatIDPath:     "1",
			body:           `{"text": "    "}`,
			mockSetup:      func(m *MockChatService) {},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:       "Not found",
			chatIDPath: "2",
			body:       `{"text": "Message text"}`,
			mockSetup: func(m *MockChatService) {
				m.On(
					"CreateMessage",
					mock.Anything,
					mock.MatchedBy(func(msg *model.Message) bool {
						return msg.ChatID == 2 && msg.Text == "Message text"
					}),
				).Return(model.ErrNotFound).Once()
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:       "Success",
			chatIDPath: "1",
			body:       `{"text": "Message text"}`,
			mockSetup: func(m *MockChatService) {
				m.On(
					"CreateMessage",
					mock.Anything,
					mock.MatchedBy(func(msg *model.Message) bool {
						return msg.ChatID == 1 && msg.Text == "Message text"
					}),
				).Run(func(args mock.Arguments) {
					msg := args.Get(1).(*model.Message)
					msg.ID = 1
					msg.ChatID = 1
					msg.Text = "Message text"
					msg.CreatedAt = time.Date(2026, 4, 20, 0, 0, 0, 0, time.UTC)
				}).Return(nil).Once()
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"id": 1, "chat_id": 1, "text": "Message text", "created_at": "2026-04-20T00:00:00Z"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			env := setupTestEnv(tc.mockSetup)
			h := handler.NewChatHandler(env.service, env.validator, env.logger)

			mux := http.NewServeMux()
			mux.HandleFunc("POST /chat/{chat_id}/message", h.CreateMessage)

			path := "/chat/" + tc.chatIDPath + "/message"
			request := httptest.NewRequest(http.MethodPost, path, strings.NewReader(tc.body))
			request.Header.Set("Content-Type", "application/json")
			recorder := httptest.NewRecorder()

			mux.ServeHTTP(recorder, request)

			require.Equal(t, tc.expectedStatus, recorder.Code)

			var response map[string]any
			err := json.Unmarshal(recorder.Body.Bytes(), &response)
			require.NoError(t, err, "response must be valid JSON")
			if tc.expectedStatus == http.StatusCreated {
				require.JSONEq(t, tc.expectedBody, recorder.Body.String())
			} else {
				require.Contains(t, response, "error", "response must contain \"error\" field")
				require.NotEmpty(t, response["error"])
			}

			env.service.AssertExpectations(t)
		})
	}
}

func TestChatHandler_GetAllMessages(t *testing.T) {
	tests := []struct {
		name           string
		chatIDPath     string
		limitQuery     string
		mockSetup      func(m *MockChatService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Invalid ID",
			chatIDPath:     "abc",
			limitQuery:     "limit=20",
			mockSetup:      func(m *MockChatService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:       "Invalid limit",
			chatIDPath: "1",
			limitQuery: "limit=abc",
			mockSetup: func(m *MockChatService) {
				m.On(
					"GetChatWithMessages",
					mock.Anything,
					uint(1),
					5,
				).Return(&model.Chat{
					ID:        1,
					Title:     "Chat title",
					Messages:  []model.Message{},
					CreatedAt: time.Date(2026, 4, 20, 0, 0, 0, 0, time.UTC),
				}, nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id": 1, "title": "Chat title", "messages": [], "created_at": "2026-04-20T00:00:00Z"}`,
		},
		{
			name:       "Not Found",
			chatIDPath: "2",
			limitQuery: "limit=20",
			mockSetup: func(m *MockChatService) {
				m.On("GetChatWithMessages", mock.Anything, uint(2), 20).Return(nil, model.ErrNotFound).Once()
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:       "Success without messages",
			chatIDPath: "1",
			limitQuery: "limit=20",
			mockSetup: func(m *MockChatService) {
				m.On(
					"GetChatWithMessages",
					mock.Anything,
					uint(1),
					20,
				).Return(&model.Chat{
					ID:        1,
					Title:     "Chat title",
					Messages:  []model.Message{},
					CreatedAt: time.Date(2026, 4, 20, 0, 0, 0, 0, time.UTC),
				}, nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id": 1, "title": "Chat title", "messages": [], "created_at": "2026-04-20T00:00:00Z"}`,
		},
		{
			name:       "Success with messages",
			chatIDPath: "1",
			limitQuery: "limit=20",
			mockSetup: func(m *MockChatService) {
				m.On(
					"GetChatWithMessages",
					mock.Anything,
					uint(1),
					20,
				).Return(&model.Chat{
					ID:    1,
					Title: "Chat title",
					Messages: []model.Message{
						{ID: 1, ChatID: 1, Text: "Message text", CreatedAt: time.Date(2026, 4, 20, 0, 0, 0, 0, time.UTC)},
					},
					CreatedAt: time.Date(2026, 4, 20, 0, 0, 0, 0, time.UTC),
				}, nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id": 1, "title": "Chat title", "messages": [{"id": 1, "chat_id": 1, "text": "Message text", "created_at": "2026-04-20T00:00:00Z"}], "created_at": "2026-04-20T00:00:00Z"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			env := setupTestEnv(tc.mockSetup)
			h := handler.NewChatHandler(env.service, env.validator, env.logger)

			mux := http.NewServeMux()
			mux.HandleFunc("GET /chat/{chat_id}", h.GetAllMessages)

			path := "/chat/" + tc.chatIDPath
			if tc.limitQuery != "" {
				path += "?" + tc.limitQuery
			}
			request := httptest.NewRequest(http.MethodGet, path, nil)
			recorder := httptest.NewRecorder()

			mux.ServeHTTP(recorder, request)

			require.Equal(t, tc.expectedStatus, recorder.Code)

			var response map[string]any
			err := json.Unmarshal(recorder.Body.Bytes(), &response)
			require.NoError(t, err, "response must be valid JSON")
			if tc.expectedStatus == http.StatusOK {
				require.JSONEq(t, tc.expectedBody, recorder.Body.String())
			} else {
				require.Contains(t, response, "error", "response must contain \"error\" field")
				require.NotEmpty(t, response["error"])
			}

			env.service.AssertExpectations(t)
		})
	}
}
