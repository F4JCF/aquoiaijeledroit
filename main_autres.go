//go:build !windows

package main

// Version de développement (Linux/macOS) : interface servie dans le navigateur.
// Usage : go run . [adresse]   (par défaut 127.0.0.1:8765)
//         go run . -site <dossier>   (génère la version site web)

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func ouvrirLien(url string) {
	if lienAutorise(url) {
		exec.Command("xdg-open", url).Start()
	}
}

func main() {
	if len(os.Args) > 2 && os.Args[1] == "-site" {
		dir := os.Args[2]
		os.MkdirAll(dir, 0o755)
		if err := os.WriteFile(filepath.Join(dir, "index.html"), []byte(construirePageWeb()), 0o644); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		os.WriteFile(filepath.Join(dir, "icone.png"), icone, 0o644)
		fmt.Println("Site généré dans", dir)
		return
	}
	adr := "127.0.0.1:8765"
	if len(os.Args) > 1 {
		adr = os.Args[1]
	}
	fmt.Println("Interface disponible sur http://" + adr + "/")
	if err := modeNavigateur(adr, false); err != nil {
		fmt.Println(err)
	}
}
