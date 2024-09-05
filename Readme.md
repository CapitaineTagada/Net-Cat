# Net-Cat
[![Made with Go](https://img.shields.io/badge/Go-1-blue?logo=go&logoColor=white)](https://golang.org "Go to Go homepage")  
Net-Cat consists on recreating the NetCat in a Server-Client Architecture that can run in a server mode on a specified port listening for incoming connections, and it can be used in client mode, trying to connect to a specified port and transmitting information to the server.

## Installation:
This project need building before you start it in order to run properly. Use this command in your terminal:
```bash
go build -o net-cat
```

## Starting server:
Now you can start your server on terminal using this command:
```bash
./net-cat #Starting your server in localhost
```

To start the client setup, you need to write this command in a new terminal:
```bash
nc localhost + port
```

## Tech:
go version go 1.22.5

Only standard packages are used

## Authors:
A. Joly

S. Cointin
