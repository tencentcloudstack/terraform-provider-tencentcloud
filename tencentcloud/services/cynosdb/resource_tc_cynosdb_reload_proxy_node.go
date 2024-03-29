package cynosdb

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCynosdbReloadProxyNode() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCynosdbReloadProxyNodeCreate,
		Read:   resourceTencentCloudCynosdbReloadProxyNodeRead,
		Delete: resourceTencentCloudCynosdbReloadProxyNodeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "cluster id.",
			},
			"proxy_group_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "proxy group id.",
			},
		},
	}
}

func resourceTencentCloudCynosdbReloadProxyNodeCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cynosdb_reload_proxy_node.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service   = CynosdbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		request   = cynosdb.NewReloadBalanceProxyNodeRequest()
		response  = cynosdb.NewReloadBalanceProxyNodeResponse()
		clusterId string
		flowId    int64
	)

	if v, ok := d.GetOk("cluster_id"); ok {
		request.ClusterId = helper.String(v.(string))
		clusterId = v.(string)
	}

	if v, ok := d.GetOk("proxy_group_id"); ok {
		request.ProxyGroupId = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCynosdbClient().ReloadBalanceProxyNode(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil {
			e = fmt.Errorf("cynosdb proxyNode not exists")
			return resource.NonRetryableError(e)
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s operate cynosdb proxyNode failed, reason:%+v", logId, err)
		return err
	}

	flowId = *response.Response.FlowId
	err = resource.Retry(6*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ok, err := service.DescribeFlow(ctx, flowId)
		if err != nil {
			if _, ok := err.(*sdkErrors.TencentCloudSDKError); !ok {
				return resource.RetryableError(err)
			} else {
				return resource.NonRetryableError(err)
			}
		}

		if ok {
			return nil
		} else {
			return resource.RetryableError(fmt.Errorf("cynosdb proxy is processing"))
		}
	})

	if err != nil {
		log.Printf("[CRITAL]%s create cynosdb proxy fail, reason:%s\n", logId, err.Error())
		return err
	}

	d.SetId(clusterId)

	return resourceTencentCloudCynosdbReloadProxyNodeRead(d, meta)
}

func resourceTencentCloudCynosdbReloadProxyNodeRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cynosdb_reload_proxy_node.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudCynosdbReloadProxyNodeDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cynosdb_reload_proxy_node.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
