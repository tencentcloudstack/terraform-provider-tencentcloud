# CLB Instances Data Source - Zones Field Example

This example demonstrates the usage of the `zones` field in the `tencentcloud_clb_instances` data source.

## Overview

The `zones` field contains zone information for VPC internal load balancers with nearby access mode. This field indicates which availability zones the CLB rules are deployed to.

## Features Demonstrated

1. **Zones Field**: Query CLB instances and retrieve their zones information
2. **Retry Logic**: The underlying API calls now have automatic retry capability for improved reliability

## Usage

Configure your Tencent Cloud credentials and run:

```bash
terraform init
terraform plan
terraform apply
```

## Example Output

### Internal CLB with Zones

```hcl
clb_instances_with_zones = [
  {
    clb_id       = "lb-xxxxxxxx"
    clb_name     = "my-internal-lb"
    network_type = "INTERNAL"
    vpc_id       = "vpc-xxxxxxxx"
    zones        = ["ap-guangzhou-1", "ap-guangzhou-2"]
  },
]
```

### CLB without Zones (null)

For CLB instances that don't have zone information (e.g., public CLB, traditional CLB), the `zones` field will be `null`:

```hcl
all_clb_instances = [
  {
    clb_id       = "lb-yyyyyyyy"
    clb_name     = "my-public-lb"
    network_type = "OPEN"
    zones        = null
  },
]
```

## Notes

- The `zones` field is only populated for VPC internal load balancers with nearby access mode
- For other CLB types, the field may return `null`
- The retry logic improvement is transparent and requires no configuration changes

## API Reference

- Data Source: [tencentcloud_clb_instances](https://registry.terraform.io/providers/tencentcloudstack/tencentcloud/latest/docs/data-sources/clb_instances)
- API: [DescribeLoadBalancers](https://cloud.tencent.com/document/product/214/30685)
