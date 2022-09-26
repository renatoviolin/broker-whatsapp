package entity

type ListPayload struct {
	MessagingProduct string      `json:"messaging_product"`
	To               string      `json:"to"`
	Type             string      `json:"type"`
	Interactive      Interactive `json:"interactive"`
}

type Interactive struct {
	Type   string `json:"type"`
	Body   Body   `json:"body"`
	Action Action `json:"action"`
}

type Action struct {
	Button   string    `json:"button"`
	Sections []Section `json:"sections"`
}

type Section struct {
	Title string `json:"title"`
	Rows  []Row  `json:"rows"`
}

type Row struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type Body struct {
	Text string `json:"text"`
}

func NewListPayload(body string, rows []Row, to string) (payload ListPayload, err error) {
	if body == "" || len(rows) == 0 {
		return payload, ErrInvalidListPayload
	}
	payload.MessagingProduct = "whatsapp"
	payload.Type = "interactive"
	payload.To = to
	payload.Interactive.Type = "list"
	payload.Interactive.Body.Text = body
	payload.Interactive.Action.Button = "Opções"
	payload.Interactive.Action.Sections = []Section{
		{Title: "Agentes Disponíveis", Rows: rows},
	}
	return payload, nil
}
