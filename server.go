package main

import (
	"github.com/go-chi/chi"
	"github.com/syntaqx/render"
	"net/http"
	"net/url"
	"os"
)

var db *DB

type Link struct {
	ID  string `json:"id" form:"id"`
	URL string `json:"url" form:"url,omitempty"`
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:        err,
		StatusCode: 400,
	}
}

func ErrNotFound(err error) render.Renderer {
	return &ErrResponse{
		Err:        err,
		StatusCode: 404,
	}
}

type ErrResponse struct {
	Err error `json:"-"` // low-level runtime error

	StatusCode int    `json:"code"`            // user-level status message
	StatusText string `json:"status"`          // user-level status message
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	e.StatusText = http.StatusText(e.StatusCode)
	e.ErrorText = e.Err.Error()
	render.Status(r, e.StatusCode)
	return nil
}

func (link *Link) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (link *Link) Bind(r *http.Request) error {
	_, err := url.Parse(link.URL)
	if err != nil {
		return err
	}
	return nil
}

func CreateServer() *chi.Mux {
	render.Respond = Respond

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
		link := &Link{}

		if err := render.Bind(r, link); err != nil {
			render.Render(w, r, ErrInvalidRequest(err))
			return
		}

		db.NewLink(link)

		render.Status(r, http.StatusCreated)
		render.Render(w, r, link)
	})

	r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
		ID := chi.URLParam(r, "id")
		link, err := db.GetLink(ID)
		if err != nil {
			render.Render(w, r, ErrNotFound(err))
			return
		}
		http.Redirect(w, r, link.URL, 302)
	})
	return r
}

func Respond(w http.ResponseWriter, r *http.Request, v interface{}) {
	// Format response based on request Accept header.
	switch render.GetAcceptedContentType(r) {
	case render.ContentTypeJSON:
		render.JSON(w, r, v)
	default:
		render.XML(w, r, v)
	}
}

func main() {
	var err error
	db, err = NewDB(os.Getenv("ES_URL"))
	if err != nil {
		panic(err)
	}
	r := CreateServer()
	http.ListenAndServe(":3000", r)
}
