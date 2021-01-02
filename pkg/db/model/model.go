package model

const (
	GreaterThan        Operation = "$gt"
	GreaterThanOrEqual Operation = "$gte"
	LessThan           Operation = "$lt"
	LessThanOrEqual    Operation = "$lte"
	Limit              Operation = "$limit"
	Facet              Operation = "$facet"
	Match              Operation = "$match"
	And                Operation = "$and"
	Text               Operation = "$text"
	Search             Operation = "$search"
	Sort               Operation = "$sort"
)

const (
	Created Property = "created"
)

type Property string
type Operation string
