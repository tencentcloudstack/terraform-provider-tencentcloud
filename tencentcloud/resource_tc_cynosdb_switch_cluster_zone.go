/*
Provides a resource to create a cynosdb switch_cluster_zone

Example Usage

```hcl
resource "tencentcloud_cynosdb_switch_cluster_zone" "switch_cluster_zone" {
  cluster_id = "cynosdbmysql-507j6phr"
  zone = "ap-guangzhou-6"
}
```

Import

cynosdb switch_cluster_zone can be imported using the id, e.g.

```
terraform import tencentcloud_cynosdb_switch_cluster_zone.switch_cluster_zone switch_cluster_zone_id
```
*/
package tencentcloud

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudCynosdbSwitchClusterZone() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCynosdbSwitchClusterZoneCreate,
		Read:   resourceTencentCloudCynosdbSwitchClusterZoneRead,
		Update: resourceTencentCloudCynosdbSwitchClusterZoneUpdate,
		Delete: resourceTencentCloudCynosdbSwitchClusterZoneDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Cluster Id.",
			},

			"zone": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Availability zone to switch to.",
			},
		},
	}
}

func resourceTencentCloudCynosdbSwitchClusterZoneCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_switch_cluster_zone.create")()
	defer inconsistentCheck(d, meta)()

	var clusterId string
	if v, ok := d.GetOk("cluster_id"); ok {
		clusterId = v.(string)
	}

	d.SetId(clusterId)

	return resourceTencentCloudCynosdbSwitchClusterZoneUpdate(d, meta)
}

func resourceTencentCloudCynosdbSwitchClusterZoneRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_switch_cluster_zone.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	clusterId := d.Id()
	clusterSlaveZone, err := service.DescribeCynosdbClusterSlaveZoneById(ctx, clusterId)
	if err != nil {
		return err
	}

	if clusterSlaveZone == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CynosdbSwitchClusterZone` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("cluster_id", clusterId)

	if clusterSlaveZone.MasterZone != nil {
		_ = d.Set("zone", clusterSlaveZone.MasterZone)
	}

	return nil
}

func resourceTencentCloudCynosdbSwitchClusterZoneUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_switch_cluster_zone.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	request := cynosdb.NewSwitchClusterZoneRequest()

	clusterId := d.Id()
	service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}
	clusterSlaveZone, err := service.DescribeCynosdbClusterSlaveZoneById(ctx, clusterId)
	if err != nil {
		return err
	}
	if clusterSlaveZone == nil {
		log.Printf("[WARN]%s resource `CynosdbSwitchClusterZone` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	request.OldZone = clusterSlaveZone.MasterZone
	request.ClusterId = &clusterId
	request.IsInMaintainPeriod = helper.String("no")

	if v, ok := d.GetOk("zone"); ok {
		request.NewZone = helper.String(v.(string))
	}

	var flowId *int64
	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCynosdbClient().SwitchClusterZone(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		flowId = result.Response.FlowId
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update cynosdb switchClusterZone failed, reason:%+v", logId, err)
		return err
	}

	conf := BuildStateChangeConf([]string{}, []string{CYNOSDB_FLOW_STATUS_SUCCESSFUL}, 10*readRetryTimeout, time.Second, service.CynosdbClusterSlaveZoneStateRefreshFunc(*flowId, []string{}))
	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudCynosdbSwitchClusterZoneRead(d, meta)
}

func resourceTencentCloudCynosdbSwitchClusterZoneDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_switch_cluster_zone.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
