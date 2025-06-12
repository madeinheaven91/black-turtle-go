package models

type DBStudyEntity struct {
	Id int
	Api_id int
	Kind StudyEntityType
	Name string
}

type DBChat struct {
	Id int64 
	Kind string 
	Name string 
	Username *string
	StudyEntityID *int
	IsBanned bool
}
