package main

import (
	"context"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"net/http"
)

func (app *application) wrap(next http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		ctx := context.WithValue(r.Context(), "params", ps)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func (app *application) routes() http.Handler {
	router := httprouter.New()
	secure := alice.New(app.checkToken)

	router.HandlerFunc(http.MethodGet, "/status", app.statusHandler)

	router.HandlerFunc(http.MethodPost, "/v1/signin", app.Signin)

	router.GET("/v1/movies", app.wrap(secure.ThenFunc(app.getAllMovies)))
	//router.HandlerFunc(http.MethodGet, "/v1/movies", app.getAllMovies)
	router.HandlerFunc(http.MethodGet, "/v1/movies/genres/:id", app.getAllMoviesByGenre)
	router.HandlerFunc(http.MethodGet, "/v1/movie/:id", app.getOneMovie)
	router.HandlerFunc(http.MethodGet, "/v1/genres", app.getAllGenres)

	router.POST("/v1/admin/editmovie", app.wrap(secure.ThenFunc(app.editMovie)))
	router.DELETE("/v1/admin/deletemvoie/:id", app.wrap(secure.ThenFunc(app.deleteMovie)))

	//router.HandlerFunc(http.MethodPost, "/v1/admin/editmovie", app.editMovie)
	//router.HandlerFunc(http.MethodDelete, "/v1/admin/deletemvoie/:id", app.deleteMovie)

	return app.enableCORS(router)
}
