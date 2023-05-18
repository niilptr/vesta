package types

func NewTwin(name string, hash string, creator string) Twin {
	return Twin{
		Name:    name,
		Hash:    hash,
		Creator: creator,
	}
}
