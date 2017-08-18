// +build !appengine

package textapi

import (
	"context"
)

// NewClientInContext returns a new client using the given auth information, 
// ignoring context, for use outside GAE.
func NewClientInContext(auth Auth, useHTTPS bool, ctx context.Context) (*Client, error) {
	return  NewClient(auth, useHTTPS)
}
