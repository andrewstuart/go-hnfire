package hnfire

import (
	"encoding/json"
	"log"
	"sync"
)

//GetFP returns an array of the front page to given depth. It will do so
//concurrently and may currently open too many connections. This will likely be
//fixed by pooling connections.
func GetFP(depth int) ([]*Item, error) {
	var fp []int
	res, err := cli.Get(hnBase.Child("topstories.json").String())
	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(res.Body).Decode(&fp)

	res.Body.Close()

	if err != nil {
		return nil, err
	}

	wg := &sync.WaitGroup{}
	wg.Add(len(fp))

	var fpI = make([]*Item, len(fp))

	for i := range fp {
		go func(i int) {
			defer wg.Done()
			item, err := NewItem(fp[i], depth)
			if err != nil {
				log.Println("Error getting item", i, err)
				return
			}
			fpI[i] = item
		}(i)
	}

	wg.Wait()

	return fpI, nil
}
