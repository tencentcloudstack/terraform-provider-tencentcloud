Use this data source to query available accelerate regions of GA2 (Global Accelerator 2)

Example Usage

Query all GA2 accelerate regions

```hcl
data "tencentcloud_ga2_accelerate_regions" "example" {}
```

Query GA2 accelerate regions and output to file

```hcl
data "tencentcloud_ga2_accelerate_regions" "example" {
  result_output_file = "ga2_accelerate_regions.json"
}
```