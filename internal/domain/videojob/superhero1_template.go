package job

import "image"

type superHero1 struct {
}

type SuperHero1Config struct {
	ImageKid image.Image
}

func (s *superHero1) Process(cnfg SuperHero1Config) (string, error) {
	return "", nil
}
