package main

import (
	"encoding/json"
	"io"
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
	mux.HandleFunc("/api/stock", func(w http.ResponseWriter, r *http.Request) {
		var d struct{ K, V string }
		if r.Method == http.MethodPost && json.NewDecoder(io.LimitReader(r.Body, 2<<20)).Decode(&d) == nil {
			ecrireStock(d.K, d.V)
		}
		w.WriteHeader(http.StatusNoContent)
	})
	mux.HandleFunc("/api/pdf", func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(io.LimitReader(r.Body, 40<<20))
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(lirePDFBase64(string(b)))
	})
	mux.HandleFunc("/api/document", func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(io.LimitReader(r.Body, 4<<20))
		var d Document
		if json.Unmarshal(b, &d) != nil {
			http.Error(w, "document illisible", 400)
			return
		}
		pdf, err := genererPDF(d)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		w.Header().Set("Content-Type", "application/pdf")
		w.Write(pdf)
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
