package boxfile

import "testing"
import "strings"
import "encoding/json"
import "fmt"


func Testparse(t *testing.T) {
  box := New([]byte("a: Easy!\nb:\n  c: 2\n  d: [3, 4]\n"))
  if box.StringValue("a") != "Easy!" {
    t.Error("boxfile parsed does not match boxfile in")
  }

  bad := New([]byte("baz:\n   cdr\nfoo::*)-> bar"))
  if bad.Valid {
    t.Error("Boxfile thinks its valid but it shoudnt be")
  }

}


func TestNode(t *testing.T) {
  box := New([]byte("web1:\n  name: site\n  type: php\n  version: 5.4\n  php_extensions:\n    - mysql\n    - gd\n    - eaccelerator\n"))
  web1 := box.Node("web1")
  if web1.StringValue("name") != "site" {
    t.Error("nested nodes dont work")
  }
  if string(web1.raw) != "name: site\nphp_extensions:\n- mysql\n- gd\n- eaccelerator\ntype: php\nversion: 5.4\n" {
    t.Error("subnodes dont create raw yaml correctly")
  }
}

func TestParsedSubParts(t *testing.T) {
  box := New([]byte("a: Easy!\nb:\n  c: 2\n  d: [3, 4]\n"))
  invalidNode := box.Node("nonya")
  if invalidNode.Parsed == nil {
    t.Error("the parsed data in a invalid node should be an empty map")
  }
}

func TestDeepNesting(t *testing.T) {
  box := New([]byte(`run.config:
  engine: php
  engine.config:
    extensions:
      - mysqli
data.mysql:
  image: nanobox/mysql
  config:
    users:
    - username: root
      meta:
        privileges:
        - privilege: ALL PRIVILEGES
          'on': "*.*"
          with_grant: true
    - username: nanobox
      meta:
        privileges:
        - privilege: ALL PRIVILEGES
          'on': gonano.*
          with_grant: true
        - privilege: ALL PRIVILEGES
          'on': testing.*
          with_grant: true
        - privilege: ALL PRIVILEGES
          'on': blah.*
          with_grant: true
        - privilege: PROCESS
          'on': "*.*"
          with_grant: false
        - privilege: SUPER
          'on': "*.*"
          with_grant: false
        databases:
        - gonano
        - testing
        - blah`))

  b, err := json.Marshal(box.Node("data.mysql").Node("config").Parsed)
  if err != nil {
    t.Error("unable to marshal nested interfaces: %s", err)
  }
  fmt.Print(string(b))
  if !strings.Contains(string(b), "blah.*") {
    t.Error("does not contain data from deeply nested interfaces")
  }
}