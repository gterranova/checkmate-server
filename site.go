package main

import (
	"fmt"
)

type ConditionValue uint

const (
	NotApplicable ConditionValue = iota
	Art22bis
	Art20
	DM10sett2010

	EUAP
	Ramsar
	ReteNatura2000
	Art142Dlgs2004
	Art142Dlgs2004Coste
	Art142Dlgs2004MontagneBoschi
)

type Condition struct {
	Title       string
	Description string
	Reference   string
	Checked     bool
	Disabled    bool
	Value       ConditionValue
}

type LocationConditions []*Condition
type SuitabilityConditions []*Condition
type NationalGuidelinesConditions []*Condition
type ScreeningGuidelinesConditions []*Condition

func NewLocationConditions() LocationConditions {
	return LocationConditions{
		&Condition{
			Reference:   "",
			Title:       "Aree agricole",
			Description: ``,
			Value:       NotApplicable,
		},
		&Condition{
			Reference:   "Art. 22-bis, comma 1, D.lgs. 199/2021",
			Title:       "Aree a destinazione industriale, artigianale e commerciale",
			Description: `*L'installazione, con qualunque modalità, di impianti fotovoltaici su terra e delle relative opere connesse e infrastrutture necessarie, ubicati nelle zone e nelle aree a destinazione industriale, artigianale e commerciale* [...] *è considerata attività di manutenzione ordinaria e non è subordinata acquisizione di permessi, autorizzazioni o atti di assenso comunque denominati, fatte salve le valutazioni ambientali di cui al titolo III della parte seconda del decreto legislativo 3 aprile 2006, n. 152, ove previste.*`,
			Value:       Art22bis,
		},
		&Condition{
			Reference:   "Art. 22-bis, comma 1, D.lgs. 199/2021",
			Title:       "Discariche o lotti di discarica chiusi e ripristinati",
			Description: `*L'installazione, con qualunque modalità, di impianti fotovoltaici su terra e delle relative opere connesse e infrastrutture necessarie, ubicati* [...] *in discariche o lotti di discarica chiusi e ripristinati* [...] *è considerata attività di manutenzione ordinaria e non è subordinata acquisizione di permessi, autorizzazioni o atti di assenso comunque denominati, fatte salve le valutazioni ambientali di cui al titolo III della parte seconda del decreto legislativo 3 aprile 2006, n. 152, ove previste.*`,
			Value:       Art22bis,
		},
		&Condition{
			Reference:   "Art. 22-bis, comma 1, D.lgs. 199/2021",
			Title:       "Cave o lotti o porzioni di cave non suscettibili di ulteriore sfruttamento",
			Description: `*L'installazione, con qualunque modalità, di impianti fotovoltaici su terra e delle relative opere connesse e infrastrutture necessarie, ubicati* [...] *in cave o lotti o porzioni di cave non suscettibili di ulteriore sfruttamento, è considerata attività di manutenzione ordinaria e non è subordinata acquisizione di permessi, autorizzazioni o atti di assenso comunque denominati, fatte salve le valutazioni ambientali di cui al titolo III della parte seconda del decreto legislativo 3 aprile 2006, n. 152, ove previste.*`,
			Value:       Art22bis,
		},
		&Condition{
			Reference:   "",
			Title:       "Altro",
			Description: ``,
			Value:       NotApplicable,
		},
	}
}

func (locationConditions *LocationConditions) ApplicableCondition() *Condition {
	for _, cond := range *locationConditions {
		if !cond.Disabled && cond.Checked {
			return cond
		}
	}
	return nil
}

func (locationConditions *LocationConditions) Art22bisArea() ConditionValue {
	if cond := locationConditions.ApplicableCondition(); cond != nil {
		return cond.Value
	}
	return NotApplicable
}

