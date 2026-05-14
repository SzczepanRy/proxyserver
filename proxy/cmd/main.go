package main

import (
	"fmt"
	"log"
	"net"
)



func main(){

	l , err := net.Listen("tcp" , ":8080")

	if err != nil {
		log.Fatalf("could not init Listener , %v" , err)
	}

	defer l.Close()

	for {
		conn , err:= l.Accept()

		if err != nil {
			log.Fatalf("could not init Listener , %v" , err)
		}






	}




	fmt.Println("proxy started")

}
