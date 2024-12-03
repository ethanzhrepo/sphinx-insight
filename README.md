# sphinx-insight

This is a test project, just for test.

Monitor the content of Binance's announcement channel, updated every 5 seconds
If new content is found, it will be pushed to the specified conversation in telegram

Perform the following operations in the code root directory

```shell
touch .env
echo "TELEGRAM_BOT_TOKEN=yourtoken" >> .env
echo "TELEGRAM_CHAT_ID=yourchatid" >> .env
echo OPENAI_API_KEY=your-openai-key >> .env
echo DEBUG=true >> .env
```

Then run go run main.go for testing

**The code is still under development, In fact, it just started...**

## Development plan

1. Extract the relevant token name from the text, connect to the api to query the real-time price, and provide a transaction link.

2. Obtain token-related public opinion data from elastic search or other retrieval interfaces, combine the search index, and try to give sentiment analysis.

3. Actions: Mainly dex and binance api fast transaction