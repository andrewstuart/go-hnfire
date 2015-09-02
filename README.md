# hnfire
--
    import "github.com/andrewstuart/hnfire"


## Usage

#### func  GetFP

```go
func GetFP(depth int) ([]*Item, error)
```
GetFP returns an array of the front page to given depth. It will do so
concurrently and may currently open too many connections. This will likely be
fixed by pooling connections.

#### func  Watch

```go
func Watch(uri string, evCh chan *Event) <-chan error
```
Watch takes an endpoint and an event channel the event channel on updates to the
resource

#### type Endpoint

```go
type Endpoint string
```

An Endpoint is a convenience function for getting nested API strings.

#### func (Endpoint) Child

```go
func (ep Endpoint) Child(format string, args ...interface{}) Endpoint
```
Child returns the current endpoint concatenated with the result of formatting
the given format string with the passed arguments.

#### func (Endpoint) String

```go
func (ep Endpoint) String() string
```
String returns the Endpoint's underlying string value

#### type Event

```go
type Event struct {
	Path          string
	URI           string
	Body          io.Reader
	OriginalEvent *sse.Event
}
```

Event is a firebase-specific structure representing the path and data for an
event.

#### type Item

```go
type Item struct {
	Type        string  `json:"type"`
	Author      string  `json:"by"`
	Title       string  `json:"title"`
	Text        string  `json:"text"`
	URL         string  `json:"url"`
	Descendants int     `json:"descendants"`
	ID          int     `json:"id"`
	Points      int     `json:"score"`
	ChildrenIDs []int   `json:"kids,omitempty"`
	Children    []*Item `json:"children"`
	Rank        int
}
```

Item represents an hn story or comment, specific to the way FireBase represents
them.

#### func  NewItem

```go
func NewItem(id int, depth int) (*Item, error)
```
NewItem returns an item and all its children to `depth`. If depth is zero, no
recursion happens.

#### func (*Item) Refresh

```go
func (i *Item) Refresh() error
```
Refresh updates the item using the FireBase api.
