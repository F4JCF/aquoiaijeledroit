package main

// Création de documents PDF (lettres types, dossier, analyse de facture, simulation de retraite)
// et enregistrement dans le dossier Documents\À quoi ai-je le droit.

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/go-pdf/fpdf"
)

var (
	//go:embed polices/DejaVuSansCondensed.ttf
	policeSans []byte
	//go:embed polices/DejaVuSansCondensed-Bold.ttf
	policeSansGras []byte
	//go:embed polices/DejaVuSansCondensed-Oblique.ttf
	policeSansItal []byte
	//go:embed polices/DejaVuSerif.ttf
	policeSerif []byte
	//go:embed polices/DejaVuSerif-Bold.ttf
	policeSerifGras []byte
	//go:embed polices/DejaVuSerif-Italic.ttf
	policeSerifItal []byte
)

// Bloc : un élément du document. Types : h1, h2, p, b (gras), i (italique), petit,
// puce, exp (expéditeur), dest (destinataire), droite, sep (trait), espace.
type Bloc struct {
	T string `json:"t"`
	X string `json:"x"`
}

type Document struct {
	Titre   string `json:"titre"`
	Fichier string `json:"fichier"`
	Style   string `json:"style"` // "lettre" ou "dossier"
	Pied    string `json:"pied"`
	Blocs   []Bloc `json:"blocs"`
}

// nettoyerTexte retire les caractères absents des polices (émojis, pictogrammes).
func nettoyerTexte(s string) string {
	var b strings.Builder
	for _, r := range s {
		switch {
		case r == '\t':
			b.WriteRune(' ')
		case r > 0xFFFF, r >= 0x2600 && r <= 0x27BF, r >= 0x2B00 && r <= 0x2BFF, r == 0xFE0F, r == 0x200D:
			continue
		case unicode.IsControl(r) && r != '\n':
			continue
		default:
			b.WriteRune(r)
		}
	}
	return strings.TrimRight(b.String(), " ")
}

