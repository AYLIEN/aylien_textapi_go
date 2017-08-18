// +build appengine

package textapi

import (
	"context"

	"google.golang.org/appengine/urlfetch"
)

// NewClientInContext returns a new client using the given auth information & context, for use with GAE.
func NewClientInContext(auth Auth, useHTTPS bool, ctx context.Context) (*Client, error) {
	client, err := NewClient(auth, useHTTPS)
	if err != nil {
		return nil, err
	}
	client.Client = urlfetch.Client(ctx)
	return client, nil
}
