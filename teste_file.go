package main

// import (
// 	"fmt"
// 	"io"
// 	"os"
// 	"time"
// )

// // func main() {
// // 	file, err := os.OpenFile("1mbFile.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
// // 	if err != nil {
// // 		panic(err)
// // 	}

// // 	defer file.Close()

// // 	// Use bufio.Writer para bufferizar as escritas no arquivo.
// // 	w := bufio.NewWriter(file)

// // 	for i := 0; i < 1024*1024*1024; i++ {
// // 		_, err := w.WriteString(fmt.Sprintf("Miguel Lucas %d\n", i))
// // 		if err != nil {
// // 			panic(err)
// // 		}
// // 	}

// // 	// Flush para garantir que o buffer seja escrito no arquivo.
// // 	err = w.Flush()
// // 	if err != nil {
// // 		panic(err)
// // 	}

// // 	fmt.Println("File created successfully")
// // }

// func main() {

// 	file, err := os.OpenFile("1mbFile.txt", os.O_RDONLY, 0666)

// 	if err != nil {
// 		fmt.Println("Error opening file:", err)
// 		panic(err)
// 	}
// 	defer file.Close()

// 	const chunkSize = 1024
// 	errChan := make(chan error)
// 	var offset int64 = 0

// 	readChan := make(chan string, chunkSize)
// 	linesProcessed := int64(0)

// 	nProcessors :- 

// 	go func () {
// 		for i := 0; i < 10; i++ {
// 		go func(chunkIndex int) {
// 			buffer := make([]byte, chunkSize)
// 			nbytesRead, err := file.ReadAt(
// 				buffer, offset)

// 			if err != nil && err != io.EOF {
// 				fmt.Println("Error reading file:", err)
// 				panic(err)
// 			}
// 			if err == io.EOF {
// 				errChan <- err
// 			}
// 			offset += int64(nbytesRead)
// 			readChan <- string(buffer)

// 		}(i)

// 	}}

// 	select {
// 	case err := <-errChan:
// 		fmt.Println("Error reading file:", err)

// 	case read := <-readChan:
// 		fmt.Println(read)
// 		linesProcessed++
// 	case <-time.After(time.Second * 5):
// 		fmt.Printf("Processed %d lines\n", linesProcessed)
// 		return

// 	default:
// 	}

// }
