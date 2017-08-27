// +build appengine

package textapi

import (
	"context"
	"time"

	"google.golang.org/appengine/urlfetch"
)

// NewClientInContext returns a new client using the given auth information & context, for use with GAE.
func NewClientInContext(auth Auth, useHTTPS bool, ctx context.Context) (*Client, error) {
	client, err := NewClient(auth, useHTTPS)
	if err != nil {
		return nil, err
	}
	ctxTO, _ := context.WithTimeout(ctx, 60*time.Second)
	client.Client = urlfetch.Client(ctxTO)
	return client, nil
}
