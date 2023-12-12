Use this data source to query detailed information of tsf application

Example Usage

```hcl
data "tencentcloud_tsf_application" "application" {
  application_type = "V"
  microservice_type = "N"
  # application_resource_type_list = [""]
  application_id_list = ["application-a24x29xv"]
}
```