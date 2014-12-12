package boxfile


import (
  "launchpad.net/goyaml"
  "io/ioutil"
)

type Boxfile struct {
  raw []byte
  parsed map[interface{}]interface{}
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
    parsed: make(map[interface{}]interface{}),
  }
  box.parse()
  return box
}

// Node returns just a specific node from the boxfile
// if the object is a sub hash it returns a boxfile object 
// this allows Node to be chained if you know the data
func (b Boxfile) Node(name interface{}) interface{} {
  switch b.parsed[name].(type) {
  default:
      return b.parsed[name]
  case map[interface{}]interface{}:
    b := Boxfile{parsed: b.parsed[name].(map[interface{}]interface{})}
    b.fillRaw()
    b.Valid = true
    return b
  }
}

// Merge puts a new boxfile data ontop of your existing boxfile
func (self *Boxfile) Merge(box Boxfile) {
  for key, val := range box.parsed {
    self.parsed[key] = val
  }
}

// MergeProc drops a procfile into the existing boxfile
func (self *Boxfile) MergeProc(box Boxfile) {
  for key, val := range box.parsed {
    self.parsed[key] = make(map[interface{}]interface{})
    self.parsed[key]["exec"] = val
  }
}

// fillRaw is used when a boxfile is create from an existing boxfile and we want to 
// see what the raw would look like
func (b *Boxfile) fillRaw() {
  b.raw , _ = goyaml.Marshal(b.parsed)
}

// parse takes raw data and converts it to a map structure
func (b *Boxfile) parse() {
  err := goyaml.Unmarshal(b.raw, &b.parsed)
  if err != nil {
    b.Valid = false
  } else {
    b.Valid = true
  }
}

