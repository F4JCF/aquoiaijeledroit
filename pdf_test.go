package main

import (
	"os"
	"strings"
	"testing"
)

func TestLirePDF(t *testing.T) {
	p := os.Getenv("AQDROIT_PDF_TEST")
	if p == "" {
		t.Skip("pas de PDF de test")
	}
	b, _ := os.ReadFile(p)
	txt, err := lireTextePDF(b)
	if err != nil || !strings.Contains(txt, "kVA") {
		t.Fatalf("%v %q", err, txt)
	}
	t.Log(txt)
	if _, err := lireTextePDF([]byte("pas un pdf")); err == nil {
		t.Fatal("un faux PDF a été accepté")
	}
}
