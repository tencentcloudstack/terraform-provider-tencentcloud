Provides an available EMR for the user.

The EMR data source fetch proper EMR from user's EMR pool.

Example Usage

```hcl
data "tencentcloud_emr" "my_emr" {
  display_strategy="clusterList"
  instance_ids=["emr-rnzqrleq"]
}
```