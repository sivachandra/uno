package token

type KindSet map[Kind]bool

func (s KindSet) Add(t Kind) {
	s[t] = true
}

func (s KindSet) Contains(t Kind) bool {
	_, e := s[t]
	return e
}

func NewKindSet(tokens []Kind) KindSet {
	s := KindSet(make(map[Kind]bool))
	for _, t := range tokens {
		s.Add(t)
	}
	return s
}

