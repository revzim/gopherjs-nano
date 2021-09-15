# gopherjs-nano | nanojs

## simple gopherjs wrapper for the [nano golang game framework](https://github.com/lonng/nano) [client sdk](https://github.com/nano-ecosystem/nano-websocket-client)

## IMPORTANT
### DO NOT INCLUDE NANO (STARX) CLIENT SIDE, THIS LIBRARY WILL AUTO INJECT THE CLIENT LIBRARY
### GOPHERJS EXAMPLE OF A PRETTY MUCH 1:1 OF THE ORIGINAL CHAT EXAMPLE
![EXAMPLE](./example.PNG)
[ORIGINAL NANO CHAT EXAMPLE](https://github.com/lonng/nano/tree/master/examples/demo/chat)
# [/example](https://github.com/revzim/gopherjs-nano/example)
* clone project `git clone https://github.com/revzim/gopherjs-nano`
* `go mod tidy`
* `cd example`
	* EDIT & REBUILD GOPHERJS NANO CLIENT EXAMPLE: `buildjs.sh`
* RUN WEBSERVER & GAMESERVER: `go run main.go`
* OPEN UP http://localhost:8080 TO TEST EXAMPLE