**Simple Pub/Sub Using Websockets**

This is a simple implementation of the pubsub pattern using the gorilla websockets library. 

Two endpoints are exposed:

`/subscribe` - Hit this endpoint to subscribe to messages

`/publish` - Send messages to this endpoint to publish them to all subscribers. 

Messages should look like this: 

`{
    "message": "hello"
}`

Run the server like so

`go run .`

Run the tests like so

`go test .`


