Use this data source to query the health status of origin groups under a TEO (EdgeOne) load balancer instance.

Example Usage

Query the health status of all origin groups under a load balancer

```hcl
data "tencentcloud_teo_origin_group_health_status" "example" {
  zone_id        = "zone-3fkff38fyw8s"
  lb_instance_id = "lb-instance-xxx"
}
```

Query the health status of specified origin groups

```hcl
data "tencentcloud_teo_origin_group_health_status" "example" {
  zone_id         = "zone-3fkff38fyw8s"
  lb_instance_id  = "lb-instance-xxx"
  origin_group_ids = ["origin-group-xxx1", "origin-group-xxx2"]
}
```
