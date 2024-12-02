package browser

import (
	"log"
	"sync"
	"testing"
)

// TestFetch
func TestFetch(t *testing.T) {

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()

		response, err := Fetch("https://www.binance.com/en/support/announcement/new-cryptocurrency-listing?c=48&navId=48&hl=en")
		if err != nil {
			t.Error(err)
		}
		log.Print("Response: ", response)
	}()

	wg.Wait()
	t.Log("Done")
}
