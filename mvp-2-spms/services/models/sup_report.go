package models

import domainaggregate "mvp-2-spms/domain-aggregate"

type SupReport struct {
	Faculty       string     //
	Dept          string     //
	StudentName   string     //
	Course        string     //
	EdProgramme   string     //
	Theme         string     //
	Items         []Ctiteria //
	Comment       string     //
	ScienceDegree string     //
	SupRewGrade   string     //
	ProfName      string     //
	Date          string     //
}

type Ctiteria struct {
	Num   string
	Name  string
	Grade string
}

type EdProg struct {
	domainaggregate.EducationalProgramme
	Faculty string
	Dept    string
}
