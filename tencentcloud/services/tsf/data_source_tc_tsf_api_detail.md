Use this data source to query detailed information of tsf api_detail

Example Usage

```hcl
data "tencentcloud_tsf_api_detail" "api_detail" {
  microservice_id = "ms-yq3jo6jd"
  path = "/printRequest"
  method = "GET"
  pkg_version = "20210625192923"
  application_id = "application-a24x29xv"
}
```