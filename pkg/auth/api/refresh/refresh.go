package refresh

import (
	"encoding/json"
	"net/http"
	"survey-api/pkg/auth"
	"survey-api/pkg/di"
	"survey-api/pkg/logger"
	"survey-api/pkg/user"
)

type deps struct {
	loggerService *logger.LoggerService
	cookieService *auth.CookieService
	tokenService  *auth.TokenService
	authRepo      *auth.AuthRepo
	userRepo      *user.UserRepo
	authHandler   *auth.AuthHandler
}

var handler func(http.ResponseWriter, *http.Request)

func Handler() func(http.ResponseWriter, *http.Request) {
	return handler
}

func Init(deps *deps) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		cookie, err := deps.cookieService.ParseSessionCookie(r)
		if err != nil {
			deps.loggerService.LogErr(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		sessionId, err := deps.cookieService.ValidateSessionCookie(cookie)
		if err != nil {
			deps.loggerService.LogErr(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		session, err := deps.authRepo.FindById(sessionId)
		if err != nil {
			deps.loggerService.LogErr(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var generatedCookie *http.Cookie
		token := session.Token
		_, err = deps.tokenService.ValidateJwtToken(token)
		if err != nil {
			newCookie, newToken, err := deps.authHandler.RefreshAuth(session)
			if err != nil {
				deps.loggerService.LogErr(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			generatedCookie = &newCookie
			token = newToken
		}

		user, err := deps.userRepo.FindById(session.UserId.Hex())
		if err != nil {
			deps.loggerService.LogErr(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		authUser := auth.AuthUser{
			Token: token,
			User:  user.ToClientUser(),
		}
		result, err := json.Marshal(authUser)
		if err != nil {
			deps.loggerService.LogErr(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if generatedCookie != nil {
			http.SetCookie(w, generatedCookie)
		}

		w.WriteHeader(http.StatusOK)
		w.Write(result)
	}
}

func init() {
	handler = Init(
		&deps{
			loggerService: di.Container().LoggerService(),
			cookieService: di.Container().CookieService(),
			tokenService:  di.Container().TokenService(),
			authRepo:      di.Container().AuthRepo(),
			userRepo:      di.Container().UserRepo(),
			authHandler:   di.Container().AuthHandler(),
		},
	)
}
