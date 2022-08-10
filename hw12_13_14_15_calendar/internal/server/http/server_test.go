package internalhttp

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sn5ake6/otus-golang/hw12_13_14_15_calendar/internal/app"
	"github.com/sn5ake6/otus-golang/hw12_13_14_15_calendar/internal/logger"
	memorystorage "github.com/sn5ake6/otus-golang/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/stretchr/testify/require"
)

func TestHTTPServerCRUDOperations(t *testing.T) {
	logg, err := logger.New("info")
	require.NoError(t, err)

	storage := memorystorage.New()
	router := NewRouter(logg, app.New(logg, storage))

	body := bytes.NewBufferString(`{
		"id":"512bc5cd-01e9-4639-99a2-d42fe25dec62",
		"title":"some title",
		"beginAt":"2022-07-24T16:00:00Z",
		"endAt":"2022-07-24T18:00:00Z",
		"description":"some description",
		"userId":"6b216e09-7ab3-41f9-ba57-cc94d45fe759",
		"notifyAt":"2022-07-24T15:00:00Z",
		"notifiedAt":"0001-01-01T00:00:00Z"
	}`)

	expectedBody := replaceUnexpectedSpaces(body.String())

	updatedBody := bytes.NewBufferString(`{
		"id":"512bc5cd-01e9-4639-99a2-d42fe25dec62",
		"title":"upated title",
		"beginAt":"2022-07-25T16:00:00Z",
		"endAt":"2022-07-25T18:00:00Z",
		"description":"updated description",
		"userId":"6b216e09-7ab3-41f9-ba57-cc94d45fe759",
		"notifyAt":"2022-07-25T15:00:00Z",
		"notifiedAt":"0001-01-01T00:00:00Z"
	}`)

	expectedUpdatedBody := replaceUnexpectedSpaces(updatedBody.String())

	t.Run("create success case", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/events/", body)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		require.Equal(t, http.StatusCreated, resp.StatusCode)

		body, err := ioutil.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Equal(t, "", string(body))
	})

	t.Run("get success case", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/events/512bc5cd-01e9-4639-99a2-d42fe25dec62/", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)

		body, err := ioutil.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Equal(t, expectedBody, string(body))
	})

	t.Run("get fail case", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/events/6b216e09-7ab3-41f9-ba57-cc94d45fe759/", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		require.NotEqual(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("update success case", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/events/512bc5cd-01e9-4639-99a2-d42fe25dec62/", updatedBody)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)

		body, err := ioutil.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Equal(t, "", string(body))

		reqGet := httptest.NewRequest(http.MethodGet, "/events/512bc5cd-01e9-4639-99a2-d42fe25dec62/", nil)
		wGet := httptest.NewRecorder()

		router.ServeHTTP(wGet, reqGet)

		respGet := wGet.Result()
		defer respGet.Body.Close()
		require.Equal(t, http.StatusOK, respGet.StatusCode)

		bodyGet, err := ioutil.ReadAll(respGet.Body)
		require.NoError(t, err)
		require.Equal(t, expectedUpdatedBody, string(bodyGet))
	})

	t.Run("update fail case", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/events/6b216e09-7ab3-41f9-ba57-cc94d45fe759/", updatedBody)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		require.NotEqual(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("delete success cases", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/events/512bc5cd-01e9-4639-99a2-d42fe25dec62/", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)

		reqGet := httptest.NewRequest(http.MethodGet, "/events/512bc5cd-01e9-4639-99a2-d42fe25dec62/", nil)
		wGet := httptest.NewRecorder()

		router.ServeHTTP(wGet, reqGet)

		respGet := wGet.Result()
		defer respGet.Body.Close()
		require.NotEqual(t, http.StatusOK, respGet.StatusCode)
	})

	t.Run("delete fail cases", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/events/6b216e09-7ab3-41f9-ba57-cc94d45fe759/", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		require.NotEqual(t, http.StatusOK, resp.StatusCode)
	})
}

func TestHTTPServerSelectOperations(t *testing.T) {
	logg, err := logger.New("info")
	require.NoError(t, err)

	storage := memorystorage.New()
	router := NewRouter(logg, app.New(logg, storage))

	body := bytes.NewBufferString(`{
		"id":"512bc5cd-01e9-4639-99a2-d42fe25dec62",
		"title":"some title",
		"beginAt":"2022-07-24T16:00:00Z",
		"endAt":"2022-07-24T18:00:00Z",
		"description":"some description",
		"userId":"6b216e09-7ab3-41f9-ba57-cc94d45fe759",
		"notifyAt":"2022-07-24T15:00:00Z",
		"notifiedAt":"0001-01-01T00:00:00Z"
	}`)

	expectedBody := replaceUnexpectedSpaces(body.String())

	t.Run("create case", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/events/", body)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		require.Equal(t, http.StatusCreated, resp.StatusCode)

		body, err := ioutil.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Equal(t, "", string(body))
	})

	t.Run("select on day exists case", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/events/?period=day&date=2022-07-24%2015:00:00", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)

		body, err := ioutil.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Equal(t, "["+expectedBody+"]", string(body))
	})

	t.Run("select on day empty case", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/events/?period=day&date=2022-07-26%2015:00:00", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)

		body, err := ioutil.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Equal(t, "[]", string(body))
	})

	t.Run("select on week exists case", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/events/?period=week&date=2022-07-24%2015:00:00", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)

		body, err := ioutil.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Equal(t, "["+expectedBody+"]", string(body))
	})

	t.Run("select on week empty case", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/events/?period=week&date=2022-07-26%2015:00:00", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)

		body, err := ioutil.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Equal(t, "[]", string(body))
	})

	t.Run("select on month exists case", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/events/?period=month&date=2022-07-24%2015:00:00", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)

		body, err := ioutil.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Equal(t, "["+expectedBody+"]", string(body))
	})

	t.Run("select on month empty case", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/events/?period=month&date=2022-07-26%2015:00:00", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)

		body, err := ioutil.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Equal(t, "[]", string(body))
	})
}

func replaceUnexpectedSpaces(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(s, "\n", ""), "\t", "")
}
