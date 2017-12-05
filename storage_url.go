package gcp

import (
	"fmt"
	"net/url"

	"cloud.google.com/go/storage"
)

// ObjectFromURL parses a URI of the form gs://bucket/path/to/object and returns a GCS ObjectHandle.
func ObjectFromURL(client *storage.Client, uri string) (*storage.ObjectHandle, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	if u.Scheme != "gs" {
		return nil, fmt.Errorf("unknown scheme %s://; must be gs://", u.Scheme)
	}
	bucket := u.Host
	path := u.Path
	return client.Bucket(bucket).Object(path[1:]), nil
}
