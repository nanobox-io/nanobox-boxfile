package boxfile


import (
  "launchpad.net/goyaml"
  "io/ioutil"
)

type Boxfile struct {
  raw []byte
  parsed map[string]interface{}
  Valid bool
} 

func NewFromPath(path string) *Boxfile {
  raw, _ := ioutil.ReadFile(path)
  return New(raw)
}

func New(raw []byte) *Boxfile {
  box := &Boxfile{
    raw: raw,
    parsed: make(map[string]interface{}),
  }
  box.Parse()
  return box
}

func (b *Boxfile) Parse() {
  err := goyaml.Unmarshal(b.raw, &b.parsed)
  if err != nil {
    b.Valid = false
  } else {
    b.Valid = true
  }
}

func (b *Boxfile) fillRaw() {
  b.raw , _ = goyaml.Marshal(b.parsed)
}

func (b *Boxfile) Node(name string) interface{} {
  switch b.parsed[name].(type) {
  default:
      return b.parsed[name]
  case map[string]interface{}:
    b := Boxfile{parsed: b.parsed[name].(map[string]interface{})}
    b.fillRaw()
    return b
  }
}