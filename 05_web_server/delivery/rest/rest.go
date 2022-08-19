package rest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dr-leee/Go/05_web_server/repository"
	"github.com/dr-leee/Go/05_web_server/types"
	"github.com/golang/gddo/httputil/header"
	"io"
	"log"
	"net/http"
	"strings"
)

//const layoutDT = "2006-01-02 15:04"

type Env struct {
	Repo repository.DbModel
}

func (env *Env) GetJsonLogin(w http.ResponseWriter, r *http.Request) {
	// If the Content-Type header is present, check that it has the value
	// application/json. Note that we are using the gddo/httputil/header
	// package to parse and extract the value here, so the check works
	// even if the client includes additional charset or boundary
	// information in the header.
	if r.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
		if value != "application/json" {
			msg := "Хэдер Content-Type не application/json"
			http.Error(w, msg, http.StatusUnsupportedMediaType)
			return
		}
	}

	// Use http.MaxBytesReader to enforce a maximum read of 1MB from the
	// response body. A request body larger than that will now result in
	// Decode() returning a "http: request body too large" error.
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	// Setup the decoder and call the DisallowUnknownFields() method on it.
	// This will cause Decode() to return a "json: unknown field ..." error
	// if it encounters any extra unexpected fields in the JSON. Strictly
	// speaking, it returns an error for "keys which do not match any
	// non-ignored, exported fields in the destination".
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	var p types.Log
	err := dec.Decode(&p)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		// Catch any syntax errors in the JSON and send an error message
		// which interpolates the location of the problem to make it
		// easier for the client to fix.
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Тело запроса содержит неправильный JSON (на позиции %d)", syntaxError.Offset)
			http.Error(w, msg, http.StatusBadRequest)

		// In some circumstances Decode() may also return an
		// io.ErrUnexpectedEOF error for syntax errors in the JSON. There
		// is an open issue regarding this at
		// https://github.com/golang/go/issues/25956.
		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := fmt.Sprintf("Тело запроса содержит неправильный JSON")
			http.Error(w, msg, http.StatusBadRequest)

		// Catch any type errors, like trying to assign a string in the
		// JSON request body to a int field in our Person struct. We can
		// interpolate the relevant field name and position into the error
		// message to make it easier for the client to fix.
		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Неправильное значение поля %q запроса (на позиции %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			http.Error(w, msg, http.StatusBadRequest)

		// Catch the error caused by extra unexpected fields in the request
		// body. We extract the field name from the error message and
		// interpolate it in our custom error message. There is an open
		// issue at https://github.com/golang/go/issues/29035 regarding
		// turning this into a sentinel error.
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Запрос содержит неизвестное поле %s", fieldName)
			http.Error(w, msg, http.StatusBadRequest)

		// An io.EOF error is returned by Decode() if the request body is
		// empty.
		case errors.Is(err, io.EOF):
			msg := "Запрос не может быть пустым"
			http.Error(w, msg, http.StatusBadRequest)

		// Catch the error caused by the request body being too large. Again
		// there is an open issue regarding turning this into a sentinel
		// error at https://github.com/golang/go/issues/30715.
		case err.Error() == "http: размер запроса превышает лимит":
			msg := "Тело запроса не должно быть больше 1MB"
			http.Error(w, msg, http.StatusRequestEntityTooLarge)

		// Otherwise default to logging the error and sending a 500 Internal
		// Server Error response.
		default:
			log.Print(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	// Call decode again, using a pointer to an empty anonymous struct as
	// the destination. If the request body only contained a single JSON
	// object this will return an io.EOF error. So if we get anything else,
	// we know that there is additional data in the request body.
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		msg := "Запрос должен содержать только один объект JSON"
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	//p.Time = string(time.Now().Format(layoutDT))
	//p.Time = time.Now()
	//fmt.Fprintf(w, "Логин: %+v\n", p)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = env.Repo.CreateLog(ctx, p.Login)
	if err != nil {
		msg := "Ошибка БД"
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
}

func (env *Env) ShowLogins(w http.ResponseWriter, _ *http.Request) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rows, err := env.Repo.GetLogs(ctx)
	if err != nil {
		msg := "Ошибка БД"
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	allLogins := make([]types.Log, 0)

	for rows.Next() {
		login := types.Log{}
		err := rows.Scan(&login.ID, &login.Login, &login.Time)
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(500), 500)
			return
		}
		allLogins = append(allLogins, login)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	var tmpl = `<tr><td>%d</td><td>%s</td><td>%s</td></tr>`
	fmt.Fprintf(w, "<table cellspacing=\"2\" border=\"1\" cellpadding=\"5\">")
	fmt.Fprintf(w, `<tr bgcolor="#cecece"><td>ID</td><td>Login</td><td>Time</td></tr>`)

	for _, login := range allLogins {
		//fmt.Fprintf(w, "%d %s %s\n", login.ID, login.Login, login.Time)
		fmt.Fprintf(w, tmpl, login.ID, login.Login, login.Time)
	}
	fmt.Fprintf(w, "</table>")
}
