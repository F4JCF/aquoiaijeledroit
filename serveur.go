package main

import (
	"net"
	"net/http"
)

// modeNavigateur : solution de secours si WebView2 est absent.
// L'interface est servie en local (127.0.0.1) et ouverte dans le navigateur par défaut.
func modeNavigateur(adresse string, ouvrir bool) error {
	ln, err := net.Listen("tcp", adresse)
	if err != nil {
		return err
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(construirePage()))
	})
	mux.HandleFunc("/api/verifierMaj", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(verifierMajJSON()))
	})
	if ouvrir {
		ouvrirLien("http://" + ln.Addr().String() + "/")
	}
	return http.Serve(ln, mux)
}
