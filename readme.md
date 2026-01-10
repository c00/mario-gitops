# My Autonomous Release Integration Operator (Mario)

Yes, the acronym makes no sense. Let it be.

This little tool listens for webhooks and updates tag names in your gitops repo.

Example flow:

- Docker image gets pushed to Docker Hub
- Docker sends webhook to `https://example.com/webhook/some-configured-uuid` (This is where mario sits)
- Mario checks the validity of the request
- Mario updates the gitops repo with the latest tag
- Gitops system (e.g. argocd or flux) updates k8s cluster.

## Todo

- Test main
  - Create a script that runs the thing with mocks
  - And then throws a bunch of curls at it.

```sh
# In one terminal
MARIO_CONFIG_PATH="e2etest.yaml" go run ./cmd/server

# In another
curl -v http://localhost:8888/webhook/e2e-endpoint
# Should return 200 OK
```