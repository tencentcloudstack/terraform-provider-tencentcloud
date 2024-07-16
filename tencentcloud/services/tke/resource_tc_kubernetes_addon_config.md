Provide a resource to configure addon that kubernetes comes with.

Example Usage

Update cluster-autoscaler addon

```hcl

resource "tencentcloud_kubernetes_addon_config" "kubernetes_addon_config" {
	cluster_id = "cls-xxxxxx"
	addon_name = "cluster-autoscaler"
	raw_values = "{\"extraArgs\":{\"scale-down-enabled\":true,\"max-empty-bulk-delete\":11,\"scale-down-delay-after-add\":\"10mm\",\"scale-down-unneeded-time\":\"10mm\",\"scale-down-utilization-threshold\":0.005,\"ignore-daemonsets-utilization\":false,\"skip-nodes-with-local-storage\":true,\"skip-nodes-with-system-pods\":true}}"
}
`

```