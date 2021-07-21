package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"workshops/rest-api/internal/entities"
	"workshops/rest-api/internal/errors"
	"workshops/rest-api/internal/filters"
	"workshops/rest-api/internal/validators"
)

type fetchTestCase struct {
	name                 string
	inputBody            string
	mockBehavior         func(m *MockTaskService, filters *filters.TaskFilter)
	expectedStatus       int
	expectedBodyResponse Response
}

func TestTaskFetchMock(t *testing.T) {

	cases := []fetchTestCase{
		{
			name:      "Successful fetch",
			inputBody: "",
			mockBehavior: func(m *MockTaskService, filters *filters.TaskFilter) {
				m.EXPECT().Fetch(context.Background(), filters)
			},
			expectedStatus:       http.StatusOK,
			expectedBodyResponse: Response{},
		},
		{
			name:      "Error fetch",
			inputBody: "",
			mockBehavior: func(m *MockTaskService, filters *filters.TaskFilter) {
				m.EXPECT().Fetch(context.Background(), filters).Return(nil, fmt.Errorf("Fetch Error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBodyResponse: Response{
				Error: fmt.Sprint("Fetch Error"),
			},
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Fetch Case %d: %s ", i, c.name), func(t *testing.T) {
			request, _ := http.NewRequest(http.MethodGet, "/tasks?category", nil)
			response := httptest.NewRecorder()
			ctrl := gomock.NewController(t)
			service := NewMockTaskService(ctrl)
			taskController := NewTask(service)
			validator := validators.TaskValidator{}
			queryFilters := filters.ValidatedTaskFilter(&validator, request.URL.Query())
			c.mockBehavior(service, &queryFilters)

			taskController.Fetch(response, request)

			got := Response{}
			dec := json.NewDecoder(response.Body)
			_ = dec.Decode(&got)
			want := c.expectedBodyResponse
			assert.Equal(t, response.Code, c.expectedStatus)
			assert.Equal(t, got, want, fmt.Sprintf("got %q, want %q", got, want))
		})
	}
}

type getTestCases struct {
	name                 string
	inputBody            string
	vars                 map[string]string
	mockBehavior         func(m *MockTaskService)
	expectedStatus       int
	expectedBodyResponse Response
}

func TestTaskGetMock(t *testing.T) {
	cases := []getTestCases{
		{
			name:      "Without error",
			inputBody: "",
			mockBehavior: func(m *MockTaskService) {
				m.EXPECT().Get(gomock.Any(), 1)
			},
			vars:                 map[string]string{"id": "1"},
			expectedStatus:       200,
			expectedBodyResponse: Response{},
		},
		{
			name:      "With atoi error",
			inputBody: "",
			mockBehavior: func(m *MockTaskService) {
				//EXPECT error before mock run
			},
			expectedStatus: http.StatusBadRequest,
			expectedBodyResponse: Response{
				Error: "strconv.Atoi: parsing \"\": invalid syntax",
			},
		},
		{
			name:      "With not found error",
			inputBody: "",
			mockBehavior: func(m *MockTaskService) {
				m.EXPECT().Get(gomock.Any(), 1).Return(nil, errors.NotFound)
			},
			vars:           map[string]string{"id": "1"},
			expectedStatus: 404,
			expectedBodyResponse: Response{
				Error: errors.NotFound.Error(),
			},
		},
	}
	for i, c := range cases {
		t.Run(fmt.Sprintf("Get Case %d: %s ", i, c.name), func(t *testing.T) {
			request, _ := http.NewRequest(http.MethodGet, "/tasks/1", bytes.NewBufferString(c.inputBody))
			request = mux.SetURLVars(request, c.vars)

			response := httptest.NewRecorder()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := NewMockTaskService(ctrl)
			c.mockBehavior(service)

			taskController := NewTask(service)
			taskController.Get(response, request)

			got := Response{}
			dec := json.NewDecoder(response.Body)
			_ = dec.Decode(&got)
			want := c.expectedBodyResponse

			assert.Equal(t, c.expectedStatus, response.Code)
			assert.Equal(t, got, want, fmt.Sprintf("got %q, want %q", got, want))
		})
	}
}

type createTestCases struct {
	name                 string
	inputBody            *entities.Task
	mockBehavior         func(m *MockTaskService)
	expectedStatus       int
	expectedBodyResponse string
}

func TestTaskCreateMock(t *testing.T) {
	task := entities.Task{
		Title:       "New task",
		Description: "Test description",
		Category:    validators.NOTE,
		Date:        time.Time{},
	}

	cases := []createTestCases{
		{
			name:      "Without error",
			inputBody: &task,
			mockBehavior: func(m *MockTaskService) {
				m.EXPECT().Create(gomock.Any(), &task).Return(&task, nil)
			},
			expectedStatus:       http.StatusOK,
			expectedBodyResponse: `{"payload":{"Id":0,"Title":"New task","Description":"Test description","Category":"note","Date":"0001-01-01T00:00:00Z"}}`,
		},
		{
			name:      "Bad request error",
			inputBody: nil,
			mockBehavior: func(m *MockTaskService) {

			},
			expectedStatus:       http.StatusBadRequest,
			expectedBodyResponse: `{"error":"Wrong data for new task"}`,
		},
		{
			name:      "Server error",
			inputBody: &task,
			mockBehavior: func(m *MockTaskService) {
				m.EXPECT().Create(gomock.Any(), &task).Return(nil, fmt.Errorf("error"))
			},
			expectedStatus:       http.StatusInternalServerError,
			expectedBodyResponse: `{"error":"error"}`,
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Create Case %d: %s ", i, c.name), func(t *testing.T) {
			data := new(bytes.Buffer)
			if c.inputBody != nil {
				json.NewEncoder(data).Encode(c.inputBody)
			} else {
				json.NewEncoder(data).Encode([]byte{})
			}

			request, _ := http.NewRequest(http.MethodPost, "/tasks", data)
			response := httptest.NewRecorder()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := NewMockTaskService(ctrl)
			c.mockBehavior(service)

			taskController := NewTask(service)
			taskController.Create(response, request)

			got := response.Body.String()
			want := c.expectedBodyResponse

			assert.Equal(t, c.expectedStatus, response.Code)
			assert.Equal(t, got, want, fmt.Sprintf("got %q, want %q", got, want))
		})
	}
}

type updateTestCases struct {
	name                 string
	inputBody            *entities.Task
	vars                 map[string]string
	mockBehavior         func(m *MockTaskService)
	expectedStatus       int
	expectedBodyResponse string
}

func TestTaskUpdateMock(t *testing.T) {
	task := entities.Task{
		Title:       "New task",
		Description: "Test description",
		Category:    validators.NOTE,
		Date:        time.Time{},
	}

	cases := []updateTestCases{
		{
			name:      "Without error",
			inputBody: &task,
			vars:      map[string]string{"id": "1"},
			mockBehavior: func(m *MockTaskService) {
				m.EXPECT().Update(gomock.Any(), 1, gomock.Any()).Return(&task, nil)
			},
			expectedStatus:       http.StatusOK,
			expectedBodyResponse: `{"payload":{"Id":0,"Title":"New task","Description":"Test description","Category":"note","Date":"0001-01-01T00:00:00Z"}}`,
		},
		{
			name:      "Bad request atoi error",
			inputBody: nil,
			vars:      map[string]string{"id": "r"},
			mockBehavior: func(m *MockTaskService) {

			},
			expectedStatus:       http.StatusBadRequest,
			expectedBodyResponse: `{"error":"Wrong Id parameter"}`,
		},
		{
			name:      "Bad request error",
			inputBody: nil,
			vars:      map[string]string{"id": "1"},
			mockBehavior: func(m *MockTaskService) {

			},
			expectedStatus:       http.StatusBadRequest,
			expectedBodyResponse: `{"error":"Cant decode task data"}`,
		},
		{
			name:      "Server error",
			vars:      map[string]string{"id": "1"},
			inputBody: &task,
			mockBehavior: func(m *MockTaskService) {
				m.EXPECT().Update(gomock.Any(), 1, gomock.Any()).Return(nil, fmt.Errorf("Task didnt update"))
			},
			expectedStatus:       http.StatusInternalServerError,
			expectedBodyResponse: `{"error":"Task didnt update"}`,
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Update Case %d: %s ", i, c.name), func(t *testing.T) {
			data := new(bytes.Buffer)
			if c.inputBody != nil {
				json.NewEncoder(data).Encode(c.inputBody)
			} else {
				json.NewEncoder(data).Encode([]byte{})
			}

			request, _ := http.NewRequest(http.MethodPost, "/tasks", data)
			request = mux.SetURLVars(request, c.vars)
			response := httptest.NewRecorder()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := NewMockTaskService(ctrl)
			c.mockBehavior(service)

			taskController := NewTask(service)
			taskController.Update(response, request)

			got := response.Body.String()
			want := c.expectedBodyResponse

			assert.Equal(t, c.expectedStatus, response.Code)
			assert.Equal(t, got, want, fmt.Sprintf("got %q, want %q", got, want))
		})
	}
}

func TestTaskDeleteMock(t *testing.T) {
	cases := []getTestCases{
		{
			name:      "Without error",
			inputBody: "",
			mockBehavior: func(m *MockTaskService) {
				m.EXPECT().Delete(gomock.Any(), 1).Return(true, nil)
			},
			vars:           map[string]string{"id": "1"},
			expectedStatus: 200,
			expectedBodyResponse: Response{
				Payload: true,
			},
		},
		{
			name:      "With atoi error",
			inputBody: "",
			mockBehavior: func(m *MockTaskService) {
				//EXPECT error before mock run
			},
			expectedStatus: http.StatusBadRequest,
			expectedBodyResponse: Response{
				Error: "Wrong Id parameter",
			},
		},
		{
			name:      "With Task didnt delete error",
			inputBody: "",
			mockBehavior: func(m *MockTaskService) {
				m.EXPECT().Delete(gomock.Any(), 1).Return(false, fmt.Errorf("Task didnt delete"))
			},
			vars:           map[string]string{"id": "1"},
			expectedStatus: 500,
			expectedBodyResponse: Response{
				Error: "Task didnt delete",
			},
		},
	}
	for i, c := range cases {
		t.Run(fmt.Sprintf("Delete Case %d: %s ", i, c.name), func(t *testing.T) {
			request, _ := http.NewRequest(http.MethodGet, "/tasks/1", bytes.NewBufferString(c.inputBody))
			request = mux.SetURLVars(request, c.vars)

			response := httptest.NewRecorder()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := NewMockTaskService(ctrl)
			c.mockBehavior(service)

			taskController := NewTask(service)
			taskController.Delete(response, request)

			got := Response{}
			dec := json.NewDecoder(response.Body)
			_ = dec.Decode(&got)
			want := c.expectedBodyResponse

			assert.Equal(t, c.expectedStatus, response.Code)
			assert.Equal(t, got, want, fmt.Sprintf("got %q, want %q", got, want))
		})
	}
}