func (locationConditions *LocationConditions) Describe() []string {
	ret := make([]string, 0)
	if cond := locationConditions.ApplicableCondition(); cond != nil {
		if len(cond.Reference) > 0 {
			ret = append(ret, fmt.Sprintf("Ai sensi dell'%v, \"%v\".", cond.Reference, cond.Description))
		}
	}
	return ret
}

func NewSuitabilityConditions() SuitabilityConditions {
	return SuitabilityConditions{
		&Condition{
			Reference:   "Art. 20, comma 8 lett. a), D.lgs. 199/2021",
			Title:       "Aree occupate da impianti della stessa fonte per interventi di modifica",
			Description: `a) i siti ove sono già installati impianti della stessa fonte e in cui vengono realizzati interventi di modifica, anche sostanziale, per rifacimento, potenziamento o integrale ricostruzione, eventualmente abbinati a sistemi di accumulo, che non comportino una variazione dell'area occupata superiore al 20 per cento. Il limite percentuale di cui al primo periodo non si applica per gli impianti fotovoltaici, in relazione ai quali la variazione dell'area occupata è soggetta al limite di cui alla lettera c-ter), numero 1);`,
			Value:       Art20,
		},
		&Condition{
			Reference:   "Art. 20, comma 8 lett. b), D.lgs. 199/2021",
			Title:       "Siti oggetto di bonifica",
			Description: `b) le aree dei siti oggetto di bonifica individuate ai sensi del Titolo V, Parte quarta, del decreto legislativo 3 aprile 2006, n. 152;`,
			Value:       Art20,
		},
		&Condition{
			Reference:   "Art. 20, comma 8 lett. c), D.lgs. 199/2021",
			Title:       "Cave e miniere cessate",
			Description: `c) le cave e miniere cessate, non recuperate o abbandonate o in condizioni di degrado ambientale, o le porzioni di cave e miniere non suscettibili di ulteriore sfruttamento;`,
			Value:       Art20,
		},
		&Condition{
			Reference:   "Art. 20, comma 8 lett. c-bis), D.lgs. 199/2021",
			Title:       "Aree RFI e concessionari autostradali",
			Description: `c-bis) i siti e gli impianti nelle disponibilità delle società del gruppo Ferrovie dello Stato italiane e dei gestori di infrastrutture ferroviarie nonché delle società concessionarie autostradali;`,
			Value:       Art20,
		},
		&Condition{
			Reference:   "Art. 20, comma 8 lett. c-bis.1), D.lgs. 199/2021",
			Title:       "Aree aeroportuali",
			Description: `c-bis.1) i siti e gli impianti nella disponibilità delle società di gestione aeroportuale all'interno dei sedimi aeroportuali, ivi inclusi quelli all'interno del perimetro di pertinenza degli aeroporti delle isole minori di cui all'allegato 1 al decreto del Ministro dello sviluppo economico 14 febbraio 2017, pubblicato nella Gazzetta Ufficiale n. 114 del 18 maggio 2017, ferme restando le necessarie verifiche tecniche da parte dell'Ente nazionale per l'aviazione civile (ENAC);`,
			Value:       Art20,
		},
		&Condition{
			Reference:   "Art. 20, comma 8 lett. c-ter) sub 1), D.lgs. 199/2021",
			Title:       "Aree agricole entro 500m da zone industriali, SIN, cave e miniere",
			Description: `1) le aree classificate agricole, racchiuse in un perimetro i cui punti distino non più di 500 metri da zone a destinazione industriale, artigianale e commerciale, compresi i siti di interesse nazionale, nonché le cave e le miniere`,
			Value:       Art20,
		},
		&Condition{
			Reference:   "Art. 20, comma 8 lett. c-ter) sub 2), D.lgs. 199/2021",
			Title:       "Aree interne a impianti industriali e stabilimenti o agricole entro 500m",
			Description: `2) le aree interne agli impianti industriali e agli stabilimenti, questi ultimi come definiti dall'articolo 268, comma 1, lettera h), del decreto legislativo 3 aprile 2006, n. 152, nonché le aree classificate agricole racchiuse in un perimetro i cui punti distino non più di 500 metri dal medesimo impianto o stabilimento`,
			Value:       Art20,
		},
		&Condition{
			Reference:   "Art. 20, comma 8 lett. c-ter) sub 3), D.lgs. 199/2021",
			Title:       "Aree entro 300m da autostrade",
			Description: `3) le aree adiacenti alla rete autostradale entro una distanza non superiore a 300 metri`,
			Value:       Art20,
		},
		&Condition{
			Reference:   "Art. 20, comma 8 lett. c-quater), D.lgs. 199/2021",
			Title:       "Altre aree (c-quater)",
			Description: `c-quater) fatto salvo quanto previsto alle lettere a), b), c), c-bis) e c-ter), le aree che non sono ricomprese nel perimetro dei beni sottoposti a tutela ai sensi del decreto legislativo 22 gennaio 2004, n. 42, incluse le zone gravate da usi civici di cui all'articolo 142, comma 1, lettera h), del medesimo decreto, né ricadono nella fascia di rispetto dei beni sottoposti a tutela ai sensi della parte seconda oppure dell'articolo 136 del medesimo decreto legislativo. Ai soli fini della presente lettera, la fascia di rispetto è determinata considerando una distanza dal perimetro di beni sottoposti a tutela di tre chilometri per gli impianti eolici e di cinquecento metri per gli impianti fotovoltaici. Resta ferma, nei procedimenti autorizzatori, la competenza del Ministero della cultura a esprimersi in relazione ai soli progetti localizzati in aree sottoposte a tutela secondo quanto previsto all'articolo 12, comma 3-bis, del decreto legislativo 29 dicembre 2003, n. 387`,
			Value:       Art20,
		},
	}
}

