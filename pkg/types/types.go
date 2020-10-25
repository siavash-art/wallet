package types

//Money is type money
type Money int64

//PaymentCategory  payments categories
type PaymentCategory string

//PaymentStatus payments status
type PaymentStatus string

//Status categories
const (
	PaymentStatusOk         PaymentStatus = "OK"
	PaymentStatusFail       PaymentStatus = "FAIL"
	PaymentStatusInProgress PaymentStatus = "INPROGRESS"
)

//Payment struct
type Payment struct {
	ID       string
	AccountID int64
	Amount   Money
	Category PaymentCategory
	Status   PaymentStatus
}

//Phone payments phone
type Phone string

//Account struct
type Account struct {
	ID      int64
	Phone   Phone
	Balance Money
}

type Favorite struct {
	ID        string
	AccountID int64
	Name      string
	Amount    Money
	Category  PaymentCategory
}

type Progress struct {
	Part int
	Result Money
}