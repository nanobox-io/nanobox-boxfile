package boxfile


import (
  "launchpad.net/goyaml"
  "io/ioutil"
  "strconv"
)

type Boxfile struct {
  raw []byte
  Parsed map[interface{}]interface{}
  Valid bool
}

// NewFromPath creates a new boxfile from a file instead of raw bytes
func NewFromPath(path string) Boxfile {
  raw, _ := ioutil.ReadFile(path)
  return New(raw)
}

// New returns a boxfile object from raw data
func New(raw []byte) Boxfile {
  box := Boxfile{
    raw: raw,
    Parsed: make(map[interface{}]interface{}),
  }
  box.parse()
  return box
}

// Node returns just a specific node from the boxfile
// if the object is a sub hash it returns a boxfile object 
// this allows Node to be chained if you know the data
func (self *Boxfile) Node(name interface{}) (box Boxfile) {
  switch self.Parsed[name].(type) {
  case map[interface{}]interface{}:
    box.Parsed = self.Parsed[name].(map[interface{}]interface{})
    box.fillRaw()
    box.Valid = true
  default:
    box.Valid = false
  }
  return
}

func (b *Boxfile) Value(name interface{}) interface{} {
  return b.Parsed[name]
}

func (b *Boxfile) StringValue(name interface{}) string {
  switch b.Parsed[name].(type) {
  default:
    return ""
  case string:
    return b.Parsed[name].(string)
  case bool:
    return strconv.FormatBool(b.Parsed[name].(bool))
  case int:
    return strconv.Itoa(b.Parsed[name].(int))
  }
}

func (b *Boxfile) IntValue(name interface{}) int {
  switch b.Parsed[name].(type) {
  default:
    return 0
  case string:
    i, _ := strconv.Atoi(b.Parsed[name].(string))
    return i
  case bool:
    if b.Parsed[name].(bool) == true {
      return 1
    }
    return 0
  case int:
    return b.Parsed[name].(int)
  }
}

func (b *Boxfile) BoolValue(name interface{}) bool {
  switch b.Parsed[name].(type) {
  default:
    return false
  case string:
    boo, _ :=strconv.ParseBool(b.Parsed[name].(string))
    return boo
  case bool:
    return b.Parsed[name].(bool)
  case int:
    return (b.Parsed[name].(int) != 0)
  }
}

// list nodes
func (b *Boxfile) Nodes() (rtn []interface{}) {
  for key, _ := range b.Parsed {
    rtn = append(rtn, key)
  }
  return
}

// Merge puts a new boxfile data ontop of your existing boxfile
func (self *Boxfile) Merge(box Boxfile) {
  for key, val := range box.Parsed {
    self.Parsed[key] = val
  }
}

// MergeProc drops a procfile into the existing boxfile
func (self *Boxfile) MergeProc(box Boxfile) {
  for key, val := range box.Parsed {
    self.Parsed[key] = map[interface{}]interface{}{"exec":val}
  }
}

// fillRaw is used when a boxfile is create from an existing boxfile and we want to 
// see what the raw would look like
func (b *Boxfile) fillRaw() {
  b.raw , _ = goyaml.Marshal(b.Parsed)
}

// parse takes raw data and converts it to a map structure
func (b *Boxfile) parse() {
  err := goyaml.Unmarshal(b.raw, &b.Parsed)
  if err != nil {
    b.Valid = false
  } else {
    b.Valid = true
  }
}

