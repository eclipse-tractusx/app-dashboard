# Installed Versions - Dashboard

This Version Dashboard provides an overview of all installed applications through [ArgoCD](https://github.com/argoproj/argo-cd/). The main purpose is to give everyone non ArgoCD focused stakeholder (Release Mgmt, Test Mgmt, Product Teams, etc.) an insight into the cluster state.

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
- The manual caching part could be done through nginx. Instead of delivering everything through the go code, only the dashboard itself could provide itself while nginx takes care of assets. Kept it as it is to show lowest dependency version

### How to maintain this code

- Keep Go kubernetes client up-to-date
  - https://github.com/kubernetes/client-go
- Update Go in the Dockerfile
- Update Simple-DataTable
  - Use the latest code from https://github.com/fiduswriter/Simple-DataTables through the CDN example (calling the CDN endpoint, which already contains the optimized version and save it locally)

### Development tips & tricks
- Use `go mod tidy` if there are dependency issues
- Font awesome was setup through https://fontawesome.com/docs/web/setup/host-yourself/webfonts
- How to get the object structure from the Kubernetes API & the proper URL/Path
  - `kubectl get applications.argoproj.io`
  - Use -v=8: `kubectl -v=8 get applications.argoproj.io -n argocd`
