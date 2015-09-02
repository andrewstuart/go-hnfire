package hnfire

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

//Item represents an hn story or comment, specific to the way FireBase
//represents them.
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

//NewItem returns an item and all its children to `depth`. If depth is zero, no
//recursion happens.
func NewItem(id int, depth int) (*Item, error) {
	uri := fmt.Sprintf("%s/item/%d.json", hnBase, id)
	req, err := http.NewRequest("GET", uri, nil)

	if err != nil {
		return nil, err
	}

	res, err := cli.Do(req)

	if err != nil {
		return nil, err
	}

	item := Item{}

	err = json.NewDecoder(res.Body).Decode(&item)

	res.Body.Close()

	if err != nil {
		return nil, err
	}

	if depth > 0 && item.ChildrenIDs != nil {
		depth--

		item.Children = make([]*Item, len(item.ChildrenIDs))
		for i := range item.ChildrenIDs {
			var child *Item
			child, err = NewItem(item.ChildrenIDs[i], depth)
			if err != nil {
				log.Println("Error getting child", item.ChildrenIDs[i], err)
			}
			item.Children[i] = child
		}
	}

	return &item, nil
}

//Refresh updates the item using the FireBase api.
func (i *Item) Refresh() error {
	iUpdated, err := NewItem(i.ID, 0)
	if err != nil {
		return err
	}

	i.Points = iUpdated.Points
	i.Descendants = iUpdated.Descendants
	i.Title = iUpdated.Title
	i.URL = iUpdated.URL
	i.ChildrenIDs = iUpdated.ChildrenIDs
	i.Author = iUpdated.Author

	return nil
}
