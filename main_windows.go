//go:build windows

package main

import (
	"os"
	"os/exec"
	"path/filepath"

	webview2 "github.com/jchv/go-webview2"
)

func ouvrirLien(url string) {
	if lienAutorise(url) {
		exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	}
}

func main() {
	dir, _ := os.UserConfigDir()
	w := webview2.NewWithOptions(webview2.WebViewOptions{
		AutoFocus: true,
		DataPath:  filepath.Join(dir, "AQuoiAiJeLeDroit", "WebView2"),
		WindowOptions: webview2.WindowOptions{
			Title:  "À quoi ai-je le droit ?",
			Width:  1180,
			Height: 820,
			IconId: 1,
			Center: true,
		},
	})
	if w == nil { // WebView2 absent : ouverture dans le navigateur
		modeNavigateur("127.0.0.1:0", true)
		return
	}
	defer w.Destroy()
	w.Bind("ouvrirLien", func(url string) { ouvrirLien(url) })
	w.Bind("stockEcrire", func(cle, valeur string) { ecrireStock(cle, valeur) })
	// La vérification se fait en arrière-plan pour ne pas figer la fenêtre.
	w.Bind("lancerVerifMaj", func() {
		go func() {
			r := verifierMajJSON()
			w.Dispatch(func() { w.Eval("window.__majResultat && window.__majResultat(" + jsChaine(r) + ")") })
		}()
	})
	w.SetSize(900, 640, webview2.HintMin)
	w.SetHtml(construirePage())
	w.Run()
}
