package devops

const DemoKey = "nice:devops"

type IService interface {
	GetAllStudent() []Student
}

type Student struct {
	ID   int
	Name string
}
