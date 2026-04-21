package test2

type IAccountBaseRepository interface {
	GetByName(name string) string
}

type AccountBaseRepository struct {
}

func (a AccountBaseRepository) GetByName(name string) string {
	return "你好" + name
}

func NewAccountBaseRepository() *AccountBaseRepository {
	return &AccountBaseRepository{}
}
