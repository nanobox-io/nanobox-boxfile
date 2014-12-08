package boxfile

import "testing"
// import "fmt"


func TestParse(t *testing.T) {
  box := New([]byte("a: Easy!\nb:\n  c: 2\n  d: [3, 4]\n"))
  if box.Node("a").(string) != "Easy!" {
    t.Error("boxfile parsed does not match boxfile in")
  }

  bad := New([]byte("baz:\n   cdr\nfoo::*)-> bar"))
  if bad.Valid {
    t.Error("Boxfile thinks its valid but it shoudnt be")
  }

}
