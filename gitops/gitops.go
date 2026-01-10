package gitops

type GitOpser interface {
	Update(filepath string, jsonpath string, newTag string) error
}
