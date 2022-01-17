package demo

import (
	devopsService "github.com/go-ddh/nice/app/provider/devops"
)

func UserModelsToUserDTOs(models []UserModel) []UserDTO {
	var ret []UserDTO
	for _, model := range models {
		t := UserDTO{
			ID:   model.UserId,
			Name: model.Name,
		}
		ret = append(ret, t)
	}
	return ret
}

func StudentsToUserDTOs(students []devopsService.Student) []UserDTO {
	var ret []UserDTO
	for _, student := range students {
		t := UserDTO{
			ID:   student.ID,
			Name: student.Name,
		}
		ret = append(ret, t)
	}
	return ret
}
