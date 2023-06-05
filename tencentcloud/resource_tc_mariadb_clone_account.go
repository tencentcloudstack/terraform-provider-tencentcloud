/*
Provides a resource to create a mariadb clone_account

Example Usage

```hcl
resource "tencentcloud_mariadb_clone_account" "clone_account" {
  instance_id = "tdsql-9vqvls95"
  src_user = "srcuser"
  src_host = "10.13.1.%"
  dst_user = "dstuser"
  dst_host = "10.13.23.%"
  dst_desc = "test clone"
}
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mariadb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mariadb/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudMariadbCloneAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMariadbCloneAccountCreate,
		Read:   resourceTencentCloudMariadbCloneAccountRead,
		Delete: resourceTencentCloudMariadbCloneAccountDelete,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},
			"src_user": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Source user account name.",
			},
			"src_host": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Source user host.",
			},
			"dst_user": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Target user account name.",
			},
			"dst_host": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Target user host.",
			},
			"dst_desc": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Target account description.",
			},
		},
	}
}

func resourceTencentCloudMariadbCloneAccountCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_clone_account.create")()
	defer inconsistentCheck(d, meta)()

	var (
		logId      = getLogId(contextNil)
		ctx        = context.WithValue(context.TODO(), logIdKey, logId)
		service    = MariadbService{client: meta.(*TencentCloudClient).apiV3Conn}
		request    = mariadb.NewCloneAccountRequest()
		instanceId string
		flowId     uint64
	)

	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = helper.String(v.(string))
		instanceId = v.(string)
	}

	if v, ok := d.GetOk("src_user"); ok {
		request.SrcUser = helper.String(v.(string))
	}

	if v, ok := d.GetOk("src_host"); ok {
		request.SrcHost = helper.String(v.(string))
	}

	if v, ok := d.GetOk("dst_user"); ok {
		request.DstUser = helper.String(v.(string))
	}

	if v, ok := d.GetOk("dst_host"); ok {
		request.DstHost = helper.String(v.(string))
	}

	if v, ok := d.GetOk("dst_desc"); ok {
		request.DstDesc = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMariadbClient().CloneAccount(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		flowId = *result.Response.FlowId
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s operate mariadb cloneAccount failed, reason:%+v", logId, err)
		return err
	}

	err = resource.Retry(10*writeRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeFlowById(ctx, int64(flowId))
		if e != nil {
			return retryError(e)
		}

		if *result.Status == MARIADB_TASK_SUCCESS {
			return nil
		} else if *result.Status == MARIADB_TASK_RUNNING {
			return resource.RetryableError(fmt.Errorf("operate mariadb cloneAccount status is running"))
		} else if *result.Status == MARIADB_TASK_FAIL {
			return resource.NonRetryableError(fmt.Errorf("operate mariadb cloneAccount status is fail"))
		} else {
			e = fmt.Errorf("operate mariadb cloneAccount status illegal")
			return resource.NonRetryableError(e)
		}
	})

	if err != nil {
		log.Printf("[CRITAL]%s operate mariadb cloneAccount task failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(instanceId)

	return resourceTencentCloudMariadbCloneAccountRead(d, meta)
}

func resourceTencentCloudMariadbCloneAccountRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_clone_account.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudMariadbCloneAccountDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_clone_account.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
