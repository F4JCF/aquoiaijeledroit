# À quoi ai-je le droit ?

Logiciel Windows qui repère, en quelques questions, les aides et droits auxquels une personne peut prétendre : famille, minima sociaux, santé, logement, handicap, autonomie et aidants, retraite (dont CNRACL), emploi, **droit du travail** (privé et fonction publique) et impôts.

- **Aucune installation** : un seul fichier `AQuoiAiJeLeDroit.exe`, à lancer par double-clic.
- **Confidentiel** : les réponses restent sur l'ordinateur. Rien n'est envoyé sur Internet.
- **Toujours à jour** : les règles se mettent à jour automatiquement à chaque lancement.

> Repérage indicatif : seul l'organisme concerné peut confirmer un droit.

## Télécharger

Dernière version : onglet **Releases** du dépôt → `AQuoiAiJeLeDroit.exe`.

Au premier lancement, si Windows affiche « Windows a protégé votre ordinateur » : **Informations complémentaires → Exécuter quand même**.

---

## Comment fonctionnent les mises à jour

Il y a deux niveaux de mise à jour.

| Ce qui change | Ce qu'on modifie | Ce que voit l'utilisateur |
|---|---|---|
| Une loi, un plafond, une nouvelle aide, une aide supprimée | **`droits.json` seulement** | Au lancement suivant, les règles sont téléchargées et installées toutes seules, avec un bandeau « Les règles ont été mises à jour » et la liste des nouveautés. **Rien à réinstaller.** |
| Nouvelle question dans le formulaire, nouvelle fonction | Le code (`index.html`, `*.go`) + une nouvelle **release** | Un bandeau « Nouvelle version du logiciel disponible » avec un bouton de téléchargement. |

Fonctionnement technique :

- Le logiciel contient une copie de `droits.json` (fonctionne hors ligne).
- À chaque lancement, il lit `https://raw.githubusercontent.com/F4JCF/aquoiaijeledroit/main/droits.json`.
- Si la version en ligne est plus récente **et valide**, elle est enregistrée dans `%APPDATA%\AQuoiAiJeLeDroit\droits.json` et utilisée immédiatement.
- Un fichier en ligne illisible ou incomplet est **ignoré** : l'utilisateur garde les règles précédentes.
- Il consulte aussi la dernière release GitHub pour signaler une nouvelle version du `.exe`.

---

## Mettre à jour les règles après une nouvelle loi

1. Ouvrir `droits.json` sur GitHub (bouton crayon ✏️).
2. Modifier, ajouter ou supprimer le droit concerné (voir le format ci-dessous).
3. **Augmenter `version`** (par exemple `2026.10.13` → `2027.1.15`) et mettre à jour `mise_a_jour` (`2027-01-15`).
4. Ajouter une ligne en tête de `nouveautes` : elle s'affiche chez les utilisateurs.
   ```json
   {"date": "2027-01-15", "v": "2027.1.15", "texte": "Revalorisation de l'ASPA et nouvelle aide MaPrimeAdapt' élargie."}
   ```
5. Valider (« Commit changes »).

Le contrôle automatique (onglet **Actions**) vérifie le fichier. Une croix rouge ❌ signale une erreur : la corriger avant que les utilisateurs ne lancent le logiciel. De toute façon, un fichier invalide est ignoré par le logiciel.

### Format d'un droit

```json
{"id":"aspa", "d":"Ressources", "t":"Allocation de solidarité aux personnes âgées (ASPA)",
 "o":"Caisse de retraite", "l":"{SP}ASPA",
 "r":"Résumé du droit…", "a":"Démarche à faire…",
 "si": {"et":[{"c":"age",">=":65},{"c":"niv","<=?":2}]},
 "n":  [{"si":{"c":"niv","=":1},"n":"probable"},{"n":"possible"}],
 "w":  "{age} ans et ressources modestes"}
```

