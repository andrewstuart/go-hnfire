package hnfire

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
)

const fakeItem = `{
	"type": "story",
	"id": 1,
	"url": "http://google.com",
	"title": "Google!",
	"descendants": 3,
	"kids": [2],
	"score": 999999,
	"author": "andrewstuart2"
}`

const fakeItem2 = `{
	"type": "story",
	"id": 1,
	"url": "http://google.com",
	"title": "Google!",
	"descendants": 3,
	"kids": [2],
	"score": 1000000,
	"author": "andrewstuart2"
}`

func TestFpGetter(t *testing.T) {
	var requests int32
	var update bool
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		switch {
		case r.URL.Path == "/topstories.json":
			fmt.Fprintln(w, `[1, 2, 3, 4, 5]`)
		case strings.HasPrefix(r.URL.Path, "/item/"):
			atomic.AddInt32(&requests, 1)
			if update {
				fmt.Fprintln(w, fakeItem2)
			} else {
				fmt.Fprintln(w, fakeItem)
			}
		}
	}))
	defer ts.Close()

	hnBase = Endpoint(ts.URL)

	fp, err := GetFP(1)

	if err != nil {
		t.Fatalf("Error getting message: %v", err)
	}

	if len(fp) != 5 {
		t.Fatalf("Wrong length of response: %v", fp)
	}

	if fp[0].ID != 1 {
		t.Fatalf("Wrong last element: %v, should have ID 1", fp[0])
	}

	if len(fp[0].ChildrenIDs) != 1 {
		t.Fatalf("Wrong length of children: %v, should have ChildrenIDs length 1", fp[0].ChildrenIDs)
	}

	if requests != 10 {
		t.Errorf("Wrong number of requests: %d", requests)
	}

	if fp[0].Children[0].ID != 1 {
		t.Fatalf("Did not retrieve past first level.")
	}

	update = true

	err = fp[0].Refresh()
	if err != nil {
		t.Fatalf("Error refreshing item: %v", err)
	}

	if fp[0].Points != 1000000 {
		t.Errorf("Did not update points after refresh.")
	}
}
