package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudCynosdbClusterSlaveZone() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCynosdbClusterSlaveZoneCreate,
		Read:   resourceTencentCloudCynosdbClusterSlaveZoneRead,
		Update: resourceTencentCloudCynosdbClusterSlaveZoneUpdate,
		Delete: resourceTencentCloudCynosdbClusterSlaveZoneDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(300 * time.Second),
			Update: schema.DefaultTimeout(300 * time.Second),
			Delete: schema.DefaultTimeout(300 * time.Second),
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The ID of cluster.",
			},

			"slave_zone": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Slave zone.",
			},
		},
	}
}

func resourceTencentCloudCynosdbClusterSlaveZoneCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_cluster_slave_zone.create")()
	defer inconsistentCheck(d, meta)()

	timeout := d.Timeout(schema.TimeoutCreate)

	logId := getLogId(contextNil)

	var (
		request   = cynosdb.NewAddClusterSlaveZoneRequest()
		response  = cynosdb.NewAddClusterSlaveZoneResponse()
		flowId    *int64
		clusterId string
		slaveZone string
	)
	if v, ok := d.GetOk("cluster_id"); ok {
		request.ClusterId = helper.String(v.(string))
		clusterId = v.(string)
	}

	if v, ok := d.GetOk("slave_zone"); ok {
		request.SlaveZone = helper.String(v.(string))
		slaveZone = v.(string)
	}

	err := resource.Retry(timeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCynosdbClient().AddClusterSlaveZone(request)
		if e != nil {
			if sdkErr, ok := e.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkErr.Code == "FailedOperation.OperationFailedError" {
					return resource.NonRetryableError(e)
				}
			}
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cynosdb clusterSlaveZone failed, reason:%+v", logId, err)
		return err
	}

	flowId = response.Response.FlowId

	if flowId == nil {
		return fmt.Errorf("delete [%s] failed, reason: FlowId is null.\n", d.Id())
	}

	service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}
	conf := BuildStateChangeConf([]string{}, []string{CYNOSDB_FLOW_STATUS_SUCCESSFUL}, timeout, 3*time.Second, service.CynosdbClusterSlaveZoneStateRefreshFunc(*flowId, []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	d.SetId(strings.Join([]string{clusterId, slaveZone}, FILED_SP))

	return resourceTencentCloudCynosdbClusterSlaveZoneRead(d, meta)
}

func resourceTencentCloudCynosdbClusterSlaveZoneRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_cluster_slave_zone.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	clusterId := idSplit[0]
	slaveZone := idSplit[1]

	clusterSlaveZone, err := service.DescribeCynosdbClusterSlaveZoneById(ctx, clusterId)
	if err != nil {
		return err
	}

	if clusterSlaveZone == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CynosdbClusterSlaveZone` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if clusterSlaveZone.ClusterId != nil {
		_ = d.Set("cluster_id", clusterSlaveZone.ClusterId)
	}

	if len(clusterSlaveZone.SlaveZones) > 0 {
		for _, zone := range clusterSlaveZone.SlaveZones {
			if *zone == slaveZone {
				_ = d.Set("slave_zone", zone)
				break
			}
		}
	}

	return nil
}

func resourceTencentCloudCynosdbClusterSlaveZoneUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_cluster_slave_zone.update")()
	defer inconsistentCheck(d, meta)()

	timeout := d.Timeout(schema.TimeoutUpdate)

	logId := getLogId(contextNil)

	request := cynosdb.NewModifyClusterSlaveZoneRequest()
	var newSlaveZone string

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	clusterId := idSplit[0]
	slaveZone := idSplit[1]

	immutableArgs := []string{"cluster_id"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("slave_zone") {
		if v, ok := d.GetOk("slave_zone"); ok {
			newSlaveZone = v.(string)
		}
	}

	request.ClusterId = &clusterId
	request.NewSlaveZone = helper.String(newSlaveZone)
	request.OldSlaveZone = helper.String(slaveZone)

	var flowId *int64
	err := resource.Retry(timeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCynosdbClient().ModifyClusterSlaveZone(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		flowId = result.Response.FlowId
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update cynosdb clusterSlaveZone failed, reason:%+v", logId, err)
		return err
	}

	if flowId == nil {
		return fmt.Errorf("delete [%s] failed, reason: FlowId is null.\n", d.Id())
	}

	service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{CYNOSDB_FLOW_STATUS_SUCCESSFUL}, timeout, time.Second, service.CynosdbClusterSlaveZoneStateRefreshFunc(*flowId, []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	// update the id
	d.SetId(strings.Join([]string{clusterId, newSlaveZone}, FILED_SP))

	return resourceTencentCloudCynosdbClusterSlaveZoneRead(d, meta)
}

func resourceTencentCloudCynosdbClusterSlaveZoneDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_cluster_slave_zone.delete")()
	defer inconsistentCheck(d, meta)()

	timeout := d.Timeout(schema.TimeoutDelete)

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	clusterId := idSplit[0]
	slaveZone := idSplit[1]

	var (
		flowId *int64
		err    error
	)

	if flowId, err = service.DeleteCynosdbClusterSlaveZoneById(ctx, clusterId, slaveZone); err != nil {
		return err
	}

	if flowId == nil {
		return fmt.Errorf("delete [%s] failed, reason: FlowId is null.\n", d.Id())
	}

	conf := BuildStateChangeConf([]string{}, []string{CYNOSDB_FLOW_STATUS_SUCCESSFUL}, timeout, time.Second, service.CynosdbClusterSlaveZoneStateRefreshFunc(*flowId, []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return nil
}
