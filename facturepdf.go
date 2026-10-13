package main

// Lecture du texte d'une facture PDF, pour pré-remplir l'analyse de facture.
// Le fichier est lu sur l'ordinateur : rien n'est envoyé sur Internet.

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"

	"strings"

	"github.com/ledongthuc/pdf"
)

const textePDFMax = 300_000

func lireTextePDF(b []byte) (texte string, err error) {
	defer func() {
		if r := recover(); r != nil {
			texte, err = "", fmt.Errorf("PDF illisible")
		}
	}()
	if !bytes.HasPrefix(bytes.TrimLeft(b, " \r\n\t"), []byte("%PDF")) {
		return "", errors.New("ce fichier n'est pas un PDF")
	}
	r, err := pdf.NewReader(bytes.NewReader(b), int64(len(b)))
	if err != nil {
		return "", err
	}
	var sb strings.Builder
	for i := 1; i <= r.NumPage() && sb.Len() < textePDFMax; i++ {
		p := r.Page(i)
		if p.V.IsNull() {
			continue
		}
		lignes, err := p.GetTextByRow()
		if err != nil {
			continue
		}
		for _, l := range lignes {
			for j, w := range l.Content {
				if j > 0 {
					sb.WriteByte(' ')
				}
				sb.WriteString(w.S)
			}
			sb.WriteByte('\n')
		}
	}
	t := sb.String()
	if strings.TrimSpace(t) == "" {
		return "", errors.New("aucun texte trouvé (facture scannée ?)")
	}
	if len(t) > textePDFMax {
		t = t[:textePDFMax]
	}
	return t, nil
}

// lirePDFBase64 : appelée par l'interface avec le contenu du fichier encodé en base64.
func lirePDFBase64(b64 string) map[string]string {
	b, err := base64.StdEncoding.DecodeString(b64)
	if err != nil || len(b) > 30<<20 {
		return map[string]string{"erreur": "fichier illisible ou trop volumineux"}
	}
	t, err := lireTextePDF(b)
	if err != nil {
		return map[string]string{"erreur": err.Error()}
	}
	return map[string]string{"texte": t}
}
