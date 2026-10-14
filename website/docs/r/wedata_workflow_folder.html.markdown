---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_workflow_folder"
sidebar_current: "docs-tencentcloud-resource-wedata_workflow_folder"
description: |-
  Provides a resource to create a wedata wedata_workflow_folder
---

# tencentcloud_wedata_workflow_folder

Provides a resource to create a wedata wedata_workflow_folder

## Example Usage

```hcl
resource "tencentcloud_wedata_workflow_folder" "wedata_workflow_folder" {
  project_id         = 2905622749543821312
  parent_folder_path = "/"
  folder_name        = "test"
}
```

## Argument Reference

The following arguments are supported:

* `folder_name` - (Required, String) Name of the folder to create.
* `parent_folder_path` - (Required, String) The absolute path of the parent folder, such as/abc/de, if it is the root directory, pass/.
* `project_id` - (Required, String, ForceNew) Project id.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateWorkflowFolder](https://cloud.tencent.com/document/api/1267/123638) - This API is used to create a workflow folder.
* Read: [ListWorkflowFolders](https://cloud.tencent.com/document/api/1267/123621) - This API is used to query the folder list.
* Update: [UpdateWorkflowFolder](https://cloud.tencent.com/document/api/1267/123614)
* Delete: [DeleteWorkflowFolder](https://cloud.tencent.com/document/api/1267/123633) - This API is used to delete a data development folder

