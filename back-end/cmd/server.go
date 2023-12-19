package cmd

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"text/template"

	"SocialNetwork/websoct"

	"github.com/gorilla/websocket"
)

func Server(port string) {
	http.HandleFunc("/home", displayFeed)
	http.HandleFunc("/groups/", getGroup)
	/* 	http.HandleFunc("/profile", requireAuth(profile))
	   	http.HandleFunc("/api/profile/{id}", func(w http.ResponseWriter, r *http.Request) {
	   		profile(w, r)
	   	}) */
	http.HandleFunc("/api/profile", func(w http.ResponseWriter, r *http.Request) {
		profile(w, r)
	})
	http.HandleFunc("/api/images", func(w http.ResponseWriter, r *http.Request) {
		println("api/images")
		imageHandler(w, r)
	})
	http.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {
		login(w, r)
	})
	http.HandleFunc("/api/register", func(w http.ResponseWriter, r *http.Request) {
		register(w, r)
	})
	http.HandleFunc("/api/logout", func(w http.ResponseWriter, r *http.Request) {
		logout(w, r)
	})
	http.HandleFunc("/api/submitContent", func(w http.ResponseWriter, r *http.Request) {
		submitContent(w, r)
	})
	http.HandleFunc("/api/messageRequest", messageRequest)
	http.HandleFunc("/api/initialLoad", initialLoad)
	http.HandleFunc("/testing", testFunction)
	http.HandleFunc("/api/validateSession", func(w http.ResponseWriter, r *http.Request) {
		validateSessionHandler(w, r)
	})

	hub := websoct.NewHub()
	websoct.NotificationHub = hub
	go hub.Run()
	http.HandleFunc("/ws/", func(responseWriter http.ResponseWriter, request *http.Request) {
		var upgrader = websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}

		username := strings.Split(request.URL.String(), "/ws/")[1]
		fmt.Println("USER JOINED WS: ", strings.Split(request.URL.String(), "/ws/")[1])

		connection, err := upgrader.Upgrade(responseWriter, request, nil)
		if err != nil {
			log.Println(err)
			return
		}

		websoct.CreateNewSocketUser(hub, connection, username)

	})

	indexTemp, err := template.ParseFiles("../front-end/public/index.html")
	if err != nil {
		log.Fatalln(err)
	}

	fileServer := http.FileServer(http.Dir("../front-end/src"))

	static := http.FileServer(http.Dir("../front-end/static"))
	http.Handle("/static/", http.StripPrefix("/static/", static))

	http.Handle("/src/", http.StripPrefix("/src/", fileServer))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := indexTemp.Execute(w, "")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
	})
	mux := http.DefaultServeMux
	fmt.Printf("Starting server at: http://localhost%s\nPress CTRL + c to shut it down.\n", port)
	err = http.ListenAndServe(port, addCorsHeaders(mux))
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}

func addCorsHeaders(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8081")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		h.ServeHTTP(w, r)
	})
}
