package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	"time"
)

const LogoAscii = "\x1b[32m" + `
         _nnnn_
        dGGGGMMb
       @p~qp~~qMb
       M|@||@) M|
       @,----.JM|
      JS^\__/  qKL
     dZP        qKRb
    dZP          qKKb
   fZP            SMMb
   HZM            MMMM
   FqM            MMMM
 __| ".        |\dS"qML
 |    '.       | '' \Zq
_)      \.___.,|     .'
\____   )MMMMMP|   .'
     '-'       '--'
` + "\x1b[0m"

var users = make(map[net.Conn]string)
var mu sync.Mutex
var messageHistory []string

func main() {
	port := "8080"
	host := "localhost"
	// 1. Start the server on port 8080
	if len(os.Args) > 1 {
		port = os.Args[2]
		host = os.Args[1]
	}
	listener, err := net.Listen("tcp4", host+":"+port)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server starting on port :" + port)
	fmt.Printf("Welcome to the Lord of Pigeon Tchat :\n %s", LogoAscii)
	// 2. Manege incoming connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		// Check if the maximum connections is reached
		if len(users) >= 3 {
			conn.Write([]byte("Server is full\n"))
			conn.Close()
			continue
		}
		// 3. Manage the every connection in a goroutine
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	// 4. logic to handle the connection
	fmt.Printf("New connection from: %s"+"\n", conn.RemoteAddr().String())
	defer conn.Close()
	CreateUsername(conn)
	sendHistory(conn)
	broadcastMessage("[" + time.Now().Format(time.DateTime) + "] " + users[conn] + " has joined the chat\n")
	logMessage := "[" + time.Now().Format(time.DateTime) + "] New connection from: " + conn.RemoteAddr().String() + " " + users[conn] + " " + "\n"
	if _, err := os.Stat("./Logs/log" + time.Now().Format("02012006") + ".txt"); os.IsNotExist(err) {
		// If the file do not exist, create it
		_, err := os.Create("./Logs/log" + time.Now().Format("02012006") + ".txt")
		if err != nil {
			fmt.Println("Error creating file:", err)
			return
		}
	}
	logToFile(logMessage) // Log the connection
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				fmt.Println(users[conn] + " disconnected succesfully ")
				break
			} else {
				fmt.Println("Error reading:", err)
				break
			}
		}
		if string(buf[:n]) == "/rename\n" {
			Lastname := users[conn]
			conn.Write([]byte("Enter your new username:"))
			scanner := bufio.NewScanner(conn)
			// 6. Read the username from the client
			if scanner.Scan() {
				users[conn] = scanner.Text()
				broadcastMessage(fmt.Sprintf("[%s] [%s] changed his name to : [%s]\n", time.Now().Format(time.DateTime), Lastname, users[conn]))
			}
			continue
		}
		if string(buf[:n]) != "\n" {
			message := fmt.Sprintf("[%s] %s: %s", time.Now().Format(time.DateTime), users[conn], string(buf[:n]))
			fmt.Print(message)
			messageHistory = append(messageHistory, message)
			broadcastMessage(message)
		}
	}

	file, err := os.OpenFile("./Logs/log"+time.Now().Format("02012006")+".txt", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()
	defer func() {
		// 9. Handle the disconnection
		broadcastMessage("[" + time.Now().Format(time.DateTime) + "] " + users[conn] + " has left the chat\n")
		logDisconnection(conn) // Log the disconnection
		delete(users, conn)    // Remove the user from the map
		conn.Close()           // Close the connection
	}()
}

func CreateUsername(conn net.Conn) {
	// 5. Ask the user for a username
	var username string
	conn.Write([]byte("Enter your username:"))
	scanner := bufio.NewScanner(conn)
	// 6. Read the username from the client
	if scanner.Scan() {
		username = scanner.Text()
		users[conn] = username
		fmt.Printf("[%s] New user connected: [%s]\n", time.Now().Format(time.DateTime), username)
	}
}

func sendHistory(conn net.Conn) {
	for _, message := range messageHistory {
		conn.Write([]byte(message))
	}
}
func broadcastMessage(message string) {
	// 7. Broadcast the message to all the connected clients
	mu.Lock()
	for conn := range users {
		_, err := conn.Write([]byte(message))
		if err != nil {
			fmt.Println("Error writing to connection:", err)
			delete(users, conn)
		}
	}
	logToFile(message)
	mu.Unlock()
}
func logToFile(message string) {
	// 8. Log the message to a file
	file, err := os.OpenFile("./Logs/log"+time.Now().Format("02012006")+".txt", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()
	if _, err := file.WriteString(message); err != nil {
		fmt.Println("Error writing to file:", err)
	}
}

func logDisconnection(conn net.Conn) {
	// Create the log message for disconnection
	logMessage := "[" + time.Now().Format(time.DateTime) + "] from [" + conn.RemoteAddr().String() + "] User disconnected: " + users[conn] + "\n"
	logToFile(logMessage) // Log the disconnection
}
