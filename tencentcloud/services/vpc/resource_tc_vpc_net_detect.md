Provides a resource to create a vpc net_detect

Example Usage

Create a basic Net Detect

```hcl
data "tencentcloud_availability_zones" "zones" {}

resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-example"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  name              = "subnet-example"
  cidr_block        = "10.0.0.0/16"
  availability_zone = data.tencentcloud_availability_zones.zones.zones.0.name
}

resource "tencentcloud_vpc_net_detect" "example" {
  net_detect_name       = "tf-example"
  vpc_id                = tencentcloud_vpc.vpc.id
  subnet_id             = tencentcloud_subnet.subnet.id
  detect_destination_ip = [
    "10.0.0.1",
    "10.0.0.2",
  ]
}
```

If `next_hop_type` is `VPN`

```hcl
resource "tencentcloud_vpn_gateway" "vpn" {
  name      = "tf-example"
  bandwidth = 100
  zone      = data.tencentcloud_availability_zones.zones.zones.0.name
  type      = "SSL"
  vpc_id    = tencentcloud_vpc.vpc.id

  tags = {
    test = "test"
  }
}

resource "tencentcloud_vpc_net_detect" "example" {
  net_detect_name       = "tf-example"
  vpc_id                = tencentcloud_vpc.vpc.id
  subnet_id             = tencentcloud_subnet.subnet.id
  next_hop_type         = "VPN"
  next_hop_destination  = tencentcloud_vpn_gateway.vpn.id
  detect_destination_ip = [
    "192.16.10.10",
    "172.16.10.22",
  ]
}
```

If `next_hop_type` is `DIRECTCONNECT`

```hcl
resource "tencentcloud_dc_gateway" "example" {
  name                = "ci-cdg-vpc-test"
  network_instance_id = tencentcloud_vpc.vpc.id
  network_type        = "VPC"
  gateway_type        = "NAT"
}

resource "tencentcloud_vpc_net_detect" "example" {
  net_detect_name       = "tf-example"
  vpc_id                = tencentcloud_vpc.vpc.id
  subnet_id             = tencentcloud_subnet.subnet.id
  next_hop_type         = "DIRECTCONNECT"
  next_hop_destination  = tencentcloud_dc_gateway.example.id
  detect_destination_ip = [
    "192.16.10.10",
    "172.16.10.22",
  ]
}
```

If `next_hop_type` is `NAT`

```hcl
resource "tencentcloud_eip" "eip_example1" {
  name = "tf_nat_gateway_eip1"
}

resource "tencentcloud_eip" "eip_example2" {
  name = "tf_nat_gateway_eip2"
}

resource "tencentcloud_nat_gateway" "example" {
  name             = "tf_example_nat_gateway"
  vpc_id           = tencentcloud_vpc.vpc.id
  bandwidth        = 100
  max_concurrent   = 1000000
  assigned_eip_set = [
    tencentcloud_eip.eip_example1.public_ip,
    tencentcloud_eip.eip_example2.public_ip,
  ]
  tags = {
    tf_tag_key = "tf_tag_value"
  }
}

resource "tencentcloud_vpc_net_detect" "example" {
  net_detect_name       = "tf-example"
  vpc_id                = tencentcloud_vpc.vpc.id
  subnet_id             = tencentcloud_subnet.subnet.id
  next_hop_type         = "NAT"
  next_hop_destination  = tencentcloud_nat_gateway.example.id
  detect_destination_ip = [
    "192.16.10.10",
    "172.16.10.22",
  ]
}
```

If `next_hop_type` is `NORMAL_CVM`

```hcl
data "tencentcloud_images" "image" {
  image_type = ["PUBLIC_IMAGE"]
  os_name    = "TencentOS Server 3.2 (Final)"
}

data "tencentcloud_instance_types" "instance_types" {
  filter {
    name   = "zone"
    values = [data.tencentcloud_availability_zones.zones.zones.0.name]
  }

  filter {
    name   = "instance-family"
    values = ["S5"]
  }

  cpu_core_count   = 2
  exclude_sold_out = true
}

resource "tencentcloud_instance" "example" {
  instance_name     = "tf_example"
  availability_zone = data.tencentcloud_availability_zones.zones.zones.0.name
  image_id          = data.tencentcloud_images.image.images.0.image_id
  instance_type     = data.tencentcloud_instance_types.instance_types.instance_types.0.instance_type
  system_disk_type  = "CLOUD_PREMIUM"
  system_disk_size  = 50
  hostname          = "user"
  project_id        = 0
  vpc_id            = tencentcloud_vpc.vpc.id
  subnet_id         = tencentcloud_subnet.subnet.id
}

resource "tencentcloud_vpc_net_detect" "example" {
  net_detect_name       = "tf-example"
  vpc_id                = tencentcloud_vpc.vpc.id
  subnet_id             = tencentcloud_subnet.subnet.id
  next_hop_type         = "NORMAL_CVM"
  next_hop_destination  = tencentcloud_instance.example.private_ip
  detect_destination_ip = [
    "192.16.10.10",
    "172.16.10.22",
  ]
}
```

If `next_hop_type` is `CCN`

```hcl
resource "tencentcloud_ccn" "example" {
  name                 = "tf-example"
  description          = "desc."
  qos                  = "AU"
  charge_type          = "POSTPAID"
  bandwidth_limit_type = "OUTER_REGION_LIMIT"
}

resource "tencentcloud_vpc_net_detect" "example" {
  net_detect_name       = "tf-example"
  vpc_id                = tencentcloud_vpc.vpc.id
  subnet_id             = tencentcloud_subnet.subnet.id
  next_hop_type         = "CCN"
  next_hop_destination  = tencentcloud_ccn.example.id
  detect_destination_ip = [
    "172.10.0.1",
    "172.10.0.2",
  ]
}
```

If `next_hop_type` is `NONEXTHOP`

```hcl
resource "tencentcloud_vpc_net_detect" "example" {
  net_detect_name       = "tf-example"
  vpc_id                = tencentcloud_vpc.vpc.id
  subnet_id             = tencentcloud_subnet.subnet.id
  next_hop_type         = "NONEXTHOP"
  detect_destination_ip = [
    "10.0.0.1",
    "10.0.0.2",
  ]
}
```

Import

vpc net_detect can be imported using the id, e.g.

```
terraform import tencentcloud_vpc_net_detect.net_detect net_detect_id
```