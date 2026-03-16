package routers

import (
	"net/http"

	"github.com/Cinnamoon-dev/blue-gopher/internal/http/handlers"
	"github.com/Cinnamoon-dev/blue-gopher/internal/middleware"
)

type MailRouter struct {
	MailHandler handlers.MailHandler
}

func NewMailRouter(mailHandler handlers.MailHandler) MailRouter {
	return MailRouter{
		MailHandler: mailHandler,
	}
}

func (ro *MailRouter) BaseRoutes() http.HandlerFunc {
	router := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			ro.MailHandler.VerifyUser(w, r)
		case http.MethodPost:
			ro.MailHandler.SendEmail(w, r)
		}
	})

	LoggingRouter := middleware.Logging(router)
	return LoggingRouter
}
