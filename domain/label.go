package domain

type Label interface {
	Key() string
	Value() interface{}
}

type BoolLabel struct {
	key   string
	value bool
}

func NewBoolLabel(key string, value bool) Label {
	return &BoolLabel{
		key:   key,
		value: value,
	}
}

func (b BoolLabel) Key() string {
	return b.key
}

func (b BoolLabel) Value() interface{} {
	return b.value
}

func Labels() []Label {
	// todo: replace dummy with real labels
	label := NewBoolLabel("bkey", true)
	return []Label{label}
}
