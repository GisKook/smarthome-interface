package main

import "log"

func main() {
  /*  messages := make(chan string)

    go func(){messages <-"ping"}()

    msg := <-messages
    log.Println(msg)*/
    var cs map[string]chan byte = make(map[string]chan byte)

        chanKey := "a"

            ci, ok := cs[chanKey]
                if !ok {
                            ci = make(chan byte)
                                    cs[chanKey] = ci
                                        }
                    
    go func(){ci <-'b'}()
                                                log.Println("make chan!")

                                                    select {
                                                            case <-ci:
                                                                    jsonRes := 1
                                                                            log.Println(jsonRes)
                                                                                }

                                                                                    cb, ok := cs[chanKey]
                                                                                        if !ok {
                                                                                                    log.Println("no chan!")
                                                                                                        } else {
                                                                                                                    ci <- 'b'
                                                                                                                            close(cb)
                                                                                                                                    delete(cs, chanKey)
                                                                                                                                            log.Println("delete existed chan!")
                                                                                                                                                }
                                                                                                                                                    log.Println("make chan!")
}

