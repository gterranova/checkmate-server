package features

import (
	"encoding/json"
	"fmt"
)

type SuitabilityFeatures struct {
	BaseFeatures
}

func NewSuitabilityFeatures() SuitabilityFeatures {
	return SuitabilityFeatures{
		BaseFeatures{
			&Feature{
				Reference:   "Art. 20, comma 8 lett. a), D.lgs. 199/2021",
				Title:       "Aree occupate da impianti della stessa fonte per interventi di modifica",
				Description: `a) i siti ove sono già installati impianti della stessa fonte e in cui vengono realizzati interventi di modifica, anche sostanziale, per rifacimento, potenziamento o integrale ricostruzione, eventualmente abbinati a sistemi di accumulo, che non comportino una variazione dell'area occupata superiore al 20 per cento. Il limite percentuale di cui al primo periodo non si applica per gli impianti fotovoltaici, in relazione ai quali la variazione dell'area occupata è soggetta al limite di cui alla lettera c-ter), numero 1);`,
				Tags:        Art20,
			},
			&Feature{
				Reference:   "Art. 20, comma 8 lett. b), D.lgs. 199/2021",
				Title:       "Siti oggetto di bonifica",
				Description: `b) le aree dei siti oggetto di bonifica individuate ai sensi del Titolo V, Parte quarta, del decreto legislativo 3 aprile 2006, n. 152;`,
				Tags:        Art20,
			},
			&Feature{
				Reference:   "Art. 20, comma 8 lett. c), D.lgs. 199/2021",
				Title:       "Cave e miniere cessate",
				Description: `c) le cave e miniere cessate, non recuperate o abbandonate o in condizioni di degrado ambientale, o le porzioni di cave e miniere non suscettibili di ulteriore sfruttamento;`,
				Tags:        Art20,
			},
			&Feature{
				Reference:   "Art. 20, comma 8 lett. c-bis), D.lgs. 199/2021",
				Title:       "Aree RFI e concessionari autostradali",
				Description: `c-bis) i siti e gli impianti nelle disponibilità delle società del gruppo Ferrovie dello Stato italiane e dei gestori di infrastrutture ferroviarie nonché delle società concessionarie autostradali;`,
				Tags:        Art20,
			},
			&Feature{
				Reference:   "Art. 20, comma 8 lett. c-bis.1), D.lgs. 199/2021",
				Title:       "Aree aeroportuali",
				Description: `c-bis.1) i siti e gli impianti nella disponibilità delle società di gestione aeroportuale all'interno dei sedimi aeroportuali, ivi inclusi quelli all'interno del perimetro di pertinenza degli aeroporti delle isole minori di cui all'allegato 1 al decreto del Ministro dello sviluppo economico 14 febbraio 2017, pubblicato nella Gazzetta Ufficiale n. 114 del 18 maggio 2017, ferme restando le necessarie verifiche tecniche da parte dell'Ente nazionale per l'aviazione civile (ENAC);`,
				Tags:        Art20,
			},
			&Feature{
				Reference:   "Art. 20, comma 8 lett. c-ter) sub 1), D.lgs. 199/2021",
				Title:       "Aree agricole entro 500m da zone industriali, SIN, cave e miniere",
				Description: `1) le aree classificate agricole, racchiuse in un perimetro i cui punti distino non più di 500 metri da zone a destinazione industriale, artigianale e commerciale, compresi i siti di interesse nazionale, nonché le cave e le miniere`,
				Tags:        Art20,
			},
			&Feature{
				Reference:   "Art. 20, comma 8 lett. c-ter) sub 2), D.lgs. 199/2021",
				Title:       "Aree interne a impianti industriali e stabilimenti o agricole entro 500m",
				Description: `2) le aree interne agli impianti industriali e agli stabilimenti, questi ultimi come definiti dall'articolo 268, comma 1, lettera h), del decreto legislativo 3 aprile 2006, n. 152, nonché le aree classificate agricole racchiuse in un perimetro i cui punti distino non più di 500 metri dal medesimo impianto o stabilimento`,
				Tags:        Art20,
			},
			&Feature{
				Reference:   "Art. 20, comma 8 lett. c-ter) sub 3), D.lgs. 199/2021",
				Title:       "Aree entro 300m da autostrade",
				Description: `3) le aree adiacenti alla rete autostradale entro una distanza non superiore a 300 metri`,
				Tags:        Art20,
			},
			&Feature{
				Reference:   "Art. 20, comma 8 lett. c-quater), D.lgs. 199/2021",
				Title:       "Altre aree (c-quater)",
				Description: `c-quater) fatto salvo quanto previsto alle lettere a), b), c), c-bis) e c-ter), le aree che non sono ricomprese nel perimetro dei beni sottoposti a tutela ai sensi del decreto legislativo 22 gennaio 2004, n. 42, incluse le zone gravate da usi civici di cui all'articolo 142, comma 1, lettera h), del medesimo decreto, né ricadono nella fascia di rispetto dei beni sottoposti a tutela ai sensi della parte seconda oppure dell'articolo 136 del medesimo decreto legislativo. Ai soli fini della presente lettera, la fascia di rispetto è determinata considerando una distanza dal perimetro di beni sottoposti a tutela di tre chilometri per gli impianti eolici e di cinquecento metri per gli impianti fotovoltaici. Resta ferma, nei procedimenti autorizzatori, la competenza del Ministero della cultura a esprimersi in relazione ai soli progetti localizzati in aree sottoposte a tutela secondo quanto previsto all'articolo 12, comma 3-bis, del decreto legislativo 29 dicembre 2003, n. 387`,
				Tags:        Art20,
			},
		},
	}
}

func (features SuitabilityFeatures) Describe() []string {
	ret := make([]string, 0)
	applicableFeatures := features.ApplicableFeatures()
	if len(applicableFeatures) == 0 {
		ret = append(ret, "Il Progetto **non ricade in aree classificate idonee** ai sensi dell'articolo 20 del decreto legislativo 8 novembre 2021, n. 199.")
		return ret
	}
	ret = append(ret, "Il Progetto ricade nelle seguenti aree, classificate idonee ai sensi dell'articolo 20 del decreto legislativo 8 novembre 2021, n. 199:")

	for _, cond := range applicableFeatures {
		ret = append(ret, fmt.Sprintf("*%v* (%v)", cond.Description, cond.Reference))
	}
	return ret
}

func (features SuitabilityFeatures) MarshalJSON() ([]byte, error) {
	return json.Marshal(features.BaseFeatures)
}
