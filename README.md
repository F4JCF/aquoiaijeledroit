# À quoi ai-je le droit ?

Logiciel Windows qui repère, en quelques questions, les aides et droits auxquels une personne peut prétendre : famille, minima sociaux, santé, logement, handicap, autonomie et aidants, retraite (dont CNRACL), emploi, **droit du travail** (privé et fonction publique), **droit pénal** (victimes, mis en cause, casier), **code de la route** (permis, amendes, points, accident, véhicule) et impôts.

- **Aucune installation** : un seul fichier `AQuoiAiJeLeDroit.exe`, à lancer par double-clic.
- **Pratique** : lettres types pré-remplies, suivi des démarches avec relances, simulateur de retraite, recherche directe, échéances avec rappels, numéros utiles, plusieurs profils, sauvegarde, grands caractères et lecture à voix haute.
- **Sur téléphone** : version web sur https://f4jcf.github.io/aquoiaijeledroit/.
- **Confidentiel** : les réponses restent sur l'ordinateur. Rien n'est envoyé sur Internet. Les réponses de l'étape « Justice » ne sont même jamais enregistrées.
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

## Mettre à jour les prix de l'énergie

Les prix du comparateur sont dans la section `energie` de `droits.json`. Ils sont donc mis à jour comme les règles, sans nouvelle version du logiciel.