func genererPDF(d Document) ([]byte, error) {
	if len(d.Blocs) == 0 {
		return nil, errors.New("document vide")
	}
	p := fpdf.New("P", "mm", "A4", "")
	p.AddUTF8FontFromBytes("sans", "", policeSans)
	p.AddUTF8FontFromBytes("sans", "B", policeSansGras)
	p.AddUTF8FontFromBytes("sans", "I", policeSansItal)
	p.AddUTF8FontFromBytes("serif", "", policeSerif)
	p.AddUTF8FontFromBytes("serif", "B", policeSerifGras)
	p.AddUTF8FontFromBytes("serif", "I", policeSerifItal)
	police, corps, inter := "sans", 10.5, 5.2
	marge := 18.0
	if d.Style == "lettre" {
		police, corps, inter, marge = "serif", 11, 5.6, 22
	}
	p.SetMargins(marge, 18, marge)
	p.SetAutoPageBreak(true, 20)
	p.SetTitle(nettoyerTexte(d.Titre), true)
	p.SetCreator("À quoi ai-je le droit ?", true)
	pied := nettoyerTexte(d.Pied)
	p.SetFooterFunc(func() {
		p.SetY(-13)
		p.SetFont("sans", "", 8)
		p.SetTextColor(110, 110, 110)
		if pied != "" {
			p.CellFormat(0, 4, pied, "", 0, "L", false, 0, "")
		}
		p.SetX(-marge - 20)
		p.CellFormat(20, 4, fmt.Sprintf("%d/{nb}", p.PageNo()), "", 0, "R", false, 0, "")
	})
	p.AliasNbPages("{nb}")
	p.AddPage()
	larg, _ := p.GetPageSize()
	utile := larg - 2*marge
	p.SetTextColor(20, 20, 20)
	for _, bl := range d.Blocs {
		x := nettoyerTexte(bl.X)
		switch bl.T {
		case "h1":
			p.SetFont("sans", "B", 17)
			p.SetTextColor(20, 82, 79)
			p.MultiCell(utile, 8, x, "", "L", false)
			p.SetTextColor(20, 20, 20)
			p.Ln(1)
		case "h2":
			p.Ln(2)
			p.SetFont("sans", "B", 12.5)
			p.SetTextColor(20, 82, 79)
			p.MultiCell(utile, 6.5, x, "", "L", false)
			p.SetDrawColor(200, 195, 185)
			p.Line(marge, p.GetY()+0.5, larg-marge, p.GetY()+0.5)
			p.SetTextColor(20, 20, 20)
			p.Ln(2)
		case "b":
			if d.Style != "lettre" {
				p.Ln(1.5)
			}
			p.SetFont(police, "B", corps)
			p.MultiCell(utile, inter, x, "", "L", false)
		case "i":
			p.SetFont(police, "I", corps)
			p.MultiCell(utile, inter, x, "", "L", false)
		case "petit":
			p.SetFont("sans", "", 8.5)
			p.SetTextColor(90, 90, 90)
			p.MultiCell(utile, 4.2, x, "", "L", false)
			p.SetTextColor(20, 20, 20)
		case "puce":
			p.SetFont(police, "", corps)
			p.SetX(marge + 3)
			p.CellFormat(4, inter, "•", "", 0, "L", false, 0, "")
			p.MultiCell(utile-7, inter, x, "", "L", false)
		case "exp":
			p.SetFont(police, "", corps)
			p.MultiCell(utile/2, inter, x, "", "L", false)
		case "dest":
			p.Ln(6)
			p.SetFont(police, "", corps)
			for _, l := range strings.Split(x, "\n") {
				p.SetX(marge + utile/2)
				p.MultiCell(utile/2, inter, l, "", "L", false)
			}
			p.Ln(6)
		case "droite":
			p.SetFont(police, "", corps)
			p.MultiCell(utile, inter, x, "", "R", false)
			p.Ln(3)
		case "sep":
			p.Ln(2)
			p.SetDrawColor(200, 195, 185)
			p.Line(marge, p.GetY(), larg-marge, p.GetY())
			p.Ln(3)
		case "espace":
			p.Ln(inter)
		default: // "p"
			p.SetFont(police, "", corps)
			p.MultiCell(utile, inter, x, "", "L", false)
			if d.Style != "lettre" {
				p.Ln(1)
			}
		}
	}
	var buf bytes.Buffer
	if err := p.Output(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

var nomInterdit = regexp.MustCompile(`[<>:"/\\|?*\x00-\x1f]+`)

func nomFichierSur(s string) string {
	s = strings.TrimSpace(nomInterdit.ReplaceAllString(nettoyerTexte(s), " "))
	s = strings.Join(strings.Fields(s), " ")
	if len([]rune(s)) > 80 {
		s = string([]rune(s)[:80])
	}
	if s == "" {
		s = "document"
	}
	return s
}

func dossierDocuments() string {
	if p := os.Getenv("AQDROIT_DOCS"); p != "" {
		return p
	}
	h, _ := os.UserHomeDir()
	base := filepath.Join(h, "Documents")
	if st, err := os.Stat(base); err != nil || !st.IsDir() {
		base = h
	}
	return filepath.Join(base, "À quoi ai-je le droit")
}

// enregistrerDocument crée le PDF et renvoie son chemin complet.
func enregistrerDocument(docJSON string) map[string]string {
	var d Document
	if err := json.Unmarshal([]byte(docJSON), &d); err != nil {
		return map[string]string{"erreur": "document illisible"}
	}
	b, err := genererPDF(d)
	if err != nil {
		return map[string]string{"erreur": err.Error()}
	}
	dir := dossierDocuments()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return map[string]string{"erreur": "impossible de créer le dossier " + dir}
	}
	nom := nomFichierSur(d.Fichier) + " - " + time.Now().Format("2006-01-02")
	chemin := filepath.Join(dir, nom+".pdf")
	for i := 2; ; i++ {
		if _, err := os.Stat(chemin); os.IsNotExist(err) {
			break
		}
		chemin = filepath.Join(dir, fmt.Sprintf("%s (%d).pdf", nom, i))
	}
	if err := os.WriteFile(chemin, b, 0o644); err != nil {
		return map[string]string{"erreur": err.Error()}
	}
	return map[string]string{"chemin": chemin, "dossier": dir, "nom": filepath.Base(chemin)}
}

// cheminDocumentAutorise : seuls les fichiers créés dans le dossier du logiciel peuvent être ouverts.
func cheminDocumentAutorise(p string) bool {
	dir, _ := filepath.Abs(dossierDocuments())
	abs, err := filepath.Abs(p)
	return err == nil && strings.HasPrefix(abs, dir+string(os.PathSeparator)) && strings.EqualFold(filepath.Ext(abs), ".pdf")
}
