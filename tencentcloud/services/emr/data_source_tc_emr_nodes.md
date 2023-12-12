Provides an available EMR for the user.

The EMR data source obtain the hardware node information by using the emr cluster ID.

Example Usage

```hcl
data "tencentcloud_emr_nodes" "my_emr_nodes" {
  node_flag="master"
  instance_id="emr-rnzqrleq"
}
```