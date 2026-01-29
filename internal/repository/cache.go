package repository

type BigCacheRepo struct{}

func (c *BigCacheRepo) Set(key string, value []byte) error {
	// TODO: Implement BigCache set
	return nil
}

func (c *BigCacheRepo) Get(key string) ([]byte, error) {
	// TODO: Implement BigCache get
	return nil, nil
}

func (c *BigCacheRepo) Delete(key string) error {
	// TODO: Implement BigCache delete
	return nil
}