func (suitabilityConditions *SuitabilityConditions) ApplicableConditions() []*Condition {
	conditions := make([]*Condition, 0)
	for _, cond := range *suitabilityConditions {
		if !cond.Disabled && cond.Checked {
			conditions = append(conditions, cond)
		}
	}
	return conditions
}

func (suitabilityConditions *SuitabilityConditions) SuitableArea() ConditionValue {
	for _, cond := range suitabilityConditions.ApplicableConditions() {
		return cond.Value
	}
	return NotApplicable
}

func (suitabilityConditions *SuitabilityConditions) Describe() []string {
	ret := make([]string, 0)
	conditions := suitabilityConditions.ApplicableConditions()
	if len(conditions) == 0 {
		ret = append(ret, "Il Progetto **non ricade in aree classificate idonee** ai sensi dell'articolo 20 del decreto legislativo 8 novembre 2021, n. 199.")
		return ret
	}
	ret = append(ret, "Il Progetto ricade nelle seguenti aree, classificate idonee ai sensi dell'articolo 20 del decreto legislativo 8 novembre 2021, n. 199:")

	for _, cond := range conditions {
		ret = append(ret, fmt.Sprintf("*%v* (%v)", cond.Description, cond.Reference))
	}
	return ret
}

func NewNationalGuidelinesConditions() NationalGuidelinesConditions {
	return NationalGuidelinesConditions{
		&Condition{
			Reference:   "Allegato 3, lett. f) DM 10/09/2010",
			Title:       "Siti UNESCO, di notevole interesse culturale e di notevole interesse pubblico ex art. 136 d.lgs. 42/2004",
			Description: `- *i siti inseriti nella lista del patrimonio mondiale dell'UNESCO, le aree ed i beni di notevole interesse culturale di cui alla Parte Seconda del d.lgs 42 del 2004, nonché gli immobili e le aree dichiarati di notevole interesse pubblico ai sensi dell'art. 136 dello stesso decreto legislativo;*`,
			Value:       DM10sett2010,
		},
		&Condition{
			Reference:   "Allegato 3, lett. f) DM 10/09/2010",
			Title:       "Zone all'interno di coni visuali",
			Description: `- *zone all'interno di coni visuali la cui immagine è storicizzata e identifica i luoghi anche in termini di notorietà internazionale di attrattività turistica;*`,
			Value:       DM10sett2010,
		},
		&Condition{
			Reference:   "Allegato 3, lett. f) DM 10/09/2010",
			Title:       "Zone in prossimità di parchi archeologici e beni di interesse culturale, storico e/o religioso",
			Description: `- *zone situate in prossimità di parchi archeologici e nelle aree contermini ad emergenze di particolare interesse culturale, storico e/o religioso;*`,
			Value:       DM10sett2010,
		},
		&Condition{
			Reference:   "Allegato 3, lett. f) DM 10/09/2010",
			Title:       "Aree naturali protette",
			Description: `- *le aree naturali protette ai diversi livelli (nazionale, regionale, locale istituite ai sensi della Legge 394/91 ed inserite nell'Elenco Ufficiale delle Aree Naturali Protette, con particolare riferimento alle aree di riserva integrale e di riserva generale orientata di cui all'articolo 12, comma 2, lettere a) e b) della legge 394/91 ed equivalenti a livello regionale;*`,
			Value:       EUAP,
		},
		&Condition{
			Reference:   "Allegato 3, lett. f) DM 10/09/2010",
			Title:       "Zone umide (Ramsar)",
			Description: `- *le zone umide di importanza internazionale designate ai sensi della Convenzione di Ramsar;*`,
			Value:       Ramsar,
		},
		&Condition{
			Reference:   "Allegato 3, lett. f) DM 10/09/2010",
			Title:       "Aree della Rete Natura 2000 (SIC, ZPS)",
			Description: `- *le aree incluse nella Rete Natura 2000 designate in base alla Direttiva 92/43/CEE (Siti di importanza Comunitaria) ed alla Direttiva 79/409/CEE (Zone di Protezione Speciale);*`,
			Value:       ReteNatura2000,
		},
		&Condition{
			Reference:   "Allegato 3, lett. f) DM 10/09/2010",
			Title:       "Important Bird Areas (IBA)",
			Description: `- *le Important Bird Areas (I.B.A.);*`,
			Value:       DM10sett2010,
		},
		&Condition{
			Reference:   "Allegato 3, lett. f) DM 10/09/2010",
			Title:       "Aree determinanti per la conservazione della biodiversità",
			Description: `- *le aree non comprese in quelle di cui ai punti precedenti ma che svolgono funzioni determinanti per la conservazione della biodiversità (fasce di rispetto o aree contigue delle aree naturali protette; istituende aree naturali protette oggetto di proposta del Governo ovvero di disegno di legge regionale approvato dalla Giunta; aree di connessione e continuità ecologico-funzionale tra i vari sistemi naturali e seminaturali; aree di riproduzione, alimentazione e transito di specie faunistiche protette; aree in cui è accertata la presenza di specie animali e vegetali soggette a tutela dalle Convezioni internazionali (Berna, Bonn, Parigi, Washington, Barcellona) e dalle Direttive comunitarie (79/409/CEE e 92/43/CEE), specie rare, endemiche, vulnerabili, a rischio di estinzione;*`,
			Value:       DM10sett2010,
		},
		&Condition{
			Reference:   "Allegato 3, lett. f) DM 10/09/2010",
			Title:       "Aree agricole interessate da produzioni di qualità",
			Description: `- *le aree agricole interessate da produzioni agricolo-alimentari di qualità (produzioni biologiche, produzioni D.O.P., I.G.P., S.T.G., D.O.C., D.O.C.G., produzioni tradizionali) e/o di particolare pregio rispetto al contesto paesaggistico-culturale, in coerenza e per le finalità di cui all'art. 12, comma 7, del decreto legislativo 387 del 2003 anche con riferimento alle aree, se previste dalla programmazione regionale, caratterizzate da un'elevata capacità d'uso del suolo;*`,
			Value:       DM10sett2010,
		},
		&Condition{
			Reference:   "Allegato 3, lett. f) DM 10/09/2010",
			Title:       "Aree in situazione di dissesto o rischio idrogeologico",
			Description: `- *le aree caratterizzate da situazioni di dissesto e/o rischio idrogeologico perimetrate nei Piani di Assetto Idrogeologico (P.A.I.) adottati dalle competenti Autorità di Bacino ai sensi del D.L. 180/98 e s.m.i.;*`,
			Value:       DM10sett2010,
		},
		&Condition{
			Reference:   "Allegato 3, lett. f) DM 10/09/2010",
			Title:       "Zone tutelate dall'art. 142 d.lgs. 42/2004",
			Description: `- *zone individuate ai sensi dell'art. 142 del d. lgs. 42 del 2004 valutando la sussistenza di particolari caratteristiche che le rendano incompatibili con la realizzazione degli impianti;*`,
			Value:       Art142Dlgs2004,
		},
	}
}

