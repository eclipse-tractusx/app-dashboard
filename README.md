# Installed Versions - Dashboard

This Version Dashboard provides an overview of all installed applications
through [ArgoCD](https://github.com/argoproj/argo-cd/). The main purpose is to give everyone non ArgoCD focused
stakeholder (Release Mgmt, Test Mgmt, Product Teams, etc.) an insight into the cluster state.

## How to install

- Use Helm chart under /chart

## Development overview

- /.github -> Github actions for building/testing
- /chart -> Helm Chart for installtion in k8s
- /web -> Asset folder for the dashboard web parts (css, img, js, ...)
- /Dockerfile -> distroless nonroot golang container image; Required by the Helm Chart
- /go.mod + /go.sum -> Go dependencies
- /main.go -> dashboard itself
  - http.ListenAndServe provides http webserver functionality
  - go func() -> Uses threads for serving the dashboard independent of the update function

### External dependencies

- https://fontawesome.com/docs/web/setup/host-yourself/webfonts
- https://github.com/fiduswriter/Simple-DataTables

### Additional dev thoughts

- The manual caching part could be done through nginx. Instead of delivering everything through the go code, only the
  dashboard itself could provide itself while nginx takes care of assets. Kept it as it is to show lowest dependency
  version

### How to maintain this code

- Keep Go kubernetes client up-to-date
  - https://github.com/kubernetes/client-go
- Update Go in the Dockerfile
- Update Simple-DataTable
  - Use the latest code from https://github.com/fiduswriter/Simple-DataTables through the CDN example (calling the CDN
    endpoint, which already contains the optimized version and save it locally)

### Development tips & tricks

- Use `go mod tidy` if there are dependency issues
- Font awesome was setup through https://fontawesome.com/docs/web/setup/host-yourself/webfonts
- How to get the object structure from the Kubernetes API & the proper URL/Path
  - `kubectl get applications.argoproj.io`
  - Use -v=8: `kubectl -v=8 get applications.argoproj.io -n argocd`

## Notice for Docker image

DockerHub: [https://hub.docker.com/r/tractusx/app-dashboard](https://hub.docker.com/r/tractusx/app-dashboard)

Eclipse Tractus-X product(s) installed within the image:

__App Dashboard__

- GitHub: https://github.com/eclipse-tractusx/app-dashboard
- Project home: https://projects.eclipse.org/projects/automotive.tractusx
- Dockerfile: https://github.com/eclipse-tractusx/app-dashboard/blob/main/Dockerfile
- Project license: [Apache License, Version 2.0](https://github.com/eclipse-tractusx/app-dashboard/blob/main/LICENSE)

__Used base image__

- [gcr.io/distroless/static:nonroot](https://github.com/GoogleContainerTools/distroless)

As with all Docker images, these likely also contain other software which may be under other licenses (such as Bash, etc
from the base distribution, along with any direct or indirect dependencies of the primary software being contained).

As for any pre-built image usage, it is the image user's responsibility to ensure that any use of this image complies
with any relevant licenses for all software contained within.
  