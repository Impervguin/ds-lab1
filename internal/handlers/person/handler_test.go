package person_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/Impervguin/ds-lab1/internal/domain"
	"github.com/Impervguin/ds-lab1/internal/handlers/person"
	persondto "github.com/Impervguin/ds-lab1/internal/handlers/person/dto"
)

func newTestRouter(repo domain.PersonRepository) http.Handler {
	r := chi.NewRouter()
	h := person.NewHandler(repo)
	h.Register(r)
	return r
}

func TestHandler_List(t *testing.T) {
	t.Run("returns persons from repository", func(t *testing.T) {
		repo := new(mockPersonRepository)
		persons := []*domain.Person{
			{ID: 1, Name: "Alice", Age: 30, Address: "Addr1", Work: "Work1"},
			{ID: 2, Name: "Bob", Age: 25, Address: "Addr2", Work: "Work2"},
		}
		repo.On("List", mock.Anything).Return(persons, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/persons/", nil)
		rec := httptest.NewRecorder()

		newTestRouter(repo).ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		var got []persondto.PersonResponse
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
		require.Len(t, got, 2)
		assert.Equal(t, "Alice", got[0].Name)
		assert.Equal(t, "Bob", got[1].Name)
		repo.AssertExpectations(t)
	})

	t.Run("returns empty list", func(t *testing.T) {
		repo := new(mockPersonRepository)
		repo.On("List", mock.Anything).Return([]*domain.Person{}, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/persons/", nil)
		rec := httptest.NewRecorder()

		newTestRouter(repo).ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		var got []persondto.PersonResponse
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
		assert.Len(t, got, 0)
	})
}

func TestHandler_Create(t *testing.T) {
	t.Run("creates person and sets location header", func(t *testing.T) {
		repo := new(mockPersonRepository)
		created := &domain.Person{ID: 42, Name: "Alice", Age: 30, Address: "Addr", Work: "Work"}
		repo.On("Create", mock.Anything, mock.MatchedBy(func(p *domain.Person) bool {
			return p.Name == "Alice" && p.Age == 30
		})).Return(created, nil)

		body := `{"name":"Alice","age":30,"address":"Addr","work":"Work"}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1/persons/", bytes.NewBufferString(body))
		rec := httptest.NewRecorder()

		newTestRouter(repo).ServeHTTP(rec, req)

		require.Equal(t, http.StatusCreated, rec.Code)
		repo.AssertExpectations(t)
	})

	t.Run("invalid body returns 400", func(t *testing.T) {
		repo := new(mockPersonRepository)

		body := `{"age":30}` // missing required name
		req := httptest.NewRequest(http.MethodPost, "/api/v1/persons/", bytes.NewBufferString(body))
		rec := httptest.NewRecorder()

		newTestRouter(repo).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		repo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	})

	t.Run("malformed json returns 400", func(t *testing.T) {
		repo := new(mockPersonRepository)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/persons/", bytes.NewBufferString("{not json"))
		rec := httptest.NewRecorder()

		newTestRouter(repo).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestHandler_Get(t *testing.T) {
	t.Run("returns person by id", func(t *testing.T) {
		repo := new(mockPersonRepository)
		p := &domain.Person{ID: 7, Name: "Alice", Age: 30, Address: "Addr", Work: "Work"}
		repo.On("GetByID", mock.Anything, int32(7)).Return(p, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/persons/7/", nil)
		rec := httptest.NewRecorder()

		newTestRouter(repo).ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		var got persondto.PersonResponse
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
		assert.Equal(t, int32(7), got.ID)
		assert.Equal(t, "Alice", got.Name)
	})

	t.Run("not found returns 404", func(t *testing.T) {
		repo := new(mockPersonRepository)
		repo.On("GetByID", mock.Anything, int32(99)).Return(nil, domain.ErrPersonNotFound)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/persons/99/", nil)
		rec := httptest.NewRecorder()

		newTestRouter(repo).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("non-numeric id returns 404", func(t *testing.T) {
		repo := new(mockPersonRepository)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/persons/abc/", nil)
		rec := httptest.NewRecorder()

		newTestRouter(repo).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		repo.AssertNotCalled(t, "GetByID", mock.Anything, mock.Anything)
	})
}

func TestHandler_Update(t *testing.T) {
	t.Run("updates person with all fields", func(t *testing.T) {
		repo := new(mockPersonRepository)
		existing := &domain.Person{ID: 7, Name: "Alice", Age: 30, Address: "Addr", Work: "Work"}
		updated := &domain.Person{ID: 7, Name: "Bob", Age: 40, Address: "NewAddr", Work: "NewWork"}
		repo.On("GetByID", mock.Anything, int32(7)).Return(existing, nil)
		repo.On("Update", mock.Anything, int32(7), mock.MatchedBy(func(p *domain.Person) bool {
			return p.Name == "Bob" && p.Age == 40 && p.Address == "NewAddr" && p.Work == "NewWork"
		})).Return(updated, nil)

		body := `{"name":"Bob","age":40,"address":"NewAddr","work":"NewWork"}`
		req := httptest.NewRequest(http.MethodPatch, "/api/v1/persons/7/", bytes.NewBufferString(body))
		rec := httptest.NewRecorder()

		newTestRouter(repo).ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		var got persondto.PersonResponse
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
		assert.Equal(t, "Bob", got.Name)
		assert.Equal(t, int32(40), got.Age)
	})

	t.Run("updates only fields present in request", func(t *testing.T) {
		repo := new(mockPersonRepository)
		existing := &domain.Person{ID: 7, Name: "Alice", Age: 30, Address: "Addr", Work: "Work"}
		updated := &domain.Person{ID: 7, Name: "Alice", Age: 45, Address: "Addr", Work: "Work"}
		repo.On("GetByID", mock.Anything, int32(7)).Return(existing, nil)
		repo.On("Update", mock.Anything, int32(7), mock.MatchedBy(func(p *domain.Person) bool {
			return p.Name == "Alice" && p.Age == 45 && p.Address == "Addr" && p.Work == "Work"
		})).Return(updated, nil)

		body := `{"age":45}`
		req := httptest.NewRequest(http.MethodPatch, "/api/v1/persons/7/", bytes.NewBufferString(body))
		rec := httptest.NewRecorder()

		newTestRouter(repo).ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		var got persondto.PersonResponse
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
		assert.Equal(t, "Alice", got.Name)
		assert.Equal(t, int32(45), got.Age)
		repo.AssertExpectations(t)
	})

	t.Run("not found returns 404", func(t *testing.T) {
		repo := new(mockPersonRepository)
		repo.On("GetByID", mock.Anything, int32(99)).Return(nil, domain.ErrPersonNotFound)

		body := `{"name":"Bob","age":40}`
		req := httptest.NewRequest(http.MethodPatch, "/api/v1/persons/99/", bytes.NewBufferString(body))
		rec := httptest.NewRecorder()

		newTestRouter(repo).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		repo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("non-numeric id returns 404", func(t *testing.T) {
		repo := new(mockPersonRepository)

		body := `{"name":"Bob","age":40}`
		req := httptest.NewRequest(http.MethodPatch, "/api/v1/persons/abc/", bytes.NewBufferString(body))
		rec := httptest.NewRecorder()

		newTestRouter(repo).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		repo.AssertNotCalled(t, "GetByID", mock.Anything, mock.Anything)
		repo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("invalid body returns 400", func(t *testing.T) {
		repo := new(mockPersonRepository)

		req := httptest.NewRequest(http.MethodPatch, "/api/v1/persons/7/", bytes.NewBufferString(`{"age":-5}`))
		rec := httptest.NewRecorder()

		newTestRouter(repo).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		repo.AssertNotCalled(t, "GetByID", mock.Anything, mock.Anything)
		repo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("malformed json returns 400", func(t *testing.T) {
		repo := new(mockPersonRepository)

		req := httptest.NewRequest(http.MethodPatch, "/api/v1/persons/7/", bytes.NewBufferString("{not json"))
		rec := httptest.NewRecorder()

		newTestRouter(repo).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestHandler_Delete(t *testing.T) {
	t.Run("deletes person", func(t *testing.T) {
		repo := new(mockPersonRepository)
		repo.On("Delete", mock.Anything, int32(7)).Return(nil)

		req := httptest.NewRequest(http.MethodDelete, "/api/v1/persons/7/", nil)
		rec := httptest.NewRecorder()

		newTestRouter(repo).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNoContent, rec.Code)
		repo.AssertExpectations(t)
	})

	t.Run("not found still returns 204", func(t *testing.T) {
		repo := new(mockPersonRepository)
		repo.On("Delete", mock.Anything, int32(99)).Return(domain.ErrPersonNotFound)

		req := httptest.NewRequest(http.MethodDelete, "/api/v1/persons/99/", nil)
		rec := httptest.NewRecorder()

		newTestRouter(repo).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNoContent, rec.Code)
	})

	t.Run("non-numeric id returns 204", func(t *testing.T) {
		repo := new(mockPersonRepository)

		req := httptest.NewRequest(http.MethodDelete, "/api/v1/persons/abc/", nil)
		rec := httptest.NewRecorder()

		newTestRouter(repo).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNoContent, rec.Code)
		repo.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything)
	})
}