func (nationalGuidelinesConditions *NationalGuidelinesConditions) ApplicableConditions() []*Condition {
	conditions := make([]*Condition, 0)
	for _, cond := range *nationalGuidelinesConditions {
		if !cond.Disabled && cond.Checked {
			conditions = append(conditions, cond)
		}
	}
	return conditions
}

func (nationalGuidelinesConditions *NationalGuidelinesConditions) NonUnsuitableDM2010Area() ConditionValue {
	for _, cond := range nationalGuidelinesConditions.ApplicableConditions() {
		return cond.Value
	}
	return NotApplicable
}

func (nationalGuidelinesConditions *NationalGuidelinesConditions) Describe() []string {
	ret := make([]string, 0)
	conditions := nationalGuidelinesConditions.ApplicableConditions()
	if len(conditions) == 0 {
		ret = append(ret, "Il Progetto non ricade in aree classificate non idonee ai sensi dell'Allegato 3, lett. f) DM 10/09/2010 (Linee Guida Nazionali).")
		return ret
	}
	ret = append(ret, "Il Progetto ricade nelle seguenti aree, classificate non idonee ai sensi dell'Allegato 3, lett. f) DM 10/09/2010 (Linee Guida Nazionali):")

	for _, cond := range conditions {
		ret = append(ret, fmt.Sprintf("\"%v\"", cond.Description))
	}
	return ret
}

