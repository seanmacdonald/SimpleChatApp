# SimpleChatApp
A command line chat application for 2 users (between a server and a client). 

### Setup 
Change directory to the sever folder and run the command: 
go run main.go

While the server is running, open a new shell and change the directory to the client forlder. Then run the command: 
go run main.go 

The two shells will now be able to pass messages back and forth through a websocket connection. 

When either the server or the client wants to end the connection they can simply type "bye" to do so.


