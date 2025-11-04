package entities

type Employee struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phone_number"`
}

type EmployeeFilter struct {
	Id int64 `form:"id" json:"id"`
}

type GetId struct {
	Id int64 `json:"id"`
}

type CreateEmployee struct {
	Name        string `json:"name" binding:"required"`
	Address     string `json:"address" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
}

type UpdateEmployee struct {
	Id          int64  `json:"id" binding:"required"`
	Name        string `json:"name"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phone_number"`
}
