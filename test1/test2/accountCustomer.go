package test2

type AccountCustomerRepository struct {
	*AccountBaseRepository
}

func NewAccountRepository() *AccountCustomerRepository {
	return &AccountCustomerRepository{
		AccountBaseRepository: NewAccountBaseRepository(),
	}
}

func (a AccountCustomerRepository) GetByName(name string) string {
	return "哈哈哈" + name
}
