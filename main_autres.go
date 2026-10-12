//go:build !windows

package main

// Version de développement (Linux/macOS) : interface servie dans le navigateur.
// Usage : go run . [adresse]   (par défaut 127.0.0.1:8765)

import (
	"fmt"
	"os"
	"os/exec"
)

func ouvrirLien(url string) {
	if lienAutorise(url) {
		exec.Command("xdg-open", url).Start()
	}
}

func main() {
	adr := "127.0.0.1:8765"
	if len(os.Args) > 1 {
		adr = os.Args[1]
	}
	fmt.Println("Interface disponible sur http://" + adr + "/")
	if err := modeNavigateur(adr, false); err != nil {
		fmt.Println(err)
	}
}
