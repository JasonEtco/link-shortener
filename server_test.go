package main

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
)

var testClient = &http.Client{
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	},
}

func GetDatabaseURL() string {
	return os.Getenv("ES_URL")
}

func InsertLinkIntoDB(link *Link) error {
	db, err := NewDB(GetDatabaseURL())
	if err != nil {
		return err
	}
	link, err = db.AddLink(link)
	if err != nil {
		return err
	}
	return nil
}

func TestLinkGetNotFound(t *testing.T) {
	require := require.New(t)
	r, err := CreateServer(GetDatabaseURL())
	require.NoError(err)
	server := httptest.NewServer(r)
	defer server.Close()

	resp, err := testClient.Get(server.URL + "/doesntexist")
	require.NoError(err)
	require.Equal(404, resp.StatusCode)
}

func TestLinkGetFound(t *testing.T) {
	require := require.New(t)

	link := Link{ID: "abc", URL: "https://example.com"}
	err := InsertLinkIntoDB(&link)
	require.NoError(err)

	r, err := CreateServer(GetDatabaseURL())
	server := httptest.NewServer(r)
	defer server.Close()

	resp, err := testClient.Get(server.URL + "/abc")
	require.NoError(err)

	require.Equal(302, resp.StatusCode)
	require.Equal("https://example.com", resp.Header.Get("Location"))

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	require.NoError(err)
	body := string(bodyBytes[:])

	require.Equal("<a href=\"https://example.com\">Found</a>.\n\n", body)
}

// Post scenarios
// 	1. URL doesn't exits in DB and ID not provided.
// 	2. ID & URL doesn't exist in DB.
// 	3. ID doesn't exits in DB but URL does.
// 	4. URL doesn't exits in DB but ID does.
// 	5. URL does exits in DB and ID not provided.
// 	6. URL and ID not provided.
// 	7. URL not provided.

func TestLinkPostJSON(t *testing.T) {
	require := require.New(t)
	r, err := CreateServer(GetDatabaseURL())
	require.NoError(err)
	server := httptest.NewServer(r)
	defer server.Close()
	json := []byte(`{"url": "https://example.com"}`)

	resp, err := http.Post(server.URL+"/new-link", "application/json", bytes.NewBuffer(json))
	require.NoError(err)
	require.Equal(201, resp.StatusCode)

	// TODO: Test headers and response body
}

func TestLinkPostFormData(t *testing.T) {
	require := require.New(t)
	r, err := CreateServer(GetDatabaseURL())
	require.NoError(err)
	server := httptest.NewServer(r)
	defer server.Close()

	json := []byte(`url=https%3A%2F%2Fexample.com`)
	resp, err := http.Post(server.URL+"/new-link", "application/x-www-form-urlencoded", bytes.NewBuffer(json))
	require.NoError(err)
	require.Equal(201, resp.StatusCode)

	// TODO: Test headers and response body
}

func TestLinkPostJSONURLNotProvided(t *testing.T) {
	require := require.New(t)
	r, err := CreateServer(GetDatabaseURL())
	require.NoError(err)
	server := httptest.NewServer(r)
	defer server.Close()

	json := []byte(`{}`)
	resp, err := http.Post(server.URL+"/new-link", "application/json", bytes.NewBuffer(json))
	require.NoError(err)
	require.Equal(400, resp.StatusCode)
}

func TestLinkPostFormDataURLNotProvided(t *testing.T) {
	require := require.New(t)
	r, err := CreateServer(GetDatabaseURL())
	require.NoError(err)
	server := httptest.NewServer(r)
	defer server.Close()

	resp, err := http.PostForm(server.URL+"/new-link", url.Values{})
	require.NoError(err)
	require.Equal(400, resp.StatusCode)
}

func TestLinkPostJSONURLAndIDNotProvided(t *testing.T) {
	require := require.New(t)
	r, err := CreateServer(GetDatabaseURL())
	require.NoError(err)
	server := httptest.NewServer(r)
	defer server.Close()
	json := []byte(`{}`)

	resp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(json))
	require.NoError(err)
	require.Equal(400, resp.StatusCode)
}

func TestLinkPostFormDataURLAndIDNotProvided(t *testing.T) {
	require := require.New(t)
	r, err := CreateServer(GetDatabaseURL())
	require.NoError(err)
	server := httptest.NewServer(r)
	defer server.Close()

	resp, err := http.PostForm(server.URL, url.Values{})
	require.NoError(err)
	require.Equal(400, resp.StatusCode)
}
