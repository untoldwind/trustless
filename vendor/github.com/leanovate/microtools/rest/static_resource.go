package rest

import "net/http"

type StaticResource struct {
	ResourceBase
	result *Result
	self   Link
}

func NewStaticResource(result *Result, self Link) Resource {
	return &StaticResource{
		result: result,
		self:   self,
	}
}

func StaticContent(content []byte, contentType string, selfHref string) Resource {
	return NewStaticResource(
		Ok().AddHeader("Content-Type", contentType).WithBody(content),
		SimpleLink(selfHref),
	)
}

func (r StaticResource) Get(request *http.Request) (interface{}, error) {
	return r.result, nil
}

func (r StaticResource) Self() Link {
	return r.self
}