| Clé | Rôle |
|---|---|
| `id` | identifiant unique (sans espace) |
| `d` | domaine affiché : Famille, Ressources, Santé, Logement, Handicap, Autonomie, Retraite, Emploi & formation, Droit du travail, Impôts (un nouveau nom crée un nouveau domaine) |
| `t`, `o` | titre et organisme |
| `l` | lien ; `{SP}`, `{CAF}`, `{MDPH}`, `{RETRAITE}`, `{CNRACL}`, `{IMPOTS}`, `{AGEES}`, `{RENOV}` sont remplacés par les adresses de la section `liens` |
| `r`, `a` | résumé et démarche |
| `si` | condition pour afficher le droit |
| `n` | `probable`, `possible` ou `info` — directement, ou une liste de cas (le premier qui correspond gagne) |
| `w` | le « Pourquoi » affiché ; `{age}`, `{enf20}`, `{eTotalEff}`, `{tauxTxt}` sont remplacés par les réponses |

### Écrire une condition

- `"ald"` : la case est cochée ; `"!ald"` : elle ne l'est pas.
- `{"c":"age", ">=":65}` : comparaison avec `=`, `!=`, `>`, `>=`, `<`, `<=`, `dans`, `hors`.
- Ajouter `?` à l'opérateur (`">=?"`, `"<=?"`) accepte aussi une réponse **non renseignée**. C'est utile pour rester prudent quand la personne n'a pas donné ses revenus.
- `{"et":[…]}`, `{"ou":[…]}`, `{"non":…}` permettent de combiner.

### Champs disponibles

| Champ | Signification |
|---|---|
| `age` | âge |
| `sit` | `seul`, `couple`, `concubin`, `separe`, `veuf` |
| `act` | `salarie`, `fonct_t`, `fonct_e`, `fonct_h`, `contractuel`, `independant`, `chomeur`, `chomeur_fin`, `invalidite`, `retraite`, `etudiant`, `sans` |
| `log` | `loc_prive`, `loc_hlm`, `proprio`, `accedant`, `heberge`, `mobile` |
| `hand` | `non`, `encours`, `m50`, `50`, `80` ; `taux` = 0, 20, 50 ou 80 |
| `e0`, `e3`, `e6`, `e11`, `e18`, `e20` | enfants par tranche d'âge ; `e05` (0-5 ans), `e618` (6-17 ans), `enf18`, `enf20`, `enf` (total à charge), `eTotalEff` (enfants élevés) |
| `revenu`, `revenuPerso`, `ruc` | revenus du foyer, revenus personnels, revenu par unité de consommation |
| `niv` | niveau de ressources 1 à 4 (seuils dans `seuils.niveaux_ruc`) |
| cases à cocher | `grossesse`, `enfantHandicap`, `pasPension`, `tempsPartiel`, `carriereLongue`, `servicesActifs`, `independant2`, `procedure`, `logementAncien`, `adapter`, `difficulte`, `cherche`, `ald`, `autonomie`, `aidant`, `aideDomicile` |
| travail | `contrat` (`cdi`, `cdd`, `interim`, `apprenti`, `saisonnier`), `effectif` (`petite`, `moyenne`, `grande`), `anciennete` (années) ; cases `licenciement`, `depart`, `finContrat`, `heuresSup`, `salaireImpaye`, `harcelement`, `arretMaladie`, `accidentTravail`, `inaptitude`, `conges`, `naissance`, `reconversion`, `teletravail`, `horairesAtypiques`, `entrepriseDifficulte`, `mobilite` |
| déduits | `couple`, `isole` (parent seul), `fonct`, `travaille`, `actif`, `locataire`, `proprio`, `handDeclare`, `handOui`, `prive` (salarié du privé), `public` (fonctionnaire ou contractuel), `exSalarie` (demandeur d'emploi), `cdi`, `precaire` |

Un **nouveau champ** (nouvelle question) demande une nouvelle version du logiciel. Dans ce cas, mettre `app_min` dans `droits.json` à cette version : les anciennes versions garderont leurs règles et proposeront la mise à jour du logiciel.

---

## Publier une nouvelle version du logiciel

1. Modifier le code, puis changer `var version = "…"` dans `maj.go` (par exemple `1.3.0`).
2. Envoyer la modification sur `main`.

GitHub Actions compile `AQuoiAiJeLeDroit.exe`, puis **crée seul la release `v1.3.0`** avec le fichier. Les utilisateurs voient ensuite le bandeau de mise à jour. Si le numéro ne change pas, rien n'est publié : les simples corrections de `droits.json` ne créent donc pas de release.

## Développement

- `go test ./...` : tests, dont la validation de `droits.json`.
- `go run .` sous Linux ou macOS : interface servie sur http://127.0.0.1:8765.
