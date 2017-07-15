# SOFY'S

SOFY'S is a tool that streams chat messages for youtube live by JSON format.

SOFY'S stands for School Of Fish in Youtube Streaming. That's why I thought streaming the messages is similar to school of fish in sea.

# Build

Create a file which have been defined `client_id` and `cliend_secret` variable like below:
```
package main

const (
  client_id = "<YOUR CLIENT ID>"
  client_secret = "<YOUR CLIENT SECRET>"
)
```
Put it on a same folder as sofys

```
$ go get github.com/bamchoh/sofys
$ go build
```

# Usage
When you launch sofys, you have to enter a token which is display in your browser first.

Then, you have to enter a video id you want to stream.

For video id, you can get it in the URL you are watching.

For example, Let's say there is a URL like below:

`https://www.youtube.com/watch?v=fok8A9mdQbM`

A part of video id in above URL is `fok8A9mdQbM`

You copy it from URL in your browser, and paste it in the command line.

Finally, sofys starts to stream chat message.
