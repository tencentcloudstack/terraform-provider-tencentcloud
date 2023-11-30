Provides a resource to create a tmp tke template

Example Usage

```hcl
resource "tencentcloud_monitor_tmp_tke_template" "foo" {
  template {
    name     = "tf-template"
    level    = "cluster"
    describe = "template"
    service_monitors {
      name   = "tf-ServiceMonitor"
      config = <<-EOT
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: example-service-monitor
  namespace: monitoring
  labels:
    k8s-app: example-service
spec:
  selector:
    matchLabels:
      k8s-app: example-service
  namespaceSelector:
    matchNames:
      - default
  endpoints:
  - port: http-metrics
    interval: 30s
    path: /metrics
    scheme: http
    bearerTokenFile: /var/run/secrets/kubernetes.io/serviceaccount/token
    tlsConfig:
      insecureSkipVerify: true
      EOT
    }

    pod_monitors {
      name   = "tf-PodMonitors"
      config = <<-EOT
apiVersion: monitoring.coreos.com/v1
kind: PodMonitor
metadata:
  name: example-pod-monitor
  namespace: monitoring
  labels:
    k8s-app: example-pod
spec:
  selector:
    matchLabels:
      k8s-app: example-pod
  namespaceSelector:
    matchNames:
      - default
  podMetricsEndpoints:
  - port: http-metrics
    interval: 30s
    path: /metrics
    scheme: http
    bearerTokenFile: /var/run/secrets/kubernetes.io/serviceaccount/token
    tlsConfig:
      insecureSkipVerify: true
EOT
    }

    pod_monitors {
      name   = "tf-RawJobs"
      config = <<-EOT
scrape_configs:
  - job_name: 'example-job'
    scrape_interval: 30s
    static_configs:
      - targets: ['example-service.default.svc.cluster.local:8080']
    metrics_path: /metrics
    scheme: http
    bearer_token_file: /var/run/secrets/kubernetes.io/serviceaccount/token
    tls_config:
      insecure_skip_verify: true
EOT
    }
  }
}
```