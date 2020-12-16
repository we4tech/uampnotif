package configs

import "log"

type Params []Param

//
// Param represents a specific parameter with name, label and default value.
//
type Param struct {
	Name         string
	Label        string
	DefaultValue string
}

//
// ForEach implements Iterator method for providing access to internal data
// structure.
//
func (p *Params) ForEach(cb func(i int, p *Param)) {
	if p.IsEmpty() {
		log.Printf("IterableParams.ForEach: Empty Params set")

		return
	}

	for i, param := range *p {
		cb(i, &param)
	}
}

func (p *Params) IsEmpty() bool {
	return p == nil || len(*p) == 0
}
