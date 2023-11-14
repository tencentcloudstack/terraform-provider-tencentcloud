/*
Provides a resource to create a dcdb modify_instance_project

Example Usage

```hcl
resource "tencentcloud_dcdb_modify_instance_project" "modify_instance_project" {
  instance_ids =
  project_id =
}
```

Import

dcdb modify_instance_project can be imported using the id, e.g.

```
terraform import tencentcloud_dcdb_modify_instance_project.modify_instance_project modify_instance_project_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dcdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dcdb/v20180411"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudDcdbModifyInstanceProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDcdbModifyInstanceProjectCreate,
		Read:   resourceTencentCloudDcdbModifyInstanceProjectRead,
		Delete: resourceTencentCloudDcdbModifyInstanceProjectDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_ids": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of IDs of instances to be modified. The ID is in the format of `tdsql-ow728lmc` and can be obtained through the `DescribeDBInstances` API.",
			},

			"project_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "ID of the project to be assigned, which can be obtained through the `DescribeProjects` API.",
			},
		},
	}
}

func resourceTencentCloudDcdbModifyInstanceProjectCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_modify_instance_project.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = dcdb.NewModifyDBInstancesProjectRequest()
		response   = dcdb.NewModifyDBInstancesProjectResponse()
		instanceId string
	)
	if v, ok := d.GetOk("instance_ids"); ok {
		instanceIdsSet := v.(*schema.Set).List()
		for i := range instanceIdsSet {
			instanceIds := instanceIdsSet[i].(string)
			request.InstanceIds = append(request.InstanceIds, &instanceIds)
		}
	}

	if v, _ := d.GetOk("project_id"); v != nil {
		request.ProjectId = helper.IntInt64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDcdbClient().ModifyDBInstancesProject(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate dcdb modifyInstanceProject failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.InstanceId
	d.SetId(instanceId)

	return resourceTencentCloudDcdbModifyInstanceProjectRead(d, meta)
}

func resourceTencentCloudDcdbModifyInstanceProjectRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_modify_instance_project.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudDcdbModifyInstanceProjectDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_modify_instance_project.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
