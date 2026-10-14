---
subcategory: "Serverless Cloud Function(SCF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_scf_layer"
sidebar_current: "docs-tencentcloud-resource-scf_layer"
description: |-
  Provide a resource to create a SCF layer.
---

# tencentcloud_scf_layer

Provide a resource to create a SCF layer.

## Example Usage

```hcl
resource "tencentcloud_scf_layer" "foo" {
  layer_name          = "foo"
  compatible_runtimes = ["Python3.6"]
  content {
    cos_bucket_name   = "test-bucket"
    cos_object_name   = "/foo.zip"
    cos_bucket_region = "ap-guangzhou"
  }
  description  = "foo"
  license_info = "foo"
}
```

## Argument Reference

The following arguments are supported:

* `compatible_runtimes` - (Required, List: [`String`]) The compatible runtimes of layer.
* `content` - (Required, List) The source code of layer.
* `layer_name` - (Required, String) The name of layer.
* `description` - (Optional, String) The description of layer.
* `license_info` - (Optional, String) The license info of layer.

The `content` object supports the following:

* `cos_bucket_name` - (Optional, String) Cos bucket name of the SCF layer, such as `cos-1234567890`, conflict with `zip_file`.
* `cos_bucket_region` - (Optional, String) Cos bucket region of the SCF layer, conflict with `zip_file`.
* `cos_object_name` - (Optional, String) Cos object name of the SCF layer, should have suffix `.zip` or `.jar`, conflict with `zip_file`.
* `zip_file` - (Optional, String) Zip file of the SCF layer, conflict with `cos_bucket_name`, `cos_object_name`, `cos_bucket_region`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `code_sha_256` - The code type of layer.
* `create_time` - The create time of layer.
* `layer_version` - The version of layer.
* `location` - The download location url of layer.
* `status` - The current status of layer.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [GetLayerVersion](https://cloud.tencent.com/document/api/583/41386) - This API is used to get the layer version details, including links used to download files in the layer.
* Create: [GetSecretValue](https://cloud.tencent.com/document/api/1140/40522) - For user-defined credentials, this API is used to get the plaintext information of a credential by specifying the credential name and version.
* Create: [PublishLayerVersion](https://cloud.tencent.com/document/api/583/41383) - This API is used to create a version for a layer by using the given .zip file or COS object. Each time this API is called with the same layer name, a new version will be generated.
* Read: [GetLayerVersion](https://cloud.tencent.com/document/api/583/41386) - This API is used to get the layer version details, including links used to download files in the layer.
* Delete: [DeleteLayerVersion](https://cloud.tencent.com/document/api/583/41387) - This API is used to delete a specified version of a specified layer. The deleted version cannot be associated with a function, but the deletion does not affect functions that are referencing this layer.

## Import

Scf layer can be imported, e.g.

```
$ terraform import tencentcloud_scf_layer.layer layerId#layerVersion
```

