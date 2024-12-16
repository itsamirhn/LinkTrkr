package internal

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/itsamirhn/linktrkr/internal/config"
	"github.com/itsamirhn/linktrkr/internal/models"
	"github.com/itsamirhn/linktrkr/pkg"
)

const PathPrefix = "/r"

type RedirectHandler struct {
	notifier   pkg.Notifier[models.Click]
	jwtService pkg.JWTService[models.Redirect]
}

func NewRedirectHandler(
	jwtService pkg.JWTService[models.Redirect],
	notifier pkg.Notifier[models.Click],
) *RedirectHandler {
	return &RedirectHandler{notifier: notifier, jwtService: jwtService}
}

func (h *RedirectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug, ok := vars["slug"]
	if !ok {
		http.Error(w, "missing slug", http.StatusBadRequest)
		return
	}

	message, err := h.jwtService.Decode(slug)
	if err != nil {
		http.Error(w, "invalid slug", http.StatusBadRequest)
		return
	}

	click := &models.Click{
		Redirect: message,
		Request:  r,
	}

	go h.notifier.Notify(*click)

	http.Redirect(w, r, message.URL, http.StatusFound)
}

func (h *RedirectHandler) Path() string {
	return PathPrefix + "/{slug}"
}

func GetRedirectURL(slug string) string {
	return fmt.Sprintf("https://"+config.GlobalConfig.Server.Endpoint+PathPrefix+"/%s", slug) // FIXME
}
