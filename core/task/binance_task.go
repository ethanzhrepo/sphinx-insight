package task

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/ethanzhrepo/sphinx-insight/core/browser"
	"github.com/ethanzhrepo/sphinx-insight/core/db"
	"github.com/ethanzhrepo/sphinx-insight/core/notifier"
	"golang.org/x/net/html"
)

// const binance_announcement
const BinanceAnnouncement = "binance_announcement"

// BinanceTask
type BinanceTask struct {
	context context.Context
	cancel  context.CancelFunc
	ps      *notifier.SimplePubSub
	db      *db.LevelDB
}

// NewBinanceTask
func NewBinanceTask(
	ps *notifier.SimplePubSub,
	db *db.LevelDB,
) *BinanceTask {
	context, cancel := chromedp.NewContext(context.Background())
	return &BinanceTask{
		context: context,
		cancel:  cancel,
		ps:      ps,
		db:      db,
	}
}

// Close
func (t *BinanceTask) Close() {
	t.cancel()
}

// fetchAnnouncement
func (t *BinanceTask) fetchAnnouncement() {
	response, e := browser.FetchDom(t.context,
		"https://www.binance.com/en/support/announcement/new-cryptocurrency-listing?c=48&navId=48&hl=en", "section>div>div:first-of-type")

	if e != nil {
		log.Fatal(e)
		return
	}

	// parse response
	doc, err := html.Parse(strings.NewReader(response))
	if err != nil {
		log.Fatal(err)
		return
	}

	// find announcement
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			// herf
			var herf string
			if len(n.Attr) > 1 && n.Attr[1].Key == "href" {
				herf = n.Attr[1].Val
			}

			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if c.Type == html.ElementNode && c.Data == "div" {
					hash := md5.Sum([]byte(c.FirstChild.Data))
					hash2hex := hex.EncodeToString(hash[:])

					// check if hash2hex exists in level db, if exists, skip
					_, err := t.db.Get([]byte(hash2hex))
					if err == nil {
						// exists
						return
					}

					log.Println("[Binance Announcement] ", c.FirstChild.Data, "Hash:", hash2hex)
					// if not exists, store to level db and publish to subscribers
					t.db.Put([]byte(hash2hex), []byte(c.FirstChild.Data))
					// publish to subscribers
					// a json object, title and herf
					obj := map[string]string{
						"content": c.FirstChild.Data,
						"link":    herf,
					}

					if data, err := json.Marshal(obj); err == nil {
						t.ps.Publish(BinanceAnnouncement, string(data))
					} else {
						log.Println("Error marshaling JSON:", err)
					}

				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(doc)
}

// Run
func (t *BinanceTask) Run() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		go t.fetchAnnouncement()
	}
}
