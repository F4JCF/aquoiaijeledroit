package main

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestMiseAJourLogiciel(t *testing.T) {
	dir := t.TempDir()
	exe := filepath.Join(dir, "AQuoiAiJeLeDroit.exe")
	os.WriteFile(exe, []byte("MZ ancienne version"), 0o755)
	ancien := cheminExe
	cheminExe = func() string { return exe }
	defer func() { cheminExe = ancien }()

	nouveau := append([]byte("MZ"), make([]byte, 2<<20)...)
	somme := sha256.Sum256(nouveau)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write(nouveau) }))
	defer srv.Close()

	// empreinte fausse : refusée
	if err := telechargerMaj(assetGitHub{URL: srv.URL, Digest: "sha256:00"}); err == nil || majPrete() {
		t.Fatal("un fichier à l'empreinte incorrecte a été accepté")
	}
	// bonne empreinte : acceptée puis installée
	a := assetGitHub{URL: srv.URL, Taille: int64(len(nouveau)), Digest: "sha256:" + hex.EncodeToString(somme[:])}
	if err := telechargerMaj(a); err != nil || !majPrete() {
		t.Fatalf("téléchargement refusé : %v", err)
	}
	if err := installerMaj(); err != nil {
		t.Fatal(err)
	}
	b, _ := os.ReadFile(exe)
	if len(b) != len(nouveau) {
		t.Fatal("le programme n'a pas été remplacé")
	}
	if _, err := os.Stat(cheminAncien()); err != nil {
		t.Fatal("l'ancienne version devrait être conservée jusqu'au prochain lancement")
	}
	nettoyerMaj()
	if _, err := os.Stat(cheminAncien()); err == nil || majPrete() {
		t.Fatal("nettoyage incomplet")
	}
}
