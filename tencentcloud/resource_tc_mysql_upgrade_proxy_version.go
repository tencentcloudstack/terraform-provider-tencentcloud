/*
Provides a resource to create a mysql upgrade_proxy_version

Example Usage

```hcl
resource "tencentcloud_mysql_upgrade_proxy_version" "upgrade_proxy_version" {
  instance_id = ""
  proxy_group_id = ""
  src_proxy_version = ""
  dst_proxy_version = ""
  upgrade_time = ""
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
	mysql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudMysqlUpgradeProxyVersion() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMysqlUpgradeProxyVersionCreate,
		Read:   resourceTencentCloudMysqlUpgradeProxyVersionRead,
		Delete: resourceTencentCloudMysqlUpgradeProxyVersionDelete,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},

			"proxy_group_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Database proxy ID.",
			},

			"src_proxy_version": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The current version of the database agent.",
			},

			"dst_proxy_version": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Database agent upgrade version.",
			},

			"upgrade_time": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Upgrade time: nowTime (upgrade completed) timeWindow (instance maintenance time).",
			},
		},
	}
}

func resourceTencentCloudMysqlUpgradeProxyVersionCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_upgrade_proxy_version.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var (
		request      = mysql.NewUpgradeCDBProxyVersionRequest()
		response     = mysql.NewUpgradeCDBProxyVersionResponse()
		instanceId   string
		proxyGroupId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("proxy_group_id"); ok {
		proxyGroupId = v.(string)
		request.ProxyGroupId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("src_proxy_version"); ok {
		request.SrcProxyVersion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("dst_proxy_version"); ok {
		request.DstProxyVersion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("upgrade_time"); ok {
		request.UpgradeTime = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMysqlClient().UpgradeCDBProxyVersion(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate mysql upgradeProxyVersion failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(instanceId + FILED_SP + proxyGroupId)

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
			return resource.RetryableError(fmt.Errorf("%s upgrade mysql roxy version status is %s", instanceId, taskStatus))
		}
		err = fmt.Errorf("%s upgrade mysql roxy version status is %s,we won't wait for it finish ,it show message:%s", instanceId, taskStatus, message)
		return resource.NonRetryableError(err)
	})

	if err != nil {
		log.Printf("[CRITAL]%s upgrade mysql roxy version fail, reason:%s\n ", logId, err.Error())
		return err
	}

	return resourceTencentCloudMysqlUpgradeProxyVersionRead(d, meta)
}

func resourceTencentCloudMysqlUpgradeProxyVersionRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_upgrade_proxy_version.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudMysqlUpgradeProxyVersionDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_upgrade_proxy_version.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
