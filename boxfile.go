package boxfile


import (
  "launchpad.net/goyaml"
  "io/ioutil"
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
func (b *Boxfile) Node(name interface{}) interface{} {
  switch b.Parsed[name].(type) {
  default:
      return b.Parsed[name]
  case map[interface{}]interface{}:
    b := Boxfile{Parsed: b.Parsed[name].(map[interface{}]interface{})}
    b.fillRaw()
    b.Valid = true
    return b
  }
}

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

