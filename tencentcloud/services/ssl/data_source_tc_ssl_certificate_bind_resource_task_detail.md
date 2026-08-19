Use this data source to query detailed information of SSL certificate bind resource task detail, including the bind details of the certificate across various cloud resources (CLB, CDN, WAF, DDoS, Live, VOD, TKE, APIGateway, TCB, TEO, TSE, COS, TDMQ, MQTT, GAAP, SCF).

Example Usage

Query all resource types bind details by task id

```hcl
data "tencentcloud_ssl_certificate_bind_resource_task_detail" "example" {
  task_id = "51403728"
}
```

Query specific resource types bind details

```hcl
data "tencentcloud_ssl_certificate_bind_resource_task_detail" "example" {
  task_id        = "51403728"
  resource_types = ["clb", "cdn"]
}
```

Query bind details with regions filter

```hcl
data "tencentcloud_ssl_certificate_bind_resource_task_detail" "example" {
  task_id        = "51403728"
  resource_types = ["clb", "cdn"]
  regions        = ["ap-guangzhou", "ap-beijing"]
}
```
