/*
Provides a resource to create a mysql switch_proxy

Example Usage

```hcl
resource "tencentcloud_mysql_switch_proxy" "switch_proxy" {
  instance_id = "cdb-fitq5t9h"
  proxy_group_id = "proxy-h1ub486b"
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

func resourceTencentCloudMysqlSwitchProxy() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMysqlSwitchProxyCreate,
		Read:   resourceTencentCloudMysqlSwitchProxyRead,
		Delete: resourceTencentCloudMysqlSwitchProxyDelete,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance id.",
			},

			"proxy_group_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Proxy group id.",
			},
		},
	}
}

func resourceTencentCloudMysqlSwitchProxyCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_switch_proxy.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var (
		request      = mysql.NewSwitchCDBProxyRequest()
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

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMysqlClient().SwitchCDBProxy(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate mysql switchProxy failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(instanceId + FILED_SP + proxyGroupId)

	service := MysqlService{client: meta.(*TencentCloudClient).apiV3Conn}
	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		proxy, err := service.DescribeMysqlProxyById(ctx, instanceId, proxyGroupId)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if *proxy.Status != "online" {
			return resource.RetryableError(fmt.Errorf("%s Switch mysql proxy status is %s", instanceId, *proxy.Status))
		}
		err = fmt.Errorf("%s Switch mysql proxy status is %s,we won't wait for it finish", instanceId, *proxy.Status)
		return resource.NonRetryableError(err)
	})

	if err != nil {
		log.Printf("[CRITAL]%s Switch mysql proxy fail, reason:%s\n ", logId, err.Error())
		return err
	}

	return resourceTencentCloudMysqlSwitchProxyRead(d, meta)
}

func resourceTencentCloudMysqlSwitchProxyRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_switch_proxy.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudMysqlSwitchProxyDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_switch_proxy.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
