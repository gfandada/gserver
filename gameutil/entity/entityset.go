package entity

import (
	"bytes"
)

func (es EntitySet) Add(entity *Entity) {
	es[entity] = struct{}{}
}

func (es EntitySet) Del(entity *Entity) {
	delete(es, entity)
}

func (es EntitySet) Contains(entity *Entity) bool {
	_, ok := es[entity]
	return ok
}

func (es EntitySet) String() string {
	b := bytes.Buffer{}
	b.WriteString("{")
	first := true
	for entity := range es {
		if !first {
			b.WriteString(", ")
		} else {
			first = false
		}
		b.WriteString(entity.String())
	}
	b.WriteString("}")
	return b.String()
}
