package daemon_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/leanovate/microtools/logging"
	"github.com/stretchr/testify/require"
	"github.com/untoldwind/trustless/daemon"
)

func TestServiceResource(t *testing.T) {
	require := require.New(t)
	resource := daemon.NewRootResource(nil, logging.NewSimpleLoggerNull())

	request, _ := http.NewRequest("GET", "/", nil)

	result, err := resource.Get(request)
	require.Nil(err)
	doc, ok := result.(*daemon.ServiceDocument)
	require.True(ok)
	require.NotNil(doc.Links)
	require.NotEmpty(doc.Links["self"].Href)
	require.NotEmpty(doc.Links["v1"].Href)

	matcher := resource.SubResources()
	require.NotNil(matcher)

	request, _ = http.NewRequest("GET", "/v1", nil)
	response := httptest.NewRecorder()
	require.True(matcher("/v1", response, request))
	require.Equal(response.Code, 200)
}
