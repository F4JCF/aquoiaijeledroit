//go:build windows

package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	webview2 "github.com/jchv/go-webview2"
)

func ouvrirLien(url string) {
	if lienAutorise(url) || strings.HasPrefix(url, "mailto:") {
		exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	}
}

func main() {
	// Mise à jour téléchargée lors d'une session précédente : on l'installe avant d'ouvrir la fenêtre.
	nettoyerMaj()
	if majPrete() && installerMaj() == nil && relancer() == nil {
		return
	}
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
	w.Bind("lirePDF", func(b64 string) map[string]string { return lirePDFBase64(b64) })
	w.Bind("enregistrerPDF", func(doc string) map[string]string { return enregistrerDocument(doc) })
	// Ouvre le PDF, ou le dossier en sélectionnant le fichier (pour le joindre à un e-mail).
	w.Bind("ouvrirDocument", func(chemin string, dossier bool) bool {
		if !cheminDocumentAutorise(chemin) {
			return false
		}
		if dossier {
			return exec.Command("explorer", "/select,", chemin).Start() == nil
		}
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", chemin).Start() == nil
	})
	// La vérification se fait en arrière-plan pour ne pas figer la fenêtre.
	w.Bind("lancerVerifMaj", func() {
		go func() {
			res := verifierMaj()
			// Nouvelle version du logiciel : téléchargement automatique avant d'informer l'interface.
			if res.App != nil && res.App.Asset != nil {
				res.App.Prete = telechargerMaj(*res.App.Asset) == nil
			}
			b, _ := json.Marshal(res)
			w.Dispatch(func() { w.Eval("window.__majResultat && window.__majResultat(" + jsChaine(string(b)) + ")") })
		}()
	})
	// « Redémarrer maintenant » : installe la version téléchargée et relance le logiciel.
	w.Bind("redemarrerMaj", func() bool {
		if installerMaj() != nil {
			return false
		}
		if relancer() != nil {
			return false
		}
		w.Dispatch(func() { w.Terminate() })
		return true
	})
	w.SetSize(900, 640, webview2.HintMin)
	w.SetHtml(construirePage())
	w.Run()
}
