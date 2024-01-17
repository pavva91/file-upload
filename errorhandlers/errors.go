package errorhandlers

import "net/http"

func InternalServerErrorHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusInternalServerError)
    w.Write([]byte("500 Internal Server Error"))
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusNotFound)
    w.Write([]byte("404 Not Found"))
}

func BadRequestHandler(w http.ResponseWriter, r *http.Request, err error) {
    w.WriteHeader(http.StatusBadRequest)
    w.Write([]byte(err.Error()))
}
