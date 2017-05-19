package main

import "fmt"

func main(){
    message:= make(chan string,2)
    message<-"buffered"
  //  message<-"channel"
 //   message<-"2222"
    fmt.Println(<-message)
    fmt.Println(<-message)
}
