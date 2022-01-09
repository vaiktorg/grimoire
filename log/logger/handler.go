package logger

import (
	"crypto/subtle"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Connection struct {
	user      string
	password  string
	connected bool
	appkey    string
}

func (l *Logger) GinRoutes(e *gin.Engine) {
	e.GET("/debug/msgs", func(c *gin.Context) {
		c.JSON(200, l.msgstore)
	})

	e.POST("/debug/msg", func(c *gin.Context) {
		msg := new(Log)

		err := c.Bind(msg)
		if err != nil {
			http.Error(c.Writer, err.Error(), 500)
			return
		}

		l.INBOUND(*msg)
	})
}

func (l *Logger) HttpRoutes(s *http.ServeMux) {
	s.HandleFunc("/debug/connect", l.BasicAuth(func(w http.ResponseWriter, r *http.Request) {
		if u, p, ok := r.BasicAuth(); ok {
			if u != l.connection.user || p != l.connection.password {
				http.Error(w, "wrong appname or appcode", 403)
				return
			}

			l.connection.connected = true
			_, _ = w.Write([]byte(l.connection.appkey))
		}
	}, l.connection.appkey))

	s.HandleFunc("/debug/msgs", l.BasicAuth(func(w http.ResponseWriter, r *http.Request) {
		//auth, err := r.Cookie("Authentication")
		//if err != nil {
		//	http.Error(w, err.Error(), 500)
		//	return
		//}
		//fmt.Println(auth)
		if r.Method != http.MethodGet {
			http.Error(w, "only GET method allowed", http.StatusMethodNotAllowed)
			return
		}

		if !l.connection.connected {
			http.Error(w, "currently not authenticated", http.StatusUnauthorized)
			return
		}

		err := json.NewEncoder(w).Encode(l.msgstore)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

	}, l.connection.appkey))

	s.HandleFunc("/debug/msg", l.BasicAuth(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "only POST method allowed", http.StatusMethodNotAllowed)
			return
		}
		msg := new(Log)

		err := json.NewDecoder(r.Body).Decode(msg)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		l.INBOUND(*msg)
	}, l.connection.appkey))
}
func (l *Logger) SetAppName(name string) {
	l.connection.user = name
}

func (l *Logger) BasicAuth(handler http.HandlerFunc, realm string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok || subtle.ConstantTimeCompare([]byte(user),
			[]byte(l.connection.user)) != 1 || subtle.ConstantTimeCompare([]byte(pass), []byte(l.connection.password)) != 1 {
			w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
			w.WriteHeader(401)
			_, _ = w.Write([]byte("Unauthorized.\n"))
			return
		}
		handler(w, r)
	}
}
