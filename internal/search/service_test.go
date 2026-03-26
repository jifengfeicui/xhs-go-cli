package search

import (
	"testing"
)

func TestParseSearchResults(t *testing.T) {
	raw := []byte(`{
	  "data": {
	    "data": {
	      "items": [
	        {"id": "feed1", "xsec_token": "token1", "title": "标题1", "author": "作者1"},
	        {"id": "feed2", "xsec_token": "token2", "title": "标题2", "author": "作者2"}
	      ]
	    }
	  }
	}`)
	items, err := parseSearchResults(raw)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("got %d items", len(items))
	}
	if items[0].FeedID != "feed1" || items[1].Title != "标题2" {
		t.Fatalf("unexpected items: %#v", items)
	}
}
