package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"net/smtp"
	"strings"

	"github.com/Cinnamoon-dev/blue-gopher/internal/services"
	"github.com/Cinnamoon-dev/blue-gopher/pkg/config"
	"github.com/golang-jwt/jwt/v5"
)

type MailHandler struct {
	AuthService services.AuthService
	UserService services.UserService
	MailService services.MailService
}

func NewMailHandler(authService services.AuthService, userService services.UserService, mailService services.MailService) MailHandler {
	return MailHandler{
		AuthService: authService,
		UserService: userService,
		MailService: mailService,
	}
}

type SendMailRequest struct {
	Email string `json:"email"`
}

func sendEmail(to string, subject string, body string) {
	env := config.NewEnv()
	from := env.MailUsername
	pass := env.MailPassword

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" +
		body

	err := smtp.SendMail("smtp.gmail.com:587", smtp.PlainAuth("", from, pass, "smtp.gmail.com"), from, []string{to}, []byte(msg))
	if err != nil {
		log.Printf("smtp error: %s", err)
		return
	}

	log.Print("mail sent!")
}

func parseToken(path string) (string, error) {
	// path = /mail/{token}
	parts := strings.Split(path, "/")

	if len(parts) < 3 {
		return "", http.ErrNotSupported
	}

	return parts[2], nil
}

func (h *MailHandler) VerifyUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	env := config.NewEnv()
	emailToken, err := parseToken(r.URL.Path)
	if err != nil {
		RespondError(w, err)
		return
	}

	claims, err := h.AuthService.DecodeMailToken(emailToken, jwt.SigningMethodHS256, []byte(env.JwtKey))
	if err != nil {
		RespondError(w, err)
		return
	}

	email := claims.Email
	user, err := h.UserService.GetByEmail(ctx, email)
	if err != nil {
		RespondError(w, err)
		return
	}

	if user.IsVerified == true {
		RespondJSON(w, http.StatusOK, map[string]string{"message": "User already verified"})
		return
	}

	user.IsVerified = true
	err = h.UserService.Update(ctx, user.ID, *user)
	if err != nil {
		RespondError(w, err)
		return
	}

	RespondJSON(w, http.StatusOK, map[string]string{"message": "User verified"})
}

func (h *MailHandler) SendEmail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var body SendMailRequest

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		RespondJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		return
	}

	_, err = h.UserService.GetByEmail(ctx, body.Email)
	if err != nil {
		RespondJSON(w, http.StatusOK, map[string]string{"message": "An email is going to be sent if there is an account with this email"})
		return
	}

	env := config.NewEnv()
	emailToken, _ := h.AuthService.CreateToken(jwt.MapClaims{
		"email": body.Email,
	}, jwt.SigningMethodHS256, []byte(env.JwtKey))
	link := env.BackendUrl + "/mail/" + emailToken

	sendEmail(body.Email, "Email Verification", link)
	RespondJSON(w, http.StatusOK, map[string]string{"message": "An email is going to be sent if there is an account with this email"})
}
