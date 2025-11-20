package entities

type Employee struct {
	Id          int64   `json:"id"`
	Name        string  `json:"name"`
	Address     string  `json:"address"`
	PhoneNumber string  `json:"phone_number"`
	Balance     float64 `json:"balance"`
}

type EmployeeFilter struct {
	Id int64 `form:"id" json:"id" uri:"id"`
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
	Id          int64   `json:"id" binding:"required"`
	Name        string  `json:"name"`
	Address     string  `json:"address"`
	PhoneNumber string  `json:"phone_number"`
	Balance     float64 `json:"balance"`
}

type TopupBalance struct {
	Id      int64   `json:"id" binding:"required"`
	Balance float64 `json:"balance" binding:"required"`
	Token   string  `json:"token" binding:"required"`
}

type SendBalance struct {
	Id      int64   `json:"id" binding:"required"`
	ToId    int64   `json:"to_id" binding:"required"`
	Balance float64 `json:"balance" binding:"required"`
}

type EmployeeBalance struct {
	Id      int64   `json:"id"`
	Name    string  `json:"name"`
	Balance float64 `json:"balance"`
	Message string  `json:"message"`
}

type ReportSendBalance struct {
	TotalData           int     `json:"total_data"`
	TotalNilaiTransaksi float64 `json:"total_nilai_transaksi"`
	TotalSukses         int     `json:"total_sukses"`
	TotalGagal          int     `json:"total_gagal"`
	DataSukses          []EmployeeBalance
	DataGagal           []EmployeeBalance
}
