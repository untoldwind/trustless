package rest

// Link is the minimal implementation of a HATOAS link
type Link struct {
	Href   string `json:"href" xml:"href"`
	Method string `json:"method,omitempty" xml:"method,omitempty"`
}

// SimpleLink create a simple HATOAS link
func SimpleLink(href string) Link {
	return Link{
		Href: href,
	}
}