func NewScreeningGuidelinesConditions() ScreeningGuidelinesConditions {
	return ScreeningGuidelinesConditions{
		&Condition{
			Reference:   "Linee Guida Screening DM 30/03/2015",
			Title:       "4.1. Cumulo con altri progetti",
			Description: ``,
			Value:       DM10sett2010,
		},
		&Condition{
			Reference:   "Linee Guida Screening DM 30/03/2015",
			Title:       "4.2. Rischio di incidenti, per quanto riguarda, in particolare, le sostanze o le tecnologie utilizzate",
			Description: ``,
			Value:       DM10sett2010,
		},
		&Condition{
			Reference:   "Linee Guida Screening DM 30/03/2015",
			Title:       "4.3.1. Zone umide",
			Description: `zone umide di importanza internazionale (Ramsar)`,
			Value:       Ramsar,
		},
		&Condition{
			Reference:   "Linee Guida Screening DM 30/03/2015",
			Title:       "4.3.2. Zone costiere",
			Description: `vincoli di cui al Codice dei beni culturali e del paesaggio (art. 142) - Aree di rispetto coste e corpi idrici`,
			Value:       Art142Dlgs2004Coste,
		},
		&Condition{
			Reference:   "Linee Guida Screening DM 30/03/2015",
			Title:       "4.3.3. Zone montuose e forestali",
			Description: `vincoli di cui al Codice dei beni culturali e del paesaggio (art. 142) - Montagne oltre 1600 o 1200 metri; piano forestale regionale/provinciale; in assenza di piano forestale vedi vincoli di cui al Codice dei beni culturali e del paesaggio (art. 142) - Boschi`,
			Value:       Art142Dlgs2004MontagneBoschi,
		},

		&Condition{
			Reference:   "Linee Guida Screening DM 30/03/2015",
			Title:       "4.3.4. Riserve e parchi naturali, zone classificate o protette ai sensi della normativa nazionale",
			Description: `Elenco ufficiale aree naturali protette (EUAP)`,
			Value:       EUAP,
		},
		&Condition{
			Reference:   "Linee Guida Screening DM 30/03/2015",
			Title:       "4.3.5. Zone protette speciali designate ai sensi delle direttive 2009/147/CE e 92/43/CEE",
			Description: `Siti di importanza comunitaria (SIC), Zone di protezione speciale (ZPS)`,
			Value:       ReteNatura2000,
		},
		&Condition{
			Reference:   "Linee Guida Screening DM 30/03/2015",
			Title:       "4.3.6. Zone nelle quali gli standard di qualità ambientale fissati dalla normativa dell'Unione europea sono già stati superati",
			Description: `dati di qualità delle acque superficiali e sotterranee`,
			Value:       DM10sett2010,
		},
		&Condition{
			Reference:   "Linee Guida Screening DM 30/03/2015",
			Title:       "4.3.7. Zone a forte densità demografica",
			Description: `zone a forte densità demografica si intendono i centri abitati, così come delimitati dagli strumenti urbanistici comunali, posti all'interno dei territori comunali con densità superiore a 500 abitanti per km2 e popolazione di almeno 50.000 abitanti (EUROSTAT)`,
			Value:       DM10sett2010,
		},
		&Condition{
			Reference:   "Linee Guida Screening DM 30/03/2015",
			Title:       "4.3.8. Zone di importanza storica, culturale o archeologica",
			Description: `immobili e le aree di cui all'art. 136 del Codice dei beni culturali e del paesaggio di cui al decreto legislativo n. 42/2004 dichiarati di notevole interesse pubblico ai sensi dell'art. 140 del medesimo decreto e gli immobili e le aree di interesse artistico, storico, archeologico o etnoantropologico di cui all'art. 10, comma 3, lettera a), del medesimo decreto`,
			Value:       DM10sett2010,
		},
	}
}

