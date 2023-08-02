/*
Provides a resource to create a cynosdb cluster_version

Example Usage

```hcl
resource "tencentcloud_cynosdb_cluster_version" "cluster_version" {
  cluster_id = "xxx"
  cynos_version = "2.0.0"
}
```

Import

cynosdb cluster_version can be imported using the id, e.g.

```
terraform import tencentcloud_cynosdb_cluster_version.cluster_version cluster_version_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudCynosdbClusterVersion() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCynosdbClusterVersionCreate,
		Read:   resourceTencentCloudCynosdbClusterVersionRead,
		Delete: resourceTencentCloudCynosdbClusterVersionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Cluster ID.",
			},

			"cynos_version": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Kernel version.",
			},
		},
	}
}

func resourceTencentCloudCynosdbClusterVersionCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_cluster_version.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var (
		request   = cynosdb.NewUpgradeClusterVersionRequest()
		response  = cynosdb.NewUpgradeClusterVersionResponse()
		clusterId string
	)
	if v, ok := d.GetOk("cluster_id"); ok {
		clusterId = v.(string)
		request.ClusterId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cynos_version"); ok {
		request.CynosVersion = helper.String(v.(string))
	}

	request.UpgradeType = helper.String("upgradeImmediate")

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCynosdbClient().UpgradeClusterVersion(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate cynosdb clusterVersion failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(clusterId)

	flowId := *response.Response.FlowId
	service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}
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
			return resource.RetryableError(fmt.Errorf("Open cynosdb clusterVersion is processing"))
		}
	})

	if err != nil {
		log.Printf("[CRITAL]%s Open cynosdb clusterVersion fail, reason:%s\n", logId, err.Error())
		return err
	}

	return resourceTencentCloudCynosdbClusterVersionRead(d, meta)
}

func resourceTencentCloudCynosdbClusterVersionRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_cluster_version.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudCynosdbClusterVersionDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_cluster_version.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
