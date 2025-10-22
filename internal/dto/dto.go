package dto

type HHResponse struct {
	Items []Vacancy `json:"items"`
	Pages int       `json:"pages"`
}

type Vacancy struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Area   Area   `json:"area"`
	Salary Salary `json:"salary"`
	Url    string `json:"alternate_url"`
}

type Area struct {
	Name string `json:"name"`
}

type Salary struct {
	From int `json:"from"`
	To   int `json:"to"`
}

type Update struct {
	UpdateID int `json:"update_id"`
	Message  struct {
		MessageID int `json:"message_id"`
		Chat      struct {
			ID int `json:"id"`
		} `json:"chat"`
		Text string `json:"text"`
		From struct {
			Username string `json:"username"`
		} `json:"from"`
	} `json:"message"`
}
