package models

type DBStudyEntity struct {
	ID int
	ApiID int
	Kind StudyEntityType
	Name string
}

type DBChat struct {
	ID int64 
	Kind string 
	Name string 
	Username *string
	StudyEntityID *int
	IsBanned bool
}
