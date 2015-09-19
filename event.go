package hnfire

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/andrewstuart/go-sse"
)

type fireEvent struct {
	Path string          `json:"path"`
	Data json.RawMessage `json:"data"`
}

//Event is a firebase-specific structure representing the path and data for an event.
type Event struct {
	Path          string
	URI           string
	Body          io.Reader
	OriginalEvent *sse.Event
	Error         error
}

//Watch takes an endpoint and an event channel the event channel on updates to
//the resource
func Watch(uri string, evCh chan *Event) <-chan error {
	errCh := make(chan error)
	if evCh == nil {
		close(errCh)
		return errCh
	}

	ch := make(chan *sse.Event)
	go sse.Notify(hnBase.Child(uri).String(), ch)

	for evt := range ch {
		if evt.Type == "put" {
			fireEvt := fireEvent{}
			err := json.NewDecoder(evt.Data).Decode(&fireEvt)

			if err != nil {
				//TODO better error handling
				select {
				case errCh <- err:
				default:
				}
				continue
			}

			evCh <- &Event{
				Path:          fireEvt.Path,
				URI:           uri,
				Body:          bytes.NewBuffer(fireEvt.Data),
				OriginalEvent: evt,
			}
		}
	}

	return errCh
}
