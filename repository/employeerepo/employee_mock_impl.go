package employeerepo

import (
	"router-template/entities"

	"github.com/go-faker/faker/v4"
)

func newEmployeeMockImpl() EmployeeRepo {
	return &mockImpl{}
}

type mockImpl struct{}

func (m mockImpl) GetEmployee() (list []entities.Employee, er error) {
	for i := 0; i < 100; i++ {
		item := entities.Employee{
			Id:          faker.RandomUnixTime(),
			Name:        faker.Name(),
			Address:     faker.GetRealAddress().Address,
			PhoneNumber: faker.E164PhoneNumber(),
		}

		list = append(list, item)
	}
	return
}
func (m mockImpl) GetEmployeeById(id int64) (employee entities.Employee, er error) {
	return entities.Employee{
		Id:          id,
		Name:        faker.Name(),
		Address:     faker.GetRealAddress().Address,
		PhoneNumber: faker.E164PhoneNumber(),
	}, nil
}

func (m mockImpl) CreateEmployee(name, address, phone_number string) (employee entities.Employee, er error) {
	return entities.Employee{
		Id:          faker.RandomUnixTime(),
		Name:        name,
		Address:     address,
		PhoneNumber: phone_number,
	}, nil
}

func (m mockImpl) UpdateEmployee(id int64, name, address, phone_number string) (employee entities.Employee, er error) {
	return entities.Employee{
		Id:          id,
		Name:        name,
		Address:     address,
		PhoneNumber: phone_number,
	}, nil
}

func (m mockImpl) DeleteEmployee(id int64) (employee entities.Employee, er error) {
	return entities.Employee{
		Id: id,
	}, nil
}