- `energie.date` et `energie.source` sont affichés au-dessus du comparateur.
- `elec.trv` contient le tarif réglementé : abonnement annuel TTC par puissance (`abo`), prix Base par puissance (`base`), `hp` et `hc`. Il change en février et en août.
- `elec.offres` contient les offres : `base`, `hp`, `hc` (€/kWh TTC pour 6 kVA), `abo_delta` (écart d'abonnement annuel avec le tarif réglementé, par option `base` / `hphc`), `type` (`fixe`, `indexe`, `variable`), `cond`, `url`. Pour les puissances de 9 kVA et plus, le prix Base est ajusté dans la même proportion que le tarif réglementé.
- `gaz.repere` contient le prix repère CRE, publié chaque mois. Au-dessus de `seuil` kWh, c'est le profil `chauffage` qui s'applique, en dessous le profil `cuisson`.
- `gaz.offres` contient les offres : `cuisson` et/ou `chauffage`, chacun avec `abo` (€/an) et `kwh`.

Source conseillée : le comparateur du Médiateur national de l'énergie (comparateur-offres.energie-info.fr).

### Télécoms et eau

- `telecom.mobile.offres` : `operateur`, `nom`, `reseau`, `data` (Go ; 9999 = illimité), `5g`, `prix`, `promo_mois` et `prix_apres` si le prix augmente, `frais` (carte SIM), `age_max`, `cond`, `url`.
- `telecom.box.offres` : mêmes champs, plus `debit`, `tv`, `age_min`/`age_max`. Le classement se fait sur le coût total sur 24 mois, frais compris.
- `eau` : `prix_moyen_m3` (moyenne nationale SISPEA) et `conso_personne` (m³ par personne et par an).

### Pièces à préparer, lettres types, échéances, numéros utiles

- `pj` (dans un droit) : liste des pièces à préparer, affichée dans les résultats et dans le dossier imprimé. `pieces_generales` : pièces souvent demandées.
- `lettres` : modèles de courrier. Chaque lettre a `id`, `titre`, `desc`, `recommande` (lettre recommandée), `si` (condition qui la met en avant ⭐), `champs` (`id`, `label`, `type` : `text`, `date`, `number`, `textarea`, `select` avec `options` `{v, t, x}` où `x` est le texte inséré), `dest`, `objet`, `ref`, `corps`, `conseil`. Dans les textes, `{nom}`, `{adresse}`, `{cpville}`, `{tel}`, `{email}`, `{ville}`, `{date}` et `{id_du_champ}` sont remplacés. Le test `go test` vérifie qu'aucun champ inconnu n'est utilisé.
- `echeances` : dates nationales `{date: "AAAA-MM-JJ", titre, info, si, lien, approx}`, affichées si la condition `si` correspond au profil. Les rappels personnels sont saisis dans le logiciel.
- `annuaire` : numéros utiles `{cat, nom, num, info, url}`.

- `retraite` : simulateur de retraite. `generations` : `{du, au, age_a, age_m, trim, cl}` (âge légal en années et mois, trimestres requis, carrière longue en texte) ; `taux_plein_auto` ; `progressive_min` ; `source`.
- Dans `lettres`, `relance` (en jours) fixe la date de relance proposée quand l'utilisateur indique avoir envoyé la lettre.

## Version téléphone (site web)

À chaque modification de `main`, GitHub Actions génère une version site web (`go run . -site public`) et la publie sur GitHub Pages : https://f4jcf.github.io/aquoiaijeledroit/. Elle lit les règles à jour directement sur GitHub et enregistre les réponses dans le navigateur du téléphone.

Réglage unique à faire dans le dépôt : **Settings → Pages → Build and deployment → Source : « GitHub Actions »**.

## Signature du logiciel (à activer si un certificat est acheté)

Sans signature, Windows affiche « Windows a protégé votre ordinateur » au premier lancement. Pour la supprimer, il faut un certificat de signature de code, par exemple le service *Trusted Signing* de Microsoft Azure (abonnement mensuel), puis ajouter une étape de signature dans `.github/workflows/compiler.yml`, avec les identifiants enregistrés comme secrets du dépôt.

## Données de l'utilisateur

Les réponses, les profils, les coordonnées des lettres et les rappels sont enregistrés uniquement sur l'ordinateur, dans `%APPDATA%\AQuoiAiJeLeDroit\donnees.json`. Les réponses de l'étape « Justice » ne sont jamais enregistrées.

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
| `d` | domaine affiché : Famille, Ressources, Santé, Logement, Handicap, Autonomie, Retraite, Emploi & formation, Énergie, Eau, Télécoms, Assurances, Banque, Consommation & garanties, Jeunes, Seniors, Décès & succession, Séparation, Bénévolat & dons, Micro-entreprise, Santé, Outre-mer, Droit du travail, Droit pénal, Code de la route, Impôts (un nouveau nom crée un nouveau domaine) |
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
| énergie | `chauffage` (`elec`, `pac`, `gaz`, `fioul`, `bois`, `collectif`), `consoElec`, `consoGaz` ; cases `factureElevee`, `changerFournisseur`, `impayeEnergie`, `coupure`, `litigeEnergie`, `demarchage`, `linky`, `changerChauffage`, `solaire` ; déduits `chauffageElec`, `ecoElec` (économie estimée en € par le comparateur) |
| fournisseurs | eau : `factureEau`, `consoEau`, cases `eauFuite`, `eauImpaye`, `eauCoupure`, `eauLitige` ; télécoms : `prixMobile`, `prixBox`, cases `mobileHausse`, `mobileResilier`, `mobileLitige`, `couverture`, `cuivre`, `boxHausse`, `boxResilier`, `boxPanne`, `boxLitige` ; assurances : cases `assurAuto`, `assurHabitation`, `assurSante`, `assurEmprunteur`, `assurObseques`, `assurHausse`, `assurResilier`, `assurSinistre`, `assurCatnat` ; déduits `prixEau`, `eauCher`, `eauSurconso`, `ecoMobile`, `ecoBox`, `telecomLitige`, `telecomHausse`, `telecomResilier` |
| logement (locataire / propriétaire) | cases `depotGarantie`, `preavis`, `chargesLoc`, `indecent`, `bailleur`, `copro` |
| banque et consommation | cases `occasionAchat`, `vicesCaches`, `extensionGarantie`, `reparation`, `livraison`, `achatPanne`, `achatDistance`, `litigeConso`, `abonnementPiege`, `fraisBancaires`, `fraudeBancaire`, `refusCompte`, `changerBanque`, `credit`, `surendettement`, `fichage` |
| étapes de vie | jeunes : `etudiantSup`, `logementEtudiant`, `jeuneSansEmploi`, `jeuneCulture` ; seniors : `ehpad`, `domicile`, `transportSenior`, `maltraitance` ; décès : `decesProche`, `decesConjoint`, `succession`, `comptesDefunt`, `donation` |
| nouveaux (v1.8) | `territoire` (`metropole`, `corse`, `guadeloupe`, `martinique`, `guyane`, `reunion`, `mayotte`) ; cases `microProjet`, `microActif`, `activiteCote`, `sansMedecin`, `depassements`, `renoncement`, `transportSanitaire`, `cure`, `dossierMedical`, `separation`, `pensionAlim`, `gardeEnfants`, `logementFamilial`, `benevole`, `dirigeantAsso`, `dons` ; déduits `outremer`, `zni` (Corse ou outre-mer), `micro` |
| route | `permis` (`aucun`, `probatoire`, `confirme`, `suspendu`) ; cases `amende`, `points`, `suspension`, `alcoolStup`, `accidentAuto`, `stationnement`, `fourriere`, `carteGrise`, `assuranceAuto`, `zfe`, `plaqueUsurpee`, `passerPermis`, `deuxRoues`, `santeConduite` |
| justice | `victime`, `violencesFamille`, `arnaque`, `harcelementLigne`, `accidentRoute`, `convoque`, `poursuivi`, `temoin`, `casier`, `procheDetenu` — **jamais enregistrés** sur l'ordinateur |
| déduits | `couple`, `isole` (parent seul), `fonct`, `travaille`, `actif`, `locataire`, `proprio`, `handDeclare`, `handOui`, `prive` (salarié du privé), `public` (fonctionnaire ou contractuel), `exSalarie` (demandeur d'emploi), `cdi`, `precaire`, `victimeQ` (victime, quelle que soit l'infraction), `misEnCause` (convoqué ou poursuivi) |

Un **nouveau champ** (nouvelle question) demande une nouvelle version du logiciel. Dans ce cas, mettre `app_min` dans `droits.json` à cette version : les anciennes versions garderont leurs règles et proposeront la mise à jour du logiciel.

---

## Publier une nouvelle version du logiciel

1. Modifier le code, puis changer `var version = "…"` dans `maj.go` (par exemple `1.3.0`).
2. Envoyer la modification sur `main`.

GitHub Actions compile `AQuoiAiJeLeDroit.exe`, puis **crée seul la release `v1.3.0`** avec le fichier. Les utilisateurs voient ensuite le bandeau de mise à jour. Si le numéro ne change pas, rien n'est publié : les simples corrections de `droits.json` ne créent donc pas de release.

## Développement

- `go test ./...` : tests, dont la validation de `droits.json`.
- `go run .` sous Linux ou macOS : interface servie sur http://127.0.0.1:8765.
