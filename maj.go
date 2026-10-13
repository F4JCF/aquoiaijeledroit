package main

// Gestion des règles (droits.json) et des mises à jour.
//
// Les règles sont intégrées au logiciel, mais une version plus récente est
// téléchargée depuis le dépôt GitHub à chaque lancement, puis conservée dans
// %APPDATA%\AQuoiAiJeLeDroit\droits.json. Le logiciel lui-même signale
// lorsqu'une nouvelle version (release GitHub) est disponible.

import (
	_ "embed"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Renseignée à la compilation : -ldflags "-X main.version=1.2.0"
var version = "1.6.0"

const (
	depot          = "F4JCF/aquoiaijeledroit"
	schemaSupporte = 1
	tailleMax      = 4 << 20
)

var (
	urlRegles  = envOu("AQDROIT_URL_REGLES", "https://raw.githubusercontent.com/"+depot+"/main/droits.json")
	urlRelease = envOu("AQDROIT_URL_RELEASE", "https://api.github.com/repos/"+depot+"/releases/latest")
)

//go:embed droits.json
var reglesIntegrees []byte

//go:embed index.html
var pageModele string

//go:embed winres/i128.png
var icone []byte

func envOu(cle, defaut string) string {
	if v := os.Getenv(cle); v != "" {
		return v
	}
	return defaut
}

type entete struct {
	Schema     int               `json:"schema"`
	Version    string            `json:"version"`
	AppMin     string            `json:"app_min"`
	Droits     []json.RawMessage `json:"droits"`
	Nouveautes []json.RawMessage `json:"nouveautes"`
}

type droitMin struct {
	ID string          `json:"id"`
	T  string          `json:"t"`
	D  string          `json:"d"`
	Si json.RawMessage `json:"si"`
}

// comparerVersions renvoie -1, 0 ou 1 (versions pointées : 2026.10.13, 1.2.0…)
func comparerVersions(a, b string) int {
	pa, pb := strings.Split(strings.TrimPrefix(a, "v"), "."), strings.Split(strings.TrimPrefix(b, "v"), ".")
	for i := 0; i < len(pa) || i < len(pb); i++ {
		var x, y int
		if i < len(pa) {
			x, _ = strconv.Atoi(pa[i])
		}
		if i < len(pb) {
			y, _ = strconv.Atoi(pb[i])
		}
		if x != y {
			if x < y {
				return -1
			}
			return 1
		}
	}
	return 0
}

// valider vérifie qu'un fichier de règles est lisible et compatible avec ce logiciel.
func valider(b []byte) (entete, error) {
	var e entete
	if err := json.Unmarshal(b, &e); err != nil {
		return e, fmt.Errorf("fichier illisible : %w", err)
	}
	if e.Version == "" {
		return e, errors.New("version absente")
	}
	if e.Schema < 1 || e.Schema > schemaSupporte {
		return e, fmt.Errorf("format %d non pris en charge", e.Schema)
	}
	if len(e.Droits) < 10 {
		return e, errors.New("liste de droits incomplète")
	}
	vus := map[string]bool{}
	for i, r := range e.Droits {
		var d droitMin
		if err := json.Unmarshal(r, &d); err != nil || d.ID == "" || d.T == "" || d.D == "" || len(d.Si) == 0 {
			return e, fmt.Errorf("droit n°%d incomplet", i+1)
		}
		if vus[d.ID] {
			return e, fmt.Errorf("identifiant en double : %s", d.ID)
		}
		vus[d.ID] = true
	}
	return e, nil
}

func cheminCache() string {
	dir, err := os.UserConfigDir()
	if err != nil {
		return ""
	}
	return filepath.Join(dir, "AQuoiAiJeLeDroit", "droits.json")
}

// reglesCourantes : la version la plus récente et compatible entre l'intégrée et le cache.
func reglesCourantes() ([]byte, entete) {
	meilleur := reglesIntegrees
	e, _ := valider(reglesIntegrees)
	if c := cheminCache(); c != "" {
		if b, err := os.ReadFile(c); err == nil {
			if ec, err := valider(b); err == nil &&
				comparerVersions(version, ec.AppMin) >= 0 &&
				comparerVersions(ec.Version, e.Version) > 0 {
				return b, ec
			}
		}
	}
	return meilleur, e
}

// construirePage injecte les règles, l'icône et la version dans l'interface.
func construirePage() string {
	regles, _ := reglesCourantes()
	p := strings.Replace(pageModele, "/*REGLES*/null", string(regles), 1)
	p = strings.ReplaceAll(p, "{{ICONE}}", "data:image/png;base64,"+base64.StdEncoding.EncodeToString(icone))
	p = strings.ReplaceAll(p, "{{VERSION}}", version)
	return p
}

var client = &http.Client{Timeout: 10 * time.Second}

func telecharger(url string) ([]byte, error) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "AQuoiAiJeLeDroit/"+version)
	req.Header.Set("Cache-Control", "no-cache")
	r, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	if r.StatusCode != 200 {
		return nil, fmt.Errorf("réponse %d", r.StatusCode)
	}
	return io.ReadAll(io.LimitReader(r.Body, tailleMax))
}

