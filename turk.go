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
		t := reflect.TypeOf(v).String()
		if t == "int32" {
			panic("int32 isn't supported. Please use int64. Use int64(x)")
		} else if t == "int" {
			panic("int isn't supported. Please use int64. Use int64(x)")
		} else if t == "float32" {
			panic("float32 isn't supported. Please use float64. Use float64(x.y)")
		}
		switch t {
		case "int64":
			q = strings.Replace(q, "{"+k+"}", strconv.FormatInt(v.(int64), 10), -1)
		case "float64":
			q = strings.Replace(q, "{"+k+"}", strconv.FormatFloat(v.(float64), 'f', 6, 64), -1)
		case "bool":
			q = strings.Replace(q, "{"+k+"}", strconv.FormatBool(v.(bool)), -1)
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
