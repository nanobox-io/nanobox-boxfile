package main

import "github.com/nanobox-core/boxfile"
import "fmt"

func main() {
  box := boxfile.New([]byte("a: Easy!\nb:\n  c: 2\n  d: [3, 4]\n"))

  fmt.Println(box.Node("a").(string))
  
}