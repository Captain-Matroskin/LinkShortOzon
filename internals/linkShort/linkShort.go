package linkShort

type LinkFull struct {
	Link string `json:"link"`
}

type LinkShort struct {
	Link string `json:"link"`
}

type ResponseLinkShort struct {
	LinkShort LinkShort `json:"link_short"`
}

type ResponseLinkFull struct {
	LinkShort LinkFull `json:"link_full"`
}

