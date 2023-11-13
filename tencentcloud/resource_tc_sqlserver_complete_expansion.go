/*
Provides a resource to create a sqlserver complete_expansion

Example Usage

```hcl
resource "tencentcloud_sqlserver_complete_expansion" "complete_expansion" {
  instance_id = "mssql-i1z41iwd"
}
```

Import

sqlserver complete_expansion can be imported using the id, e.g.

```
terraform import tencentcloud_sqlserver_complete_expansion.complete_expansion complete_expansion_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudSqlserverCompleteExpansion() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSqlserverCompleteExpansionCreate,
		Read:   resourceTencentCloudSqlserverCompleteExpansionRead,
		Delete: resourceTencentCloudSqlserverCompleteExpansionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "ID of imported target instance.",
			},
		},
	}
}

func resourceTencentCloudSqlserverCompleteExpansionCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_complete_expansion.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = sqlserver.NewCompleteExpansionRequest()
		response   = sqlserver.NewCompleteExpansionResponse()
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSqlserverClient().CompleteExpansion(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate sqlserver completeExpansion failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.InstanceId
	d.SetId(instanceId)

	return resourceTencentCloudSqlserverCompleteExpansionRead(d, meta)
}

func resourceTencentCloudSqlserverCompleteExpansionRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_complete_expansion.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudSqlserverCompleteExpansionDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_complete_expansion.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
