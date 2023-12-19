package cdb

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mysql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudMysqlSwitchProxy() *schema.Resource {
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
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_switch_proxy.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

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

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMysqlClient().SwitchCDBProxy(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate mysql switchProxy failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(instanceId + tccommon.FILED_SP + proxyGroupId)

	service := MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
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
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_switch_proxy.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudMysqlSwitchProxyDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_switch_proxy.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
