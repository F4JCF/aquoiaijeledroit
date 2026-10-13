package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"
)

func TestReglesIntegreesValides(t *testing.T) {
	e, err := valider(reglesIntegrees)
	if err != nil {
		t.Fatal(err)
	}
	if comparerVersions(version, e.AppMin) < 0 {
		t.Fatalf("droits.json exige le logiciel %s, version actuelle %s", e.AppMin, version)
	}
}

func TestComparerVersions(t *testing.T) {
	cas := []struct {
		a, b string
		r    int
	}{{"1.2.0", "1.1.9", 1}, {"v1.1.0", "1.1", 0}, {"2026.9.30", "2026.10.1", -1}, {"2026.10.13", "2026.10.13", 0}}
	for _, c := range cas {
		if g := comparerVersions(c.a, c.b); g != c.r {
			t.Errorf("%s vs %s = %d, attendu %d", c.a, c.b, g, c.r)
		}
	}
}

func serveur(regles, release string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "release") {
			w.Write([]byte(release))
			return
		}
		w.Write([]byte(regles))
	}))
}

func modifier(f func(m map[string]any)) string {
	var m map[string]any
	json.Unmarshal(reglesIntegrees, &m)
	f(m)
	b, _ := json.Marshal(m)
	return string(b)
}

func TestMiseAJour(t *testing.T) {
	os.Setenv("XDG_CONFIG_HOME", t.TempDir())
	os.Setenv("APPDATA", os.Getenv("XDG_CONFIG_HOME"))
	plusRecente := modifier(func(m map[string]any) { m["version"] = "2099.1.1" })
	tropExigeante := modifier(func(m map[string]any) { m["version"] = "2099.2.1"; m["app_min"] = "99.0.0" })
	casse := modifier(func(m map[string]any) { m["version"] = "2099.3.1"; m["droits"] = []any{} })

	essai := func(regles, release, statut string) resultatMaj {
		s := serveur(regles, release)
		defer s.Close()
		urlRegles, urlRelease = s.URL+"/droits.json", s.URL+"/release"
		r := verifierMaj()
		if r.Statut != statut {
			t.Fatalf("statut %q, attendu %q (%s)", r.Statut, statut, r.Message)
		}
		return r
	}
	essai(string(reglesIntegrees), `{}`, "a_jour")
	essai(casse, `{}`, "invalide")
	essai(tropExigeante, `{}`, "app_trop_ancienne")
	r := essai(plusRecente, `{"tag_name":"v9.0.0","html_url":"https://x"}`, "maj")
	if r.App == nil || r.App.Version != "9.0.0" {
		t.Fatal("nouvelle version du logiciel non détectée")
	}
	if _, e := reglesCourantes(); e.Version != "2099.1.1" {
		t.Fatal("le cache n'est pas utilisé au prochain lancement")
	}
	essai(plusRecente, `{}`, "a_jour")
	urlRegles = "http://127.0.0.1:1/x"
	if verifierMaj().Statut != "hors_ligne" {
		t.Fatal("hors ligne non géré")
	}
}

// Les lettres types ne doivent utiliser que des champs connus.
func TestLettresTypes(t *testing.T) {
	var d struct {
		Lettres []struct {
			ID, Titre, Objet, Dest, Ref, Corps string
			Champs                             []struct {
				ID, Label, Type string
				Options         []struct{ V, T, X string }
			}
		}
		Echeances []struct{ Date, Titre string }
		Annuaire  []struct{ Cat, Nom, Num string }
	}
	if err := json.Unmarshal(reglesIntegrees, &d); err != nil {
		t.Fatal(err)
	}
	base := map[string]bool{"nom": true, "adresse": true, "cpville": true, "tel": true, "email": true, "ville": true, "date": true}
	re := regexp.MustCompile(`\{(\w+)\}`)
	for _, l := range d.Lettres {
		if l.ID == "" || l.Titre == "" || l.Corps == "" {
			t.Fatalf("lettre incomplète : %+v", l.ID)
		}
		connus := map[string]bool{}
		for k := range base {
			connus[k] = true
		}
		for _, c := range l.Champs {
			connus[c.ID] = true
			if c.Type == "select" && len(c.Options) == 0 {
				t.Errorf("%s : liste %s vide", l.ID, c.ID)
			}
		}
		for _, txt := range []string{l.Objet, l.Dest, l.Ref, l.Corps} {
			for _, m := range re.FindAllStringSubmatch(txt, -1) {
				if !connus[m[1]] {
					t.Errorf("lettre %s : champ inconnu {%s}", l.ID, m[1])
				}
			}
		}
	}
	for _, e := range d.Echeances {
		if _, err := time.Parse("2006-01-02", e.Date); err != nil || e.Titre == "" {
			t.Errorf("échéance invalide : %+v", e)
		}
	}
	for _, a := range d.Annuaire {
		if a.Nom == "" || a.Num == "" || a.Cat == "" {
			t.Errorf("annuaire incomplet : %+v", a)
		}
	}
}
