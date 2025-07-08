package handlers

// PageItem represents a single page item in the pagination.
type PageItem struct {
	Href  string `json:"href"`
	Title string `json:"title"`
}

// Pagination is a struct that represents pagination information.
type Pagination struct {
	PathParts []*PageItem `json:"path_parts,omitempty"`
	Total     int         `json:"total"`
	Count     int         `json:"count"`
	First     *PageItem   `json:"first,omitempty"`
	Prev      *PageItem   `json:"prev,omitempty"`
	Next      *PageItem   `json:"next,omitempty"`
	Last      *PageItem   `json:"last,omitempty"`
}
