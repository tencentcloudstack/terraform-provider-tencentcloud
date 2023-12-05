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

Import

image instance can be imported using the id, e.g.

```
$ terraform import tencentcloud_image.image_snap img-gf7jspk6
```