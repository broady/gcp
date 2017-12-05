package gcp

import (
	"context"
	"io"
	"os"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"

	"cloud.google.com/go/storage"
)

var ctx = context.Background()

func TestWatchObject(t *testing.T) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		t.Fatal(err)
	}
	obj := client.Bucket(Bucket(t)).Object("test_watch_object")

	ch := WatchObject(ctx, obj, nil)
	update1 := <-ch
	if update1.Err != nil {
		t.Fatal(update1.Err)
	}

	go func() {
		w := obj.NewWriter(ctx)
		io.WriteString(w, uuid.NewV4().String())
		w.Close()
	}()

	select {
	case <-time.After(10 * time.Second):
		t.Fatal("timed out waiting for change")
	case update2 := <-ch:
		if update2.Err != nil {
			t.Fatal(update2.Err)
		}
		if update1.Generation == update2.Generation {
			t.Fatal("want different generation; got same")
		}
	}
}

func ProjectID(t *testing.T) string {
	proj := os.Getenv("GCP_TEST_PROJECT")
	if proj == "" {
		t.Skip("GCP_TEST_PROJECT not set")
	}
	return proj
}

func Bucket(t *testing.T) string {
	bkt := os.Getenv("GCP_TEST_BUCKET")
	if bkt == "" {
		t.Skip("GCP_TEST_BUCKET not set")
	}
	return bkt
}
