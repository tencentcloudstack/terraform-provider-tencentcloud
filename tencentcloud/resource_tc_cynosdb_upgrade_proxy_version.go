/*
Provides a resource to create a cynosdb upgrade_proxy_version

Example Usage

```hcl
resource "tencentcloud_cynosdb_upgrade_proxy_version" "upgrade_proxy_version" {
  cluster_id = "cynosdbmysql-bws8h88b"
  dst_proxy_version = "1.3.7"
}
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
)

func resourceTencentCloudCynosdbUpgradeProxyVersion() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCynosdbUpgradeProxyVersionCreate,
		Read:   resourceTencentCloudCynosdbUpgradeProxyVersionRead,
		Update: resourceTencentCloudCynosdbUpgradeProxyVersionUpdate,
		Delete: resourceTencentCloudCynosdbUpgradeProxyVersionDelete,

		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Cluster ID.",
			},
			"dst_proxy_version": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Database Agent Upgrade Version.",
			},
		},
	}
}

func resourceTencentCloudCynosdbUpgradeProxyVersionCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_upgrade_proxy_version.create")()
	defer inconsistentCheck(d, meta)()

	var (
		logId           = getLogId(contextNil)
		ctx             = context.WithValue(context.TODO(), logIdKey, logId)
		service         = CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}
		srcProxyVersion string
		clusterId       string
	)

	if v, ok := d.GetOk("cluster_id"); ok {
		clusterId = v.(string)
	}

	upgradeProxyGroup, err := service.DescribeCynosdbUpgradeProxyVersionById(ctx, clusterId)
	if err != nil {
		return err
	}

	if upgradeProxyGroup == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CynosdbUpgradeProxyVersion` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	srcProxyVersion = *upgradeProxyGroup.CurrentProxyVersion

	d.SetId(strings.Join([]string{clusterId, srcProxyVersion}, FILED_SP))

	return resourceTencentCloudCynosdbUpgradeProxyVersionUpdate(d, meta)
}

func resourceTencentCloudCynosdbUpgradeProxyVersionRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_upgrade_proxy_version.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	clusterId := idSplit[0]

	upgradeProxyGroup, err := service.DescribeCynosdbUpgradeProxyVersionById(ctx, clusterId)
	if err != nil {
		return err
	}

	if upgradeProxyGroup == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CynosdbUpgradeProxyVersion` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if upgradeProxyGroup.ClusterId != nil {
		_ = d.Set("cluster_id", upgradeProxyGroup.ClusterId)
	}

	if upgradeProxyGroup.CurrentProxyVersion != nil {
		_ = d.Set("dst_proxy_version", upgradeProxyGroup.CurrentProxyVersion)
	}

	return nil
}

func resourceTencentCloudCynosdbUpgradeProxyVersionUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_upgrade_proxy_version.update")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}
		request = cynosdb.NewUpgradeProxyVersionRequest()
		flowId  int64
	)

	if d.HasChange("dst_proxy_version") {
		idSplit := strings.Split(d.Id(), FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", idSplit)
		}
		clusterId := idSplit[0]
		srcProxyVersion := idSplit[1]

		request.ClusterId = &clusterId
		request.SrcProxyVersion = &srcProxyVersion
		request.IsInMaintainPeriod = helper.String("no")

		if v, ok := d.GetOk("dst_proxy_version"); ok {
			request.DstProxyVersion = helper.String(v.(string))
		}

		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseCynosdbClient().UpgradeProxyVersion(request)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			flowId = *result.Response.FlowId
			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s update cynosdb upgradeProxyVersion failed, reason:%+v", logId, err)
			return err
		}

		err = resource.Retry(6*readRetryTimeout, func() *resource.RetryError {
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
				return resource.RetryableError(fmt.Errorf("update cynosdb upgradeProxyVersion is processing"))
			}
		})

		if err != nil {
			log.Printf("[CRITAL]%s update cynosdb upgradeProxyVersion fail, reason:%s\n", logId, err.Error())
			return err
		}
	}

	return resourceTencentCloudCynosdbUpgradeProxyVersionRead(d, meta)
}

func resourceTencentCloudCynosdbUpgradeProxyVersionDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_upgrade_proxy_version.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
