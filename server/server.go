package server

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"

	"service-boilerplate/authentication"
	"service-boilerplate/constants"
)

type Server struct {
	router  chi.Router
	authSvc authentication.Service
}

func New(
	authSvc authentication.Service,
) *Server {
	s := &Server{
		authSvc: authSvc,
	}
	r := chi.NewRouter()

	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
	r.Use(cors.Handler)

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.StripSlashes)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Group(func(r chi.Router) {

		r.Route("/v1", func(r chi.Router) {
			// implements

		})
	})

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		render.Respond(w, r, SuccessResponse(nil, "OK"))
	})

	s.router = r

	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		route = strings.Replace(route, "/*/", "/", -1)
		log.Printf("%s %s\n", method, route) // Walk and print out all routes
		return nil
	}

	if err := chi.Walk(r, walkFunc); err != nil {
		log.Panicf("Logging err: %s\n", err.Error()) // panic if there is an error
	}
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status,omitempty"`  // user-level status message
	AppCode    int64  `json:"code,omitempty"`    // application-specific error code
	Message    string `json:"message,omitempty"` // application-level error message, for debugging
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusBadRequest,
		AppCode:        constants.CodeError,
		Message:        err.Error(),
	}
}

type ApiResponse struct {
	HTTPStatusCode int `json:"-"` // http response status code

	AppCode int64       `json:"code,omitempty"` // application-specific error code
	Data    interface{} `json:"data,omitempty"` // application-specific error code
	Message string      `json:"message"`        // application-level error message, for debugging
}

func (e *ApiResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}
func SuccessResponse(data interface{}, msg string) render.Renderer {
	return &ApiResponse{
		HTTPStatusCode: http.StatusOK,
		AppCode:        constants.CodeSuccess,
		Data:           data,
		Message:        msg,
	}
}
