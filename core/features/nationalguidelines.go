package features

import (
	"encoding/json"
	"fmt"
)

type NationalGuidelinesFeatures struct {
	BaseFeatures
}

func NewNationalGuidelinesFeatures() NationalGuidelinesFeatures {
	return NationalGuidelinesFeatures{
		BaseFeatures{
			&Feature{
				Reference:   "Allegato 3, lett. f) DM 10/09/2010",
				Title:       "Siti UNESCO, di notevole interesse culturale e di notevole interesse pubblico ex art. 136 d.lgs. 42/2004",
				Description: `- *i siti inseriti nella lista del patrimonio mondiale dell'UNESCO, le aree ed i beni di notevole interesse culturale di cui alla Parte Seconda del d.lgs 42 del 2004, nonché gli immobili e le aree dichiarati di notevole interesse pubblico ai sensi dell'art. 136 dello stesso decreto legislativo;*`,
				Tags:        DM10sett2010,
			},
			&Feature{
				Reference:   "Allegato 3, lett. f) DM 10/09/2010",
				Title:       "Zone all'interno di coni visuali",
				Description: `- *zone all'interno di coni visuali la cui immagine è storicizzata e identifica i luoghi anche in termini di notorietà internazionale di attrattività turistica;*`,
				Tags:        DM10sett2010,
			},
			&Feature{
				Reference:   "Allegato 3, lett. f) DM 10/09/2010",
				Title:       "Zone in prossimità di parchi archeologici e beni di interesse culturale, storico e/o religioso",
				Description: `- *zone situate in prossimità di parchi archeologici e nelle aree contermini ad emergenze di particolare interesse culturale, storico e/o religioso;*`,
				Tags:        DM10sett2010,
			},
			&Feature{
				Reference:   "Allegato 3, lett. f) DM 10/09/2010",
				Title:       "Aree naturali protette",
				Description: `- *le aree naturali protette ai diversi livelli (nazionale, regionale, locale istituite ai sensi della Legge 394/91 ed inserite nell'Elenco Ufficiale delle Aree Naturali Protette, con particolare riferimento alle aree di riserva integrale e di riserva generale orientata di cui all'articolo 12, comma 2, lettere a) e b) della legge 394/91 ed equivalenti a livello regionale;*`,
				Tags:        EUAP,
			},
			&Feature{
				Reference:   "Allegato 3, lett. f) DM 10/09/2010",
				Title:       "Zone umide (Ramsar)",
				Description: `- *le zone umide di importanza internazionale designate ai sensi della Convenzione di Ramsar;*`,
				Tags:        Ramsar,
			},
			&Feature{
				Reference:   "Allegato 3, lett. f) DM 10/09/2010",
				Title:       "Aree della Rete Natura 2000 (SIC, ZPS)",
				Description: `- *le aree incluse nella Rete Natura 2000 designate in base alla Direttiva 92/43/CEE (Siti di importanza Comunitaria) ed alla Direttiva 79/409/CEE (Zone di Protezione Speciale);*`,
				Tags:        ReteNatura2000,
			},
			&Feature{
				Reference:   "Allegato 3, lett. f) DM 10/09/2010",
				Title:       "Important Bird Areas (IBA)",
				Description: `- *le Important Bird Areas (I.B.A.);*`,
				Tags:        DM10sett2010,
			},
			&Feature{
				Reference:   "Allegato 3, lett. f) DM 10/09/2010",
				Title:       "Aree determinanti per la conservazione della biodiversità",
				Description: `- *le aree non comprese in quelle di cui ai punti precedenti ma che svolgono funzioni determinanti per la conservazione della biodiversità (fasce di rispetto o aree contigue delle aree naturali protette; istituende aree naturali protette oggetto di proposta del Governo ovvero di disegno di legge regionale approvato dalla Giunta; aree di connessione e continuità ecologico-funzionale tra i vari sistemi naturali e seminaturali; aree di riproduzione, alimentazione e transito di specie faunistiche protette; aree in cui è accertata la presenza di specie animali e vegetali soggette a tutela dalle Convezioni internazionali (Berna, Bonn, Parigi, Washington, Barcellona) e dalle Direttive comunitarie (79/409/CEE e 92/43/CEE), specie rare, endemiche, vulnerabili, a rischio di estinzione;*`,
				Tags:        DM10sett2010,
			},
			&Feature{
				Reference:   "Allegato 3, lett. f) DM 10/09/2010",
				Title:       "Aree agricole interessate da produzioni di qualità",
				Description: `- *le aree agricole interessate da produzioni agricolo-alimentari di qualità (produzioni biologiche, produzioni D.O.P., I.G.P., S.T.G., D.O.C., D.O.C.G., produzioni tradizionali) e/o di particolare pregio rispetto al contesto paesaggistico-culturale, in coerenza e per le finalità di cui all'art. 12, comma 7, del decreto legislativo 387 del 2003 anche con riferimento alle aree, se previste dalla programmazione regionale, caratterizzate da un'elevata capacità d'uso del suolo;*`,
				Tags:        DM10sett2010,
			},
			&Feature{
				Reference:   "Allegato 3, lett. f) DM 10/09/2010",
				Title:       "Aree in situazione di dissesto o rischio idrogeologico",
				Description: `- *le aree caratterizzate da situazioni di dissesto e/o rischio idrogeologico perimetrate nei Piani di Assetto Idrogeologico (P.A.I.) adottati dalle competenti Autorità di Bacino ai sensi del D.L. 180/98 e s.m.i.;*`,
				Tags:        DM10sett2010,
			},
			&Feature{
				Reference:   "Allegato 3, lett. f) DM 10/09/2010",
				Title:       "Zone tutelate dall'art. 142 d.lgs. 42/2004",
				Description: `- *zone individuate ai sensi dell'art. 142 del d. lgs. 42 del 2004 valutando la sussistenza di particolari caratteristiche che le rendano incompatibili con la realizzazione degli impianti;*`,
				Tags:        Art142Dlgs2004,
			},
		},
	}
}

func (features NationalGuidelinesFeatures) Describe() []string {
	ret := make([]string, 0)
	applicableFeatures := features.ApplicableFeatures()
	if len(applicableFeatures) == 0 {
		ret = append(ret, "Il Progetto non ricade in aree classificate non idonee ai sensi dell'Allegato 3, lett. f) DM 10/09/2010 (Linee Guida Nazionali).")
		return ret
	}
	ret = append(ret, "Il Progetto ricade nelle seguenti aree, classificate non idonee ai sensi dell'Allegato 3, lett. f) DM 10/09/2010 (Linee Guida Nazionali):")

	for _, cond := range applicableFeatures {
		ret = append(ret, fmt.Sprintf("\"%v\"", cond.Description))
	}
	return ret
}

func (features NationalGuidelinesFeatures) MarshalJSON() ([]byte, error) {
	return json.Marshal(features.BaseFeatures)
}
