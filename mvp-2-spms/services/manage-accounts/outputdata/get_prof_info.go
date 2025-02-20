package outputdata

import (
	entities "mvp-2-spms/domain-aggregate"
	"strconv"
)

type GetProfessorInfo struct {
	Id            int    `json:"id"`
	Login         string `json:"login"`
	Name          string `json:"name"`
	ScienceDegree string `json:"science_degree"`
	University    string `json:"university"`
}

func MapToGetAccountInfo(prof entities.Professor, uni entities.University) GetProfessorInfo {
	pId, _ := strconv.Atoi(prof.Id)

	return GetProfessorInfo{
		Id:            pId,
		Login:         "", ////////////////////////////
		Name:          prof.FullNameToString(),
		ScienceDegree: prof.ScienceDegree,
		University:    uni.Name,
	}
}
