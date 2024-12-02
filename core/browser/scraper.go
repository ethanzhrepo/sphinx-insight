package browser

import (
	"context"
	"log"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

func Fetch(url string) (string, error) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	err := chromedp.Run(ctx, network.Enable())
	if err != nil {
		log.Fatal(err)
	}

	var res string
	err = chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Text("body", &res),
	)
	if err != nil {
		return "", err
	}

	log.Println("result: ", res)
	return res, nil
}

func FetchDom(ctx context.Context, url string, selector string) (string, error) {
	var data string

	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.OuterHTML(selector, &data, chromedp.ByQuery),
	)

	if err != nil {
		log.Fatal(err)
		return "", err
	}

	return data, nil
}
