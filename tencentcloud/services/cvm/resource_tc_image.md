Provide a resource to manage image.

Example Usage

```hcl
resource "tencentcloud_image" "image_snap" {
	image_name   		= "image-snapshot-keep"
	snapshot_ids 		= ["snap-nbp3xy1d", "snap-nvzu3dmh"]
	force_poweroff 		= true
	image_description 	= "create image with snapshot"
}
```

Use image family

```hcl
resource "tencentcloud_image" "image_family" {
  image_description = "create image with snapshot 12323"
  image_family      = "business-daily-update"
  image_name        = "image-family-test123"
  data_disk_ids     = []
  snapshot_ids      = [
    "snap-7uuvrcoj",
  ]
}
```

Import

image instance can be imported using the id, e.g.

```
$ terraform import tencentcloud_image.image_snap img-gf7jspk6
```