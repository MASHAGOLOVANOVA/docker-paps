package domainaggregate

type Professor struct {
	Person
	ScienceDegree   string
	AvailablePlaces []AvailablePlace // make private?
	UniversityId    string
}

type AvailablePlace struct {
	EducationalProgrammeId string
	Cource                 uint
	PlacesCount            uint
}
