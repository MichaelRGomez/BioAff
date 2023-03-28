// BIOAFF/backend/cmd/api/middleware.go

package main

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// recoverPanic tries to return to a normal state otherwise,
// it starts the gradual shutdown
func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (app *application) rateLimit(next http.Handler) http.Handler {
	//Create a client type
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	//Launch a background goroutine that removes old entries
	//from the clients map once every minute
	go func() {
		for {
			time.Sleep(time.Minute)
			//Lock before starting to clean
			mu.Lock()
			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.config.limiter.enabled {
			//Get the ip address of the request
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				app.serverErrorResponse(w, r, err)
				return
			}
			//Update the last seen time of the client
			clients[ip].lastSeen = time.Now()

			//check if request allowed
			if !clients[ip].limiter.Allow() {
				mu.Unlock()
				app.rateLimitExceededResponse(w, r)
				return
			}
			mu.Unlock()
		}
		next.ServeHTTP(w, r)
	})
}

// Authentication
func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//Add a "Vary: Authorization header to the response"
		//A note to caches that no response may vary

		w.Header().Add("Vary", "Authorization")

		//Retrieve the value of the Authorization header form the request

		authorizationHeader := r.Header.Get("Authorization")

		//if no authorization found then we will create an anonymous user
		if authorizationHeader == "" {
			r = app.contextSetUser(r, data.anonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		//check if the provided authorization header is in the right format
		headerParts := strings.Spit(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		//Extract the token
		token := headerParts[1]

		//Validate the token
		v := validator.New()
		if data.ValidateTokenPlaintext(v, token); !v.valid() {
			app.invalidCredentialsResponse(w, r)
			return
		}

		//Retrive details about the user
		user, err := app.models.Users.GetForToken(data.ScopeAuthentication, token)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				app.invalidCredentialsResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
			}
			return
		}
		//Add the user infromation to the request context
		r = app.contextSetUser(r, user)

		//call the next handler
		next.ServeHTTP(w, r)
	})
}

// Check for activated user
func (app *application) requireAutheniticatedUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandleFunc(func(w http.ResponseWriter, r *http.Request) {
		//get the user
		user := app.contextGetUser(r)

		//check for anonymous user
		if user.IsAnonymous() {
			app.authenticationRequiredResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Check for activated user
func (app *application) requireActivatedUser(next http.HandlerFunc) http.HandlerFunc {
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//Get the user
		user := app.contextGetUser(r)

		//check if user is activated
		if !user.activated {
			app.inactiveAccountResponse(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
	return app.requireAutheniticatedUser(fn)
}

// Check for activated user
func (app *application) requirePermission(code string, next http.HandlerFunc) http.HandlerFunc {
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//Get the user
		user := app.contextGetUser(r)

		//get the permission slice for the user
		permissions, err := app.models.Permissions.GetAllForUser(user.ID)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		//check for the permission
		if !permission.Include(code) {
			app.notPermittedResponse(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
	return app.requireActivatedUser(fn)
}

// Enable CORS
func (app *application) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//add the "Vary-Origin" headers
		w.Header().add("Vary", "Origin")

		//Get the value of the request's origin header
		origin := r.Header.Get("Origin")
		//check if origin header is present
		if origin != "" {
			for i := range app.config.cors.trustedOrigins {
				if origin == app.config.cors.trustedOrigins[i] {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					break
				}
			}
		}
		next.ServeHTTP(w, r)
	})
}
