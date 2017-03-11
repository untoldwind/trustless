package rest

type Link struct {
	Href   string `json:"href" xml:"href"`
	Method string `json:"method,omitempty" xml:"method,omitempty"`
}

func SimpleLink(href string) Link {
	return Link{
		Href: href,
	}
}