type resultatMaj struct {
	Statut  string          `json:"statut"` // maj, a_jour, hors_ligne, invalide, app_trop_ancienne
	Message string          `json:"message,omitempty"`
	Version string          `json:"version,omitempty"`
	Regles  json.RawMessage `json:"regles,omitempty"`
	App     *infoApp        `json:"app,omitempty"`
}

type infoApp struct {
	Version string `json:"version"`
	URL     string `json:"url"`
}

// verifierMaj télécharge les règles et vérifie la présence d'une nouvelle version du logiciel.
func verifierMaj() resultatMaj {
	res := resultatMaj{}

	// 1. Nouvelle version du logiciel ?
	if b, err := telecharger(urlRelease); err == nil {
		var rel struct {
			Tag string `json:"tag_name"`
			URL string `json:"html_url"`
		}
		if json.Unmarshal(b, &rel) == nil && rel.Tag != "" && comparerVersions(rel.Tag, version) > 0 {
			res.App = &infoApp{Version: strings.TrimPrefix(rel.Tag, "v"), URL: rel.URL}
		}
	}

	// 2. Nouvelles règles ?
	b, err := telecharger(urlRegles)
	if err != nil {
		res.Statut, res.Message = "hors_ligne", "Impossible de joindre le serveur de mises à jour."
		return res
	}
	e, err := valider(b)
	if err != nil {
		res.Statut, res.Message = "invalide", "Le fichier de règles en ligne a été ignoré ("+err.Error()+")."
		return res
	}
	_, courant := reglesCourantes()
	if comparerVersions(e.Version, courant.Version) <= 0 {
		res.Statut, res.Version = "a_jour", courant.Version
		return res
	}
	if comparerVersions(version, e.AppMin) < 0 {
		res.Statut, res.Version = "app_trop_ancienne", e.Version
		res.Message = "De nouvelles règles existent mais nécessitent la version " + e.AppMin + " du logiciel."
		return res
	}
	if c := cheminCache(); c != "" {
		os.MkdirAll(filepath.Dir(c), 0o755)
		tmp := c + ".tmp"
		if os.WriteFile(tmp, b, 0o644) == nil {
			os.Rename(tmp, c)
		}
	}
	res.Statut, res.Version, res.Regles = "maj", e.Version, b
	return res
}

func verifierMajJSON() string {
	b, _ := json.Marshal(verifierMaj())
	return string(b)
}

func lienAutorise(u string) bool {
	return strings.HasPrefix(u, "https://") || strings.HasPrefix(u, "http://")
}

func jsChaine(s string) string {
	b, _ := json.Marshal(s)
	return string(b)
}
