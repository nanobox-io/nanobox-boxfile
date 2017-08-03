package boxfile_test

import (
	"io/ioutil"
	"os"
	// "strings"
	"testing"

	"github.com/nanobox-io/nanobox-boxfile"
)

// TestNew tests parsing raw yaml bytes
func TestNew(t *testing.T) {
	box := boxfile.New([]byte(testBoxfile))
	t.Logf("%q\n", box.Parsed)
	if box.Parsed["run.config"].(map[string]interface{})["engine"] != "none" {
		t.Errorf("Failed to return parsed boxfile - %v", box.Parsed["run.config"])
		t.FailNow()
	}

	box = boxfile.New([]byte(badBoxfile))
	t.Logf("%t\n", box.Valid)
	if box.Valid {
		t.Errorf("Failed to fail parsing boxfile - %v", box.Parsed)
		t.FailNow()
	}
}

// TestNewFromPath tests fetching a boxfile from a path
func TestNewFromPath(t *testing.T) {
	// remove /tmp/boxfile.yml
	err := os.RemoveAll("/tmp/boxfile.yml")
	if err != nil {
		t.Errorf("Failed to remove /tmp/boxfile.yml - %s", err.Error())
		t.FailNow()
	}

	box := boxfile.NewFromPath("/tmp/boxfile.yml")
	if box.Valid {
		t.Errorf("Should have failed to find boxfile.yml - %v", box)
	}

	// write invalid boxfile to /tmp/boxfile.yml
	err = ioutil.WriteFile("/tmp/boxfile.yml", []byte(testBoxfile), 0644)
	if err != nil {
		t.Errorf("Failed to write invalid yaml boxfile - %s", err.Error())
		t.FailNow()
	}

	box = boxfile.NewFromPath("/tmp/boxfile.yml")
	if !box.Valid {
		t.Errorf("Should have found boxfile.yml - %v", box)
	}

	// remove /tmp/boxfile.yml
	err = os.RemoveAll("/tmp/boxfile.yml")
	if err != nil {
		t.Errorf("Failed to remove /tmp/boxfile.yml - %s", err.Error())
		t.FailNow()
	}
}

// TestNewFromFile ensures we can differentiate between a missing boxfile.yml(err)
// or an invalid one(!.Valid) when fetching a boxfile from a path
func TestNewFromFile(t *testing.T) {
	// remove /tmp/boxfile.yml
	err := os.RemoveAll("/tmp/boxfile.yml")
	if err != nil {
		t.Errorf("Failed to remove /tmp/boxfile.yml - %s", err.Error())
		t.FailNow()
	}

	box, err := boxfile.NewFromFile("/tmp/boxfile.yml")
	if box != nil || err == nil {
		t.Error("Should have failed to find boxfile.yml")
		t.FailNow()
	}
	t.Logf("Error - %s", err.Error())

	// write invalid boxfile to /tmp/boxfile.yml
	err = ioutil.WriteFile("/tmp/boxfile.yml", []byte(badBoxfile), 0644)
	if err != nil {
		t.Errorf("Failed to write invalid yaml boxfile - %s", err.Error())
		t.FailNow()
	}

	box, err = boxfile.NewFromFile("/tmp/boxfile.yml")
	if err != nil {
		t.Error("Should have found a boxfile.yml")
		t.FailNow()
	}
	if box.Valid {
		t.Errorf("Should have failed validity check - %q", box.Parsed)
		t.FailNow()
	}
	t.Logf("Boxfile - %s", box.Raw)

	// remove /tmp/boxfile.yml
	err = os.RemoveAll("/tmp/boxfile.yml")
	if err != nil {
		t.Errorf("Failed to remove /tmp/boxfile.yml - %s", err.Error())
		t.FailNow()
	}
}

// TestNode ensures we can get boxfile sub-hashes
func TestNode(t *testing.T) {
	box := boxfile.New([]byte(testBoxfile2))
	webApi := box.Node("web.api")
	t.Log(webApi.Parsed)
	startStrings := ""
	for _, v := range webApi.Parsed["start"].([]interface{}) {
		if start, ok := v.(string); ok {
			startStrings += start
		}
	}
	if startStrings != "run thinglog thing" {
		t.Errorf("Failed to return parsed boxfile - %q", startStrings)
		t.FailNow()
	}

	// start := box.Node("web.api").Node("start")
	// t.Log(start.Parsed)

	start := box.Node("web.db").Node("start")
	t.Log(start.Parsed)
}

func TestNodes(t *testing.T) {
	box := boxfile.New([]byte(testBoxfile2))
	nodes := box.Nodes()
	t.Log(nodes)

	nodes = box.Nodes("web")
	t.Log(nodes)
	nodes = box.Nodes("container")
	t.Log(nodes)
}

// func Testparse(t *testing.T) {
//   box := boxfile.New([]byte("a: Easy!\nb:\n  c: 2\n  d: [3, 4]\n"))
//   if box.StringValue("a") != "Easy!" {
//     t.Error("boxfile parsed does not match boxfile in")
//   }

//   bad := boxfile.New([]byte("baz:\n   cdr\nfoo::*)-> bar"))
//   if bad.Valid {
//     t.Error("Boxfile thinks its valid but it shoudnt be")
//   }
// }

// func TestNode(t *testing.T) {
//   box := boxfile.New([]byte("web1:\n  name: site\n  type: php\n  version: 5.4\n  php_extensions:\n    - mysql\n    - gd\n    - eaccelerator\n"))
//   web1 := box.Node("web1")
//   if web1.StringValue("name") != "site" {
//     t.Error("nested nodes dont work")
//   }
//   // if string(web1.raw) != "name: site\nphp_extensions:\n- mysql\n- gd\n- eaccelerator\ntype: php\nversion: 5.4\n" {
//   //   t.Error("subnodes dont create raw yaml correctly")
//   // }
// }

// func TestParsedSubParts(t *testing.T) {
//   box := boxfile.New([]byte("a: Easy!\nb:\n  c: 2\n  d: [3, 4]\n"))
//   invalidNode := box.Node("nonya")
//   if invalidNode.Parsed == nil {
//     t.Error("the parsed data in a invalid node should be an empty map")
//   }
// }

var (
	testBoxfile string = `
run.config:
  engine: none

data.db:
  image: postgresql:9.6
`

	testBoxfile2 string = `
run.config:
  engine: none

web.api:
  start:
    - run thing
    - log thing

web.db:
  start:
    db: run thing
    log: log thing
`

	badBoxfile string = `
run.config:
    engine: none

data.db;
  image: postgresql:9.6
`
)
