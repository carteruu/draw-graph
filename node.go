package draw_graph

type (
	Node struct {
		ID     int    `json:"id"`
		Parent int    `json:"parent"`
		Text   string `json:"text"`
	}
)