func (screeningGuidelinesConditions *ScreeningGuidelinesConditions) ApplicableConditions() []*Condition {
	conditions := make([]*Condition, 0)
	for _, cond := range *screeningGuidelinesConditions {
		if !cond.Disabled && cond.Checked {
			conditions = append(conditions, cond)
		}
	}
	return conditions
}

func (screeningGuidelinesConditions *ScreeningGuidelinesConditions) ShouldHalveThreshold() ConditionValue {
	for _, cond := range screeningGuidelinesConditions.ApplicableConditions() {
		return cond.Value
	}
	return NotApplicable
}

func (screeningGuidelinesConditions *ScreeningGuidelinesConditions) Describe() []string {
	ret := make([]string, 0)
	conditions := screeningGuidelinesConditions.ApplicableConditions()
	if len(conditions) == 0 {
		return ret
	}
	ret = append(ret, "Ai sensi delle Linee Guida per la verifica di assoggettabilità a VIA, Allegate al DM 30/03/2015, la **soglia per lo screening è soggetta a dimezzamento** in applicazione dei seguenti criteri:")

	for _, cond := range conditions {
		if len(cond.Description) > 0 {
			ret = append(ret, fmt.Sprintf("- %v (%v)", cond.Title, cond.Description))
		} else {
			ret = append(ret, fmt.Sprintf("- %v", cond.Title))
		}
	}
	return ret
}
