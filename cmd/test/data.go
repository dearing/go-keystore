package main

type MTGCard struct {
	Artist                string   `json:"artist"`
	ArtistIds             []string `json:"artistIds"`
	Availability          []string `json:"availability"`
	BorderColor           string   `json:"borderColor"`
	ColorIdentity         []string `json:"colorIdentity"`
	Colors                []string `json:"colors"`
	ConvertedManaCost     float64  `json:"convertedManaCost"`
	EdhrecRank            int      `json:"edhrecRank"`
	FaceConvertedManaCost float64  `json:"faceConvertedManaCost"`
	FaceManaValue         float64  `json:"faceManaValue"`
	FaceName              string   `json:"faceName"`
	Finishes              []string `json:"finishes"`
	FlavorText            string   `json:"flavorText"`
	ForeignData           []any    `json:"foreignData"`
	FrameEffects          []string `json:"frameEffects"`
	FrameVersion          string   `json:"frameVersion"`
	HasFoil               bool     `json:"hasFoil"`
	HasNonFoil            bool     `json:"hasNonFoil"`
	Identifiers           struct {
		MtgjsonV4ID            string `json:"mtgjsonV4Id"`
		ScryfallID             string `json:"scryfallId"`
		ScryfallIllustrationID string `json:"scryfallIllustrationId"`
		ScryfallOracleID       string `json:"scryfallOracleId"`
		TcgplayerProductID     string `json:"tcgplayerProductId"`
	} `json:"identifiers"`
	IsStarter  bool     `json:"isStarter"`
	Keywords   []string `json:"keywords"`
	Language   string   `json:"language"`
	Layout     string   `json:"layout"`
	Legalities struct {
	} `json:"legalities"`
	ManaCost     string   `json:"manaCost"`
	ManaValue    float64  `json:"manaValue"`
	Name         string   `json:"name"`
	Number       string   `json:"number"`
	OtherFaceIds []string `json:"otherFaceIds"`
	Printings    []string `json:"printings"`
	PurchaseUrls struct {
		Tcgplayer string `json:"tcgplayer"`
	} `json:"purchaseUrls"`
	Rarity        string   `json:"rarity"`
	SecurityStamp string   `json:"securityStamp"`
	SetCode       string   `json:"setCode"`
	Side          string   `json:"side"`
	Subtypes      []any    `json:"subtypes"`
	Supertypes    []any    `json:"supertypes"`
	Text          string   `json:"text"`
	Type          string   `json:"type"`
	Types         []string `json:"types"`
	UUID          string   `json:"uuid"`
}
