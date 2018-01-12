sima
---
<a href="https://godoc.org/github.com/theodesp/sima">
<img src="https://godoc.org/github.com/theodesp/sima?status.svg" alt="GoDoc">
</a>

<a href="https://opensource.org/licenses/MIT" rel="nofollow">
<img src="https://img.shields.io/github/license/mashape/apistatus.svg" alt="License"/>
</a>

<a href="https://travis-ci.org/theodesp/sima" rel="nofollow">
<img src="https://travis-ci.org/theodesp/sima.svg?branch=master" />
</a>

<a href="https://codecov.io/gh/theodesp/sima">
  <img src="https://codecov.io/gh/theodesp/sima/branch/master/graph/badge.svg" />
</a>

Sima is a simple object to object or broadcast dispatching system. 
Any number of interested parties can subscribe to events. 
Signal receives can receive also signals from specific senders.

## Installation
```bash
$ go get -u github.com/theodesp/sima
```

## Usage
1. Create a topic factory and re-use it for all signals yo want to create:

```go
tf := NewTopicFactory()
onStart := NewSima(tf)
onEnd := NewSima(tf)
```

2. After you have created your signals, just connect handlers for a particular topic or not. If you don't specify a topic then the handler will be assigned a ALL topic that defaults as a broadcast address.
```go
// Subscribe to ALL
onStart.Connect(func(context context.Context, sender *Topic) interface{} {
		fmt.PrintF("OnStart called from Sender %+v", sender)
    return sender
	}, nil)

// Subscribe to specific sender/topic
onEnd.Connect(func(context context.Context, sender *Topic) interface{} {
		fmt.PrintF("onEnd called from Sender %+v", sender)
    return sender
}, "on-end-sender")
```

3. Now just send some messages and any registered participant will call the handler.
```go
response := onStart.Dispatch(context.Background(), nil) // will handle
response := onStart.Dispatch(context.Background(), "on-start-sender") // will not handle

response := onEnd.Dispatch(context.Background(), nil) // will not handle
response := onEnd.Dispatch(context.Background(), "on-end-sender") // will handle
```

## LICENCE
Copyright © 2017 Theo Despoudis MIT license
