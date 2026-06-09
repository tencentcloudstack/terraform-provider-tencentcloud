## 1. Schema Fix

- [x] 1.1 In `tencentcloud/services/tke/resource_tc_kubernetes_cluster_attachment.go`, locate the `docker_graph_path` field inside `worker_config` (around line 87): remove `Default: "/var/lib/docker"`, add `Computed: true`, and update `Description` to `"Docker graph path. Default is determined by the platform (currently /var/lib/containerd for containerd-based nodes)."`
- [x] 1.2 In the same file, locate the second `docker_graph_path` field inside `worker_config_overrides` (around line 269): remove `Default: "/var/lib/docker"`, add `Computed: true`, and update `Description` to `"Docker graph path. Default is determined by the platform (currently /var/lib/containerd for containerd-based nodes)."`
- [x] 1.3 Run `gofmt -w tencentcloud/services/tke/resource_tc_kubernetes_cluster_attachment.go` to format the file

## 2. Documentation

- [x] 2.1 Update `tencentcloud/services/tke/resource_tc_kubernetes_cluster_attachment.md` to remove references to `/var/lib/docker` as the default value for `docker_graph_path` in example configs
- [x] 2.2 Run `make doc` to regenerate `website/docs/` markdown documentation

## 3. Verification

- [x] 3.1 Run `go build ./tencentcloud/services/tke/` to verify the file compiles without errors
