package domain

type GetMemoRepoEtcd interface {
	Get(id string) (Memo, error)
	GetByUrl(url string) (Memo, error)
}

type GetMemoService struct {
	getEtcd    GetMemoRepoEtcd
	createEtcd CreateMemoRepoEtcd
}

func NewGetMemoService(getEtcd GetMemoRepoEtcd, createEtcd CreateMemoRepoEtcd) GetMemoService {
	return GetMemoService{
		getEtcd:    getEtcd,
		createEtcd: createEtcd,
	}
}

func (s GetMemoService) Get(id string) (Memo, error) {
	memo, err := s.getEtcd.Get(id)
	if err != nil {
		return Memo{}, err
	}
	s.createEtcd.Put(memo)

	return memo, nil
}
