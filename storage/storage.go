package storage

type Storage interface {
	Store(string) string
	Retrieve(string) (string, error)
}
