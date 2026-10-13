package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDocumentPDF(t *testing.T) {
	dir := t.TempDir()
	os.Setenv("AQDROIT_DOCS", dir)
	defer os.Unsetenv("AQDROIT_DOCS")
	d := Document{Titre: "Lettre", Fichier: "Réclamer: dépôt/garantie", Style: "lettre", Pied: "Créé avec À quoi ai-je le droit ?",
		Blocs: []Bloc{{"exp", "Jean Dupont\n12 rue des Fleurs\n06250 Mougins"}, {"dest", "M. Martin\n3 avenue du Port"}, {"droite", "À Mougins, le 13 octobre 2026"},
			{"b", "Objet : Mise en demeure ✉️"}, {"p", "Madame, Monsieur,\n\nLe dépôt de garantie de 900 € ne m'a pas été restitué. Œuvre « déjà » été."}, {"espace", ""}, {"p", "Jean Dupont"}}}
	j, _ := json.Marshal(d)
	r := enregistrerDocument(string(j))
	if r["erreur"] != "" {
		t.Fatal(r["erreur"])
	}
	if filepath.Dir(r["chemin"]) != dir || !strings.HasPrefix(r["nom"], "Réclamer dépôt garantie - ") {
		t.Fatalf("chemin inattendu : %v", r)
	}
	b, _ := os.ReadFile(r["chemin"])
	txt, err := lireTextePDF(b)
	if err != nil {
		t.Fatal(err)
	}
	for _, m := range []string{"dépôt de garantie", "900", "Mougins", "« déjà »", "Objet : Mise en demeure"} {
		if !strings.Contains(txt, m) {
			t.Errorf("texte absent du PDF : %q\n%s", m, txt)
		}
	}
	if r2 := enregistrerDocument(string(j)); !strings.Contains(r2["nom"], "(2)") {
		t.Errorf("le second fichier aurait dû être numéroté : %v", r2)
	}
	if !cheminDocumentAutorise(r["chemin"]) || cheminDocumentAutorise("/etc/passwd") || cheminDocumentAutorise(filepath.Join(dir, "x.exe")) {
		t.Error("contrôle des chemins incorrect")
	}
}
