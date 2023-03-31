// Filename: BIOAFF/backend/cmd/api/middleware.go
package main

import (
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

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
			//lock()
			mu.Lock()
			//check if the IP address is in the map
			if _, found := clients[ip]; !found {
				clients[ip] = &client{limiter: rate.NewLimiter(
					rate.Limit(app.config.limiter.rps), app.config.limiter.burst)}
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
/*func (app *application) authenicate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//Add a "Vary: Authorization header to the reponse"
		//A note to caches that no reponse may vary
		w.Header().Add("Vary", "Authorization")
		//Retrieve the value of the Authorization header form the request
		authorizationHeader := r.Header.Get("Authorization")
		//if no authorization found then we will created an anonymous user
		if authorizationHeader == "" {
			r = app.contextSetUser(r, data.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		//Check if the provided authorization header is in the right format
		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			app.invalidAuthenticationTokenReponse(w, r)
			return
		}
		//Extract the token
		token := headerParts[1]
		//Validate the token
		v := validator.New()
		if data.ValidateTokenPlaintext(v, token); !v.Valid() {
			app.invalidCredntialsResponse(w, r)
			return
		}

		//Retrieve details about the user
		user, err := app.models.Users.GetForToken(data.ScopeAuthentication, token)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				app.invalidCredntialsResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
			}
			return
		}
		//Add the user infromation to the request context
		r = app.contextSetUser(r, user)

		//Call the next handler
		next.ServeHTTP(w, r)
	})
}

// Check for activated user
func (app *application) requireAuthenitcatedUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//Get the user
		user := app.contextGetUser(r)

		//Check for anonymous user
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

		//Check if user is activated
		if !user.Activated {
			app.inactiveAccountResponse(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})

	return app.requireAuthenitcatedUser(fn)
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
		//check for the permissions
		if !permissions.Include(code) {
			app.notPermittedResponse(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
	return app.requireActivatedUser(fn)
}*/

// Enable CORS
func (app *application) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//Add the "Vary:Origin" headers
		w.Header().Add("Vary", "Origin")

		//Get the value of the request's origin headers
		origin := r.Header.Get("Origin")
		//check if Origin header is present
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
