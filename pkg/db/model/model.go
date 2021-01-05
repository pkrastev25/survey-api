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
	Set                Operation = "$set"
	AddToSet           Operation = "$addToSet"
	Increment          Operation = "$inc"
	NotIn              Operation = "$nin"
	LookUp             Operation = "$lookup"
)

const (
	DB string = "survey"

	UserCollection    string = "user"
	SessionCollection string = "session"

	Created Property = "created"
)

type Property string
type Operation string
