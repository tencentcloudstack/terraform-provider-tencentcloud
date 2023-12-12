Provides a resource to create a tem environment

Example Usage

```hcl
resource "tencentcloud_tem_environment" "environment" {
  environment_name = "demo"
  description      = "demo for test"
  vpc              = "vpc-2hfyray3"
  subnet_ids       = ["subnet-rdkj0agk", "subnet-r1c4pn5m", "subnet-02hcj95c"]
  tags = {
    "created" = "terraform"
  }
}

```
Import

tem environment can be imported using the id, e.g.
```
$ terraform import tencentcloud_tem_environment.environment environment_id
```