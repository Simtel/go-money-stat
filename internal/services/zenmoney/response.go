package zenmoney

type Response struct {
	Account     []Account     `json:"account"`
	Instrument  []Instrument  `json:"instrument"`
	Transaction []Transaction `json:"transaction"`
	Tag         []Tag         `json:"tag"`
}

func (r *Response) GetIndexedTags() map[string]Tag {
	tags := make(map[string]Tag)
	for _, tag := range r.Tag {
		tags[tag.Id] = tag
	}
	return tags
}
