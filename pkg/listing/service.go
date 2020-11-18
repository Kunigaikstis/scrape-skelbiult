package listing

type Service struct {
	repo Repository
}

type Repository interface {
	GetById(id int) (Ad, error)
	Add(listing Ad) error
}

func NewService(repo *Storage) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetNewListings(searchUrl string, scraper *Scraper) ([]Ad, error) {
	latestAds := scraper.GetListings(searchUrl)
	newAds := make([]Ad, 0)

	for _, newAd := range latestAds {
		_, err := s.repo.GetById(newAd.Id)

		adFoundInDB := err == nil
		if !adFoundInDB {
			err = s.repo.Add(newAd)

			if err != nil {
				return nil, err
			}

			newAds = append(newAds, newAd)
		}
	}

	return newAds, nil
}
