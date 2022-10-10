package user

type RepositoryInterface interface {
	Get(id uint32) (User, error)
	GetByLogin(login string) (User, error)
	Store(login, password string) (id uint32, err error)
}
