/*
Provides a resource to create a mariadb create_tmp_instance

Example Usage

```hcl
resource "tencentcloud_mariadb_create_tmp_instance" "create_tmp_instance" {
  instance_id   = "tdsql-9vqvls95"
  rollback_time = "2023-06-05 01:00:00"
}
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	mariadb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mariadb/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudMariadbCreateTmpInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMariadbCreateTmpInstanceCreate,
		Read:   resourceTencentCloudMariadbCreateTmpInstanceRead,
		Delete: resourceTencentCloudMariadbCreateTmpInstanceDelete,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},
			"rollback_time": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Rollback time.",
			},
		},
	}
}

func resourceTencentCloudMariadbCreateTmpInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_create_tmp_instance.create")()
	defer inconsistentCheck(d, meta)()

	var (
		logId      = getLogId(contextNil)
		ctx        = context.WithValue(context.TODO(), logIdKey, logId)
		service    = MariadbService{client: meta.(*TencentCloudClient).apiV3Conn}
		request    = mariadb.NewCreateTmpInstancesRequest()
		instanceId string
		flowId     int64
	)

	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceIds = common.StringPtrs([]string{v.(string)})
		instanceId = v.(string)
	}

	if v, ok := d.GetOk("rollback_time"); ok {
		request.RollbackTime = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMariadbClient().CreateTmpInstances(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		flowId = *result.Response.FlowId
		return nil

	})

	if err != nil {
		log.Printf("[CRITAL]%s operate mariadb createTmpInstance failed, reason:%+v", logId, err)
		return err
	}

	err = resource.Retry(10*writeRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeFlowById(ctx, flowId)
		if e != nil {
			return retryError(e)
		}

		if *result.Status == MARIADB_TASK_SUCCESS {
			return nil
		} else if *result.Status == MARIADB_TASK_RUNNING {
			return resource.RetryableError(fmt.Errorf("operate mariadb createTmpInstance status is running"))
		} else if *result.Status == MARIADB_TASK_FAIL {
			return resource.NonRetryableError(fmt.Errorf("operate mariadb createTmpInstance status is fail"))
		} else {
			e = fmt.Errorf("operate mariadb createTmpInstance status illegal")
			return resource.NonRetryableError(e)
		}
	})

	if err != nil {
		log.Printf("[CRITAL]%s operate mariadb createTmpInstance task failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(instanceId)

	return resourceTencentCloudMariadbCreateTmpInstanceRead(d, meta)
}

func resourceTencentCloudMariadbCreateTmpInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_create_tmp_instance.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudMariadbCreateTmpInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_create_tmp_instance.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
