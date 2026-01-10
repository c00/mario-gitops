package gitops

import "log/slog"

var _ GitOpser = (*MockOps)(nil)

type MockOps struct {
}

func (g *MockOps) Update(filepath string, jsonpath string, newTag string) error {
	slog.Info("Mockups.Update()", "filepath", filepath, "jsonpath", jsonpath, "newTag", newTag)
	return nil
}
