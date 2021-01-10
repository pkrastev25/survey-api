package db

const (
	DbSurvey = "survey"
)

const (
	PropertyId           = "_id"
	PropertyCreated      = "created"
	PropertyLastModified = "last_modified"
)

const (
	GreaterThan        = "$gt"
	GreaterThanOrEqual = "$gte"
	LessThan           = "$lt"
	LessThanOrEqual    = "$lte"
)

const (
	operationSet       = "$set"
	operationAddToSet  = "$addToSet"
	operationIncrement = "$inc"
	operationNotIn     = "$nin"
	operationSort      = "$sort"
	operationLimit     = "$limit"
)

const (
	operationMatch  = "$match"
	operationFacet  = "$facet"
	operationLookUp = "$lookup"
	operationText   = "$text"
	operationSearch = "$search"
)
