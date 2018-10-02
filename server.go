package main

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"math/rand"
	"net/http"
)

type Link struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

type LinkRequest struct {
	*Link
	ProtectedID string `json:"id"` // override 'id' json to have more control
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}

type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func NewLinkResponse(link *Link) *LinkResponse {
	resp := &LinkResponse{Link: link}
	return resp
}

func (rd *LinkResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	rd.Elapsed = 10
	return nil
}

type LinkResponse struct {
	*Link

	// We add an additional field to the response here.. such as this
	// elapsed computed property
	Elapsed int64 `json:"elapsed"`
}

func (l *LinkRequest) Bind(r *http.Request) error {
	// l.Link is nil if no Link fields are sent in the request. Return an
	// error to avoid a nil pointer dereference.
	if l.Link == nil {
		return errors.New("missing required Link fields.")
	}

	// just a post-process after a decode..
	l.ProtectedID = "" // unset the protected ID
	return nil
}

var links = []*Link{}

func dbNewLink(link *Link) (*Link, error) {
	link.ID = fmt.Sprintf("%d", rand.Intn(100)+10)
	links = append(links, link)
	return link, nil
}

func CreateServer() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(400), 400)
		return
	})
	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(400), 400)
		return
	})
	r.Post("/{id}", func(w http.ResponseWriter, r *http.Request) {
		data := &LinkRequest{}

		if err := render.Bind(r, data); err != nil {
			render.Render(w, r, ErrInvalidRequest(err))
			return
		}

		dbNewLink(data.Link)

		render.Status(r, http.StatusCreated)
		render.Render(w, r, NewLinkResponse(data.Link))
	})

	r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "exists" {
			http.Redirect(w, r, "https://example.com", 302)
			return
		}
		http.Error(w, http.StatusText(404), 404)
	})
	return r
}

func main() {
	r := CreateServer()
	http.ListenAndServe(":3000", r)
}
