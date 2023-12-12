Provides a resource to create a ckafka route

Example Usage

```hcl
resource "tencentcloud_ckafka_route" "route" {
	instance_id = "ckafka-xxxxxx"
	vip_type = 3
	vpc_id = "vpc-xxxxxx"
	subnet_id = "subnet-xxxxxx"
	access_type = 0
	public_network = 3
}
```