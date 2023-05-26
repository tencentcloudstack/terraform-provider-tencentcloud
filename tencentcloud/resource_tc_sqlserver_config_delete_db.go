/*
Provides a resource to create a sqlserver config_delete_db

Example Usage

```hcl
resource "tencentcloud_sqlserver_config_delete_db" "config_delete_db" {
  instance_id = "mssql-i1z41iwd"
  name =
}
```
*/
package tencentcloud

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudSqlserverConfigDeleteDB() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSqlserverConfigDeleteDBCreate,
		Read:   resourceTencentCloudSqlserverConfigDeleteDBRead,
		Update: resourceTencentCloudSqlserverConfigDeleteDBUpdate,
		Delete: resourceTencentCloudSqlserverConfigDeleteDBDelete,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},
			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "collection of database name.",
			},
		},
	}
}

func resourceTencentCloudSqlserverConfigDeleteDBCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_config_delete_db.create")()
	defer inconsistentCheck(d, meta)()

	var (
		instanceId string
		name       string
	)

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	}

	d.SetId(strings.Join([]string{instanceId, name}, FILED_SP))

	return resourceTencentCloudSqlserverConfigDeleteDBUpdate(d, meta)
}

func resourceTencentCloudSqlserverConfigDeleteDBRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_config_delete_db.read")()
	defer inconsistentCheck(d, meta)()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken, id is %s", d.Id())
	}
	instanceId := idSplit[0]
	name := idSplit[1]

	_ = d.Set("instance_id", instanceId)
	_ = d.Set("name", name)

	return nil
}

func resourceTencentCloudSqlserverConfigDeleteDBUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_config_delete_db.update")()
	defer inconsistentCheck(d, meta)()

	var (
		logId       = getLogId(contextNil)
		request     = sqlserver.NewDeleteDBRequest()
		flowRequest = sqlserver.NewDescribeFlowStatusRequest()
		flowId      int64
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken, id is %s", d.Id())
	}
	instanceId := idSplit[0]

	request.InstanceId = &instanceId

	if v, ok := d.GetOk("name"); ok {
		request.Names = []*string{helper.String(v.(string))}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSqlserverClient().DeleteDB(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		flowId = *result.Response.FlowId
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update sqlserver configDeleteDB failed, reason:%+v", logId, err)
		return err
	}

	flowRequest.FlowId = &flowId
	err = resource.Retry(10*writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSqlserverClient().DescribeFlowStatus(flowRequest)
		if e != nil {
			return retryError(e)
		}

		if *result.Response.Status == SQLSERVER_TASK_SUCCESS {
			return nil
		} else if *result.Response.Status == SQLSERVER_TASK_RUNNING {
			return resource.RetryableError(fmt.Errorf("sqlserver configDeleteDB status is running"))
		} else if *result.Response.Status == int64(SQLSERVER_TASK_FAIL) {
			return resource.NonRetryableError(fmt.Errorf("sqlserver configDeleteDB status is fail"))
		} else {
			e = fmt.Errorf("sqlserver configDeleteDB status illegal")
			return resource.NonRetryableError(e)
		}
	})

	if err != nil {
		log.Printf("[CRITAL]%s sqlserver configDeleteDB failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudSqlserverConfigDeleteDBRead(d, meta)
}

func resourceTencentCloudSqlserverConfigDeleteDBDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_config_delete_db.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
