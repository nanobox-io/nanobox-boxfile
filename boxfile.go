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

func NewFromPath(path string) Boxfile {
  raw, _ := ioutil.ReadFile(path)
  return New(raw)
}

func New(raw []byte) Boxfile {
  box := Boxfile{
    raw: raw,
    parsed: make(map[interface{}]interface{}),
  }
  box.parse()
  return box
}

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

func (self *Boxfile) Merge(box Boxfile) {
  for key, val := range box.parsed {
    self.parsed[key] = val
  }
}

func (self *Boxfile) MergeProc(box Boxfile) {
  for key, val := range box.parsed {
    self.parsed[key] = make(map[interface{}]interface{})
    self.parsed[key]["exec"] = val
  }
}

func (b *Boxfile) fillRaw() {
  b.raw , _ = goyaml.Marshal(b.parsed)
}

func (b *Boxfile) parse() {
  err := goyaml.Unmarshal(b.raw, &b.parsed)
  if err != nil {
    b.Valid = false
  } else {
    b.Valid = true
  }
}

