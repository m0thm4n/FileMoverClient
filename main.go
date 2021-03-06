package main

import (
	"bufio"
    "bytes"
    "fmt"
    "io"
    "log"
    "net"
    "os"
    "strings"
)

func main() {
    //get port and ip address to dial

    if len(os.Args) != 3 {
        fmt.Println("useage example: tcpClient 127.0.0.1 7005")
        return
    }

    var ip string = os.Args[1]
    var port string = os.Args[2]

    connection, err := net.Dial("tcp", ip+":"+port)
    if err != nil {
        fmt.Println("There was an error making a connection")
    }

    //read from
    reader := bufio.NewReader(os.Stdin)
    fmt.Print("Please enter 'get <filename>' or 'send <filename>' to transfer files to the server\n\n")
    inputFromUser, _ := reader.ReadString('\n')
    arrayOfCommands := strings.Split(inputFromUser, " ")

    if arrayOfCommands[0] == "get" {
        getFileFromServer(arrayOfCommands[1], connection)

    } else if arrayOfCommands[0] == "send" {
        sendFileToServer(arrayOfCommands[1], connection)
    } else {
        fmt.Println("Bad Command")
    }
}

const BUFFER_SIZE = 1024

func sendFileToServer(fileName string, connection net.Conn) {

    var currentByte int64 = 0
    fmt.Println("send to client")
    fileBuffer := make([]byte, BUFFER_SIZE)

    var err error

    //file to read
    file, err := os.Open(strings.TrimSpace(fileName)) // For read access.
    if err != nil {
        connection.Write([]byte("-1"))
        log.Fatal(err)
    }
    connection.Write([]byte(fileName))
    //read file until there is an error
    for err == nil || err != io.EOF {

        _, err = file.ReadAt(fileBuffer, currentByte)
        currentByte += BUFFER_SIZE
        fmt.Println(fileBuffer)
        connection.Write(fileBuffer)
    }

    file.Close()
    connection.Close()

}

func getFileFromServer(fileName string, connection net.Conn) {

    var currentByte int64 = 0

    fileBuffer := make([]byte, BUFFER_SIZE)

    var err error
    file, err := os.Create(strings.TrimSpace(fileName))
    if err != nil {
        log.Fatal(err)
    }
    connection.Write([]byte("get " + fileName))
    for {

        connection.Read(fileBuffer)
        cleanedFileBuffer := bytes.Trim(fileBuffer, "\x00")

        _, err = file.WriteAt(cleanedFileBuffer, currentByte)

        currentByte += BUFFER_SIZE

        if err == io.EOF {
            break
        }

    }

    file.Close()
    return

}