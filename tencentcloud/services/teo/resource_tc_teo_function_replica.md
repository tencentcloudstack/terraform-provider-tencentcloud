Provides a resource to create a TEO edge function replica

Example Usage

```hcl
resource "tencentcloud_teo_function_replica" "example" {
  zone_id      = "zone-2qtuhspy7cr6"
  function_id  = "ef-2qlxy8s7o96e"
  replica_name = "replica-example"
  content      = "addEventListener('fetch', event => { event.respondWith(new Response('hello world')) })"
  remark       = "example replica"
}
```

Import

TEO function replica can be imported using the zone_id#function_id#replica_name, e.g.

```
terraform import tencentcloud_teo_function_replica.example zone-2qtuhspy7cr6#ef-2qlxy8s7o96e#replica-example
```
