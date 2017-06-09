package turk

import (
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/sethgrid/pester"
)

/**
*  Created by Galileo on 9/6/17.
*
*  A super simple GraphQL client
 */

// GraphQL represents the query and its props
type GraphQL struct {
	Statement  string
	Parameters map[string]interface{}
}

// Props is used to replace parameters in the query
type Props map[string]interface{}

// Build returns the GraphQL query string after parsing the Props
func (g GraphQL) Build() string {
	q := g.Statement
	for k, v := range g.Parameters {
		switch reflect.TypeOf(k).String() {
		case "int", "int32", "int64", "bool", "float64", "float32":
			q = strings.Replace(q, "{"+k+"}", v.(string), -1)
		case "string":
			q = strings.Replace(q, "{"+k+"}", `"`+v.(string)+`"`, -1)
		default:
			q = strings.Replace(q, "{"+k+"}", `"`+v.(string)+`"`, -1)
		}
	}
	return q
}

// Turk contains the GraphQL server URL, query and a pester HTTP client
type Turk struct {
	url   string
	query *string
	p     *pester.Client
}

// NewTurkClient returns a Turk instance
func NewTurkClient(url string, query *string) *Turk {

	client := pester.New()
	// for resiliency
	client.Concurrency = 3
	client.MaxRetries = 5
	client.Backoff = pester.ExponentialJitterBackoff

	return &Turk{
		query: query,
		p:     client,
		url:   url,
	}
}

// Send 'POST's the GraphQL query using pester HTTP client
func (g *Turk) Send() (*http.Response, error) {
	var str = `{"query":` + strconv.QuoteToASCII(*g.query) + `}`
	return g.p.Post(g.url, "application/graphql", strings.NewReader(str))
}
