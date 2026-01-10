# Monitoring and Autonomous Release Integration Operator (Mario)

Yes, the acronym makes no sense. Let it be.

This little tool listens for webhooks and updates tag names in your gitops repo.

Example flow:

- Docker image gets pushed to Docker Hub
- Docker sends webhook to `https://example.com/webhook/some-configured-uuid` (This is where mario sits)
- Mario checks the validity of the request
- Mario updates the gitops repo with the latest tag
- Gitops system (e.g. argocd or flux) updates k8s cluster.

## Todo

- Test all the things
- Make mocks for validator and gitopser

