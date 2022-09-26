package dto

type WebhookOutput struct {
	InternalStatusCode int    `json:"internal_status"`
	Message            string `json:"message"`
}

type WebhookInput struct {
	Object string  `json:"object"`
	Entry  []Entry `json:"entry"`
}

type Entry struct {
	ID      string   `json:"id"`
	Changes []Change `json:"changes"`
}

type Change struct {
	Value Value  `json:"value"`
	Field string `json:"field"`
}

type Value struct {
	MessagingProduct string    `json:"messaging_product"`
	Metadata         Metadata  `json:"metadata"`
	Contacts         []Contact `json:"contacts"`
	Messages         []Message `json:"messages"`
	Statuses         []Status  `json:"statuses"`
}

type Contact struct {
	Profile Profile `json:"profile"`
	WaID    string  `json:"wa_id"`
}

type Profile struct {
	Name string `json:"name"`
}

type Message struct {
	Context     Context     `json:"context"`
	From        string      `json:"from"`
	ID          string      `json:"id"`
	Timestamp   string      `json:"timestamp"`
	Text        Text        `json:"text"`
	Type        string      `json:"type"`
	Interactive Interactive `json:"interactive"`
}

type Text struct {
	Body string `json:"body"`
}

type Metadata struct {
	DisplayPhoneNumber string `json:"display_phone_number"`
	PhoneNumberID      string `json:"phone_number_id"`
}

type Status struct {
	ID           string       `json:"id"`
	Status       string       `json:"status"`
	Timestamp    string       `json:"timestamp"`
	RecipientID  string       `json:"recipient_id"`
	Conversation Conversation `json:"conversation"`
	Pricing      Pricing      `json:"pricing"`
}

type Conversation struct {
	ID     string `json:"id"`
	Origin Origin `json:"origin"`
}

type Origin struct {
	Type string `json:"type"`
}

type Pricing struct {
	Billable     bool   `json:"billable"`
	PricingModel string `json:"pricing_model"`
	Category     string `json:"category"`
}

type Context struct {
	From string `json:"from"`
	ID   string `json:"id"`
}

type Interactive struct {
	Type      string    `json:"type"`
	ListReply ListReply `json:"list_reply"`
}

type ListReply struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}
