Use this resource to create tcr instance.

Example Usage

Create a basic tcr instance.

```hcl
resource "tencentcloud_tcr_instance" "example" {
  name              = "tf-example-tcr"
  instance_type		= "basic"

  tags = {
    "createdBy" = "terraform"
  }
}
```

Create instance with the public network access whitelist.

```hcl
resource "tencentcloud_tcr_instance" "example" {
  name                  = "tf-example-tcr"
  instance_type		    = "basic"
  open_public_operation = true
  security_policy {
    cidr_block = "10.0.0.1/24"
  }
  security_policy {
    cidr_block = "192.168.1.1"
  }
}
```

Create instance with Replications.

```hcl
resource "tencentcloud_tcr_instance" "example" {
  name                  = "tf-example-tcr"
  instance_type		    = "premium"
  replications {
    region_id = var.tcr_region_map["ap-guangzhou"] # 1
  }
  replications {
    region_id = var.tcr_region_map["ap-singapore"] # 9
  }
}

variable "tcr_region_map" {
  default = {
    "ap-guangzhou"     = 1
    "ap-shanghai"      = 4
    "ap-hongkong"      = 5
    "ap-beijing"       = 8
    "ap-singapore"     = 9
    "na-siliconvalley" = 15
    "ap-chengdu"       = 16
    "eu-frankfurt"     = 17
    "ap-seoul"         = 18
    "ap-chongqing"     = 19
    "ap-mumbai"        = 21
    "na-ashburn"       = 22
    "ap-bangkok"       = 23
    "eu-moscow"        = 24
    "ap-tokyo"         = 25
    "ap-nanjing"       = 33
    "ap-taipei"        = 39
    "ap-jakarta"       = 72
  }
}
```

Import

tcr instance can be imported using the id, e.g.

```
$ terraform import tencentcloud_tcr_instance.foo instance_id
```