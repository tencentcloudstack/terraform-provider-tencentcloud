package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mysql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudMysqlVerifyRootAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMysqlVerifyRootAccountCreate,
		Read:   resourceTencentCloudMysqlVerifyRootAccountRead,
		Delete: resourceTencentCloudMysqlVerifyRootAccountDelete,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The instance ID, in the format: cdb-c1nl9rpv, is the same as the instance ID displayed on the cloud database console page.",
			},

			"password": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The password of the ROOT account of the instance.",
			},
		},
	}
}

func resourceTencentCloudMysqlVerifyRootAccountCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_verify_root_account.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var (
		request    = mysql.NewVerifyRootAccountRequest()
		response   = mysql.NewVerifyRootAccountResponse()
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("password"); ok {
		request.Password = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMysqlClient().VerifyRootAccount(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate mysql verifyRootAccount failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(instanceId)

	asyncRequestId := *response.Response.AsyncRequestId
	service := MysqlService{client: meta.(*TencentCloudClient).apiV3Conn}
	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		taskStatus, message, err := service.DescribeAsyncRequestInfo(ctx, asyncRequestId)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if taskStatus == MYSQL_TASK_STATUS_SUCCESS {
			return nil
		}
		if taskStatus == MYSQL_TASK_STATUS_INITIAL || taskStatus == MYSQL_TASK_STATUS_RUNNING {
			return resource.RetryableError(fmt.Errorf("%s verify rootAccount status is %s", instanceId, taskStatus))
		}
		err = fmt.Errorf("%s verify rootAccount status is %s,we won't wait for it finish ,it show message:%s", instanceId, taskStatus, message)
		return resource.NonRetryableError(err)
	})

	if err != nil {
		log.Printf("[CRITAL]%s verify rootAccount fail, reason:%s\n ", logId, err.Error())
		return err
	}

	return resourceTencentCloudMysqlVerifyRootAccountRead(d, meta)
}

func resourceTencentCloudMysqlVerifyRootAccountRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_verify_root_account.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudMysqlVerifyRootAccountDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_verify_root_account.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
