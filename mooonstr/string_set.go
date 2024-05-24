// Package mooonstr
// Wrote by yijian on 2024/05/24
package mooonstr

type StringSet map[string]struct{}

func NewStringSet() StringSet {
    return make(StringSet)
}

func NewStringSetWithValues(values []string) *StringSet {
    ss := make(StringSet)
    ss.BatchAdd(values)
    return &ss
}

func (s StringSet) Add(value string) {
    s[value] = struct{}{}
}

func (s StringSet) BatchAdd(values []string) {
    for _, value := range values {
        s.Add(value)
    }
}

func (s StringSet) Remove(value string) {
    delete(s, value)
}

func (s StringSet) BatchRemove(values []string) {
    for _, value := range values {
        s.Remove(value)
    }
}

func (s StringSet) Contains(value string) bool {
    _, exists := s[value]
    return exists
}

func (s StringSet) Count() int {
    return len(s)
}