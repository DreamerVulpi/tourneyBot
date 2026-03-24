package challonge

type ErrorResponse struct {
	Errors `json:"errors"`
}

type Errors struct {
	Detail string `json:"detail"`
	Status int    `json:"status"`
}
