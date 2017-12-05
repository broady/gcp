package gcp

import (
	"testing"

	"cloud.google.com/go/storage"
)

func TestObjectFromURL(t *testing.T) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		t.Fatal(err)
	}

	obj, err := ObjectFromURL(client, "gs://"+Bucket(t)+"/test_watch_object")
	if err != nil {
		t.Fatal(err)
	}

	attrs, err := obj.Attrs(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := attrs.Bucket, Bucket(t); got != want {
		t.Error("attrs.Bucket: got %q, want %q", got, want)
	}
	if got, want := attrs.Name, "test_watch_object"; got != want {
		t.Error("attrs.Name: got %q, want %q", got, want)
	}
}
