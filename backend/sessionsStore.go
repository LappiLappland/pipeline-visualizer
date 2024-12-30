package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/imacks/bitflags-go"
	"github.com/markbates/goth/gothic"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var sessionKey []byte

var store *sessions.CookieStore

func initSessionsStore() {
	sessionKey = []byte(os.Getenv("SESSION_KEY"))
	store = sessions.NewCookieStore([]byte(sessionKey))
	store.Options = &sessions.Options{
		Path:     "/",
		Domain:   os.Getenv("HOST_DOMAIN"),
		MaxAge:   86400 * 7,
		HttpOnly: false,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	}
}

type LoginData struct {
	Token string `json:"token"`
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("name")
		password := r.FormValue("password")

		user, err := AuthenticateUser(username, password)
		if err != nil {
			http.Error(w, "Invalid login credentials", http.StatusUnauthorized)
			return
		}

		session, _ := store.Get(r, "session-name")
		session.Values["id"] = user.ID
		session.Values["name"] = user.Name
		session.Values["provider"] = "host"
		println(user.ID, user.Name, session.Values["name"].(string), session.Name())

		err = session.Save(r, w)
		if err != nil {
			log.Println(err)
			return
		}

		http.Redirect(w, r, HostAddress+"/pipeline/run/1", http.StatusSeeOther)
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	LogoutUser(w, r)
	gothic.Logout(w, r)

	http.Redirect(w, r, HostAddress+"/login", http.StatusSeeOther)
}

func LogoutUser(w http.ResponseWriter, r *http.Request) error {
	session, _ := store.Get(r, "session-name")
	session.Values["id"] = nil
	session.Values["name"] = nil
	session.Values["provider"] = nil
	session.Options.MaxAge = -1 // Mark session as expired
	return session.Save(r, w)
}

func AuthenticateUser(name string, password string) (*User, error) {
	var user User
	err := db.Where("name = ?", name).First(&user).Error
	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &user, nil
}

func requirePipelineAuth(permission UserPermissions) Adapter {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Try to get session
			session, _ := store.Get(r, "session-name")
			userID, _ := session.Values["id"].(uint)

			vars := mux.Vars(r)
			pipelineId := vars["id"]

			var pipelineWorker PipelineWorker
			db.Preload("Pipeline").Preload("Pipeline.Users", func(db *gorm.DB) *gorm.DB {
				return db.Where("user_id = ?", userID)
			}).Preload("Pipeline.Users.Role").First(&pipelineWorker, pipelineId)

			// Pipeline does not exist
			pipeline := pipelineWorker.Pipeline
			if pipeline.ID == 0 {
				w.WriteHeader(404)
				return
			}

			// Serve pipeline if it's public and requires only read permission
			if !pipeline.IsPrivate && !bitflags.HasAny(permission, write, execute, admin) {
				next.ServeHTTP(w, r)
				return
			}

			// Deny if pipeline is private and user has no roles
			if len(pipeline.Users) == 0 && pipeline.IsPrivate {
				w.WriteHeader(403)
				return
			}

			userPermissions := pipeline.Users[0].Role.Permissions
			// Check read
			if pipeline.IsPrivate {
				if bitflags.Has(permission, read) && !bitflags.Has(userPermissions, read) {
					w.WriteHeader(403)
					return
				}
			}

			// Check write
			if bitflags.Has(permission, write) && !bitflags.Has(userPermissions, write) {
				w.WriteHeader(403)
				return
			}

			// Check execute
			if bitflags.Has(permission, execute) && !bitflags.Has(userPermissions, execute) {
				w.WriteHeader(403)
				return
			}

			// Check admin
			if bitflags.Has(permission, admin) && !bitflags.Has(userPermissions, admin) {
				w.WriteHeader(403)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
