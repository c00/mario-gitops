package gitops

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_updateYaml(t *testing.T) {
	got, err := updateYaml(baseYaml, "$.deeper.still.str", "updated")
	assert.Nil(t, err)
	expected := "deeper:\n  still:\n    str: updated\nintProp: 1\nintSlice:\n  - 1\n  - 2\n  - 3\nstrProp: foo\nstrSlice:\n  - foo\n  - bar\n  - baz\n"

	assert.Equal(t, expected, got)
}

func Test_updateYamlKustomization(t *testing.T) {
	const kustomization = `
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
namespace: prod
images:
  - name: crispyduck/nap-api
    newTag: build-1
`
	const expected = `apiVersion: kustomize.config.k8s.io/v1beta1
images:
  - name: crispyduck/nap-api
    newTag: build-2
kind: Kustomization
namespace: prod
`
	got, err := updateYaml(kustomization, "$.images[?(@.name=='crispyduck/nap-api')].newTag", "build-2")
	assert.Nil(t, err)

	assert.Equal(t, expected, got)
}

func Test_updateYamlHelmValues(t *testing.T) {
	const kustomization = `image:
  repository: docker.io/crispyduck/nap-keycloak-extension
  tag: build-1
`
	const expected = `image:
  repository: docker.io/crispyduck/nap-keycloak-extension
  tag: build-2
`
	got, err := updateYaml(kustomization, "$.image.tag", "build-2")
	assert.Nil(t, err)

	assert.Equal(t, expected, got)
}

const baseYaml = `
deeper: 
  still:
    str: flop
intProp: 1
intSlice: [1, 2, 3]
strProp: foo
strSlice: [foo, bar, baz]
`
