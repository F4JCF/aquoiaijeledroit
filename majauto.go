package main

// Mise à jour automatique du logiciel lui-même.
//
// 1. Quand une nouvelle version est publiée sur GitHub, le .exe est téléchargé en
//    arrière-plan à côté du programme (AQuoiAiJeLeDroit.exe.new).
// 2. Il est vérifié : empreinte SHA-256 publiée par GitHub, taille, en-tête Windows.
// 3. Il remplace l'ancien au prochain démarrage, ou tout de suite sur demande.
//    L'ancien est gardé sous le nom .old puis supprimé au lancement suivant.

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const nomAsset = "AQuoiAiJeLeDroit.exe"

var (
	cheminExe = func() string {
		p, _ := os.Executable()
		r, err := filepath.EvalSymlinks(p)
		if err == nil {
			return r
		}
		return p
	}
	clientDL = &http.Client{Timeout: 5 * time.Minute}
	dlMu     sync.Mutex
)

type assetGitHub struct {
	Nom    string `json:"name"`
	URL    string `json:"browser_download_url"`
	Taille int64  `json:"size"`
	Digest string `json:"digest"` // « sha256:… »
}

func cheminNouveau() string { return cheminExe() + ".new" }
func cheminAncien() string  { return cheminExe() + ".old" }

// verifierBinaire contrôle un exécutable téléchargé.
func verifierBinaire(b []byte, a assetGitHub) error {
	if len(b) < 1<<20 {
		return errors.New("fichier trop petit")
	}
	if a.Taille > 0 && int64(len(b)) != a.Taille {
		return errors.New("taille inattendue")
	}
	if !bytes.HasPrefix(b, []byte("MZ")) {
		return errors.New("ce n'est pas un programme Windows")
	}
	if strings.HasPrefix(a.Digest, "sha256:") {
		s := sha256.Sum256(b)
		if hex.EncodeToString(s[:]) != strings.TrimPrefix(a.Digest, "sha256:") {
			return errors.New("empreinte SHA-256 incorrecte")
		}
	}
	return nil
}

// telechargerMaj récupère et vérifie la nouvelle version ; renvoie nil si elle est prête.
func telechargerMaj(a assetGitHub) error {
	dlMu.Lock()
	defer dlMu.Unlock()
	if a.URL == "" {
		return errors.New("aucun fichier")
	}
	req, _ := http.NewRequest("GET", a.URL, nil)
	req.Header.Set("User-Agent", "AQuoiAiJeLeDroit/"+version)
	r, err := clientDL.Do(req)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	if r.StatusCode != 200 {
		return fmt.Errorf("réponse %d", r.StatusCode)
	}
	b, err := io.ReadAll(io.LimitReader(r.Body, 64<<20))
	if err != nil {
		return err
	}
	if err := verifierBinaire(b, a); err != nil {
		return err
	}
	tmp := cheminNouveau() + ".part"
	if err := os.WriteFile(tmp, b, 0o755); err != nil {
		return err // dossier non modifiable : la mise à jour restera manuelle
	}
	return os.Rename(tmp, cheminNouveau())
}

// majPrete indique si une nouvelle version téléchargée attend d'être installée.
func majPrete() bool {
	st, err := os.Stat(cheminNouveau())
	return err == nil && st.Size() > 1<<20
}

// installerMaj remplace le programme par la version téléchargée.
// Windows autorise à renommer un programme en cours d'exécution.
func installerMaj() error {
	if !majPrete() {
		return errors.New("aucune mise à jour prête")
	}
	exe := cheminExe()
	os.Remove(cheminAncien())
	if err := os.Rename(exe, cheminAncien()); err != nil {
		return err
	}
	if err := os.Rename(cheminNouveau(), exe); err != nil {
		os.Rename(cheminAncien(), exe) // retour arrière
		return err
	}
	return nil
}

// relancer démarre le programme (nouvelle version) et laisse l'appelant se fermer.
func relancer() error {
	return exec.Command(cheminExe()).Start()
}

// nettoyerMaj supprime l'ancienne version laissée par une mise à jour précédente.
func nettoyerMaj() {
	os.Remove(cheminAncien())
	os.Remove(cheminNouveau() + ".part")
}
