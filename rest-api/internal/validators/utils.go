package validators

type stringList []string

func (c stringList) Contains(v string) bool {
	for _, a := range c {
		if a == v {
			return true
		}
	}

	return false
}
