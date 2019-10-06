package repo

type RepositoryResource interface {
	InitRepo(name string) error
}
