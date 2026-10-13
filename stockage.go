package main

// Stockage des réponses de l'utilisateur dans un fichier local :
// %APPDATA%\AQuoiAiJeLeDroit\donnees.json (clé → valeur texte).
// La page est chargée sans adresse web fixe : le stockage du navigateur intégré
// n'y est pas fiable, d'où ce fichier géré par le logiciel.

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

var stockMu sync.Mutex

func cheminStock() string {
	if p := os.Getenv("AQDROIT_STOCK"); p != "" {
		return p
	}
	dir, err := os.UserConfigDir()
	if err != nil {
		return ""
	}
	return filepath.Join(dir, "AQuoiAiJeLeDroit", "donnees.json")
}

func lireStock() map[string]string {
	m := map[string]string{}
	if c := cheminStock(); c != "" {
		if b, err := os.ReadFile(c); err == nil {
			json.Unmarshal(b, &m)
		}
	}
	return m
}

// ecrireStock enregistre une valeur ; une valeur vide supprime la clé.
func ecrireStock(cle, valeur string) {
	if cle == "" || len(cle) > 200 || len(valeur) > 1<<20 {
		return
	}
	stockMu.Lock()
	defer stockMu.Unlock()
	c := cheminStock()
	if c == "" {
		return
	}
	m := lireStock()
	if valeur == "" {
		delete(m, cle)
	} else {
		m[cle] = valeur
	}
	b, _ := json.Marshal(m)
	os.MkdirAll(filepath.Dir(c), 0o755)
	tmp := c + ".tmp"
	if os.WriteFile(tmp, b, 0o600) == nil {
		os.Rename(tmp, c)
	}
}

func stockJSON() string {
	stockMu.Lock()
	defer stockMu.Unlock()
	b, _ := json.Marshal(lireStock())
	return string(b)
}
