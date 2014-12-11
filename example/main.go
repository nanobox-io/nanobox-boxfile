package main

import "github.com/nanobox-core/boxfile"
import "fmt"

func main() {
  box := boxfile.New([]byte("a: Easy!\nb:\n  c: 2\n  d: [3, 4]\n"))
  fmt.Println(box.Node("a").(string))

  box2 := boxfile.NewFromPath("example/Boxfile")
  fmt.Println(box2.Node("web1").(boxfile.Boxfile).Node("php_extensions"))
  box2.Merge(box)
  // fmt.Printf("%+v\n", box2)
}