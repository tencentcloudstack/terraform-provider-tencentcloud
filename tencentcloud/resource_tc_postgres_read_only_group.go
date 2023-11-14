/*
Provides a resource to create a postgres read_only_group

Example Usage

```hcl
resource "tencentcloud_postgres_read_only_group" "read_only_group" {
  master_d_b_instance_id = "postgres-xxxx"
  name = "test-rg"
  project_id = 0
  vpc_id = "vpc-e0tfm161"
  subnet_id = "subnet-443a3lv6"
  replay_lag_eliminate = 0
  replay_latency_eliminate = 0
  max_replay_lag = 5000
  max_replay_latency = 32
  min_delay_eliminate_reserve = 1
  security_group_ids =
}
```

Import

postgres read_only_group can be imported using the id, e.g.

```
terraform import tencentcloud_postgres_read_only_group.read_only_group read_only_group_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	postgres "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"time"
)

func resourceTencentCloudPostgresReadOnlyGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPostgresReadOnlyGroupCreate,
		Read:   resourceTencentCloudPostgresReadOnlyGroupRead,
		Update: resourceTencentCloudPostgresReadOnlyGroupUpdate,
		Delete: resourceTencentCloudPostgresReadOnlyGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"master_d_b_instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Primary instance ID.",
			},

			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "RO group name.",
			},

			"project_id": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Project ID.",
			},

			"vpc_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Vpc ID.",
			},

			"subnet_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Subnet ID.",
			},

			"replay_lag_eliminate": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Whether to remove a read-only replica from an RO group if the delay between the read-only replica and the primary instance exceeds the threshold. Valid values:0 (no), 1 (yes).",
			},

			"replay_latency_eliminate": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Whether to remove a read-only replica from an RO group if the sync log size difference between the read-only replica and the primary instance exceeds the threshold. Valid values:0 (no), 1 (yes).",
			},

			"max_replay_lag": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Delay threshold in ms.",
			},

			"max_replay_latency": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Delayed log size threshold in MB.",
			},

			"min_delay_eliminate_reserve": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "The minimum number of read-only replicas that must be retained in an RO group.",
			},

			"security_group_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Security group ID.",
			},
		},
	}
}

func resourceTencentCloudPostgresReadOnlyGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgres_read_only_group.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request         = postgres.NewCreateReadOnlyGroupRequest()
		response        = postgres.NewCreateReadOnlyGroupResponse()
		readOnlyGroupId string
	)
	if v, ok := d.GetOk("master_d_b_instance_id"); ok {
		request.MasterDBInstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("project_id"); ok {
		request.ProjectId = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		request.VpcId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("subnet_id"); ok {
		request.SubnetId = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("replay_lag_eliminate"); ok {
		request.ReplayLagEliminate = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("replay_latency_eliminate"); ok {
		request.ReplayLatencyEliminate = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("max_replay_lag"); ok {
		request.MaxReplayLag = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("max_replay_latency"); ok {
		request.MaxReplayLatency = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("min_delay_eliminate_reserve"); ok {
		request.MinDelayEliminateReserve = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("security_group_ids"); ok {
		securityGroupIdsSet := v.(*schema.Set).List()
		for i := range securityGroupIdsSet {
			securityGroupIds := securityGroupIdsSet[i].(string)
			request.SecurityGroupIds = append(request.SecurityGroupIds, &securityGroupIds)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UsePostgresClient().CreateReadOnlyGroup(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create postgres ReadOnlyGroup failed, reason:%+v", logId, err)
		return err
	}

	readOnlyGroupId = *response.Response.ReadOnlyGroupId
	d.SetId(readOnlyGroupId)

	service := PostgresService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"ok"}, 100*readRetryTimeout, time.Second, service.PostgresReadOnlyGroupStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudPostgresReadOnlyGroupRead(d, meta)
}

func resourceTencentCloudPostgresReadOnlyGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgres_read_only_group.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := PostgresService{client: meta.(*TencentCloudClient).apiV3Conn}

	readOnlyGroupId := d.Id()

	ReadOnlyGroup, err := service.DescribePostgresReadOnlyGroupById(ctx, readOnlyGroupId)
	if err != nil {
		return err
	}

	if ReadOnlyGroup == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `PostgresReadOnlyGroup` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if ReadOnlyGroup.MasterDBInstanceId != nil {
		_ = d.Set("master_d_b_instance_id", ReadOnlyGroup.MasterDBInstanceId)
	}

	if ReadOnlyGroup.Name != nil {
		_ = d.Set("name", ReadOnlyGroup.Name)
	}

	if ReadOnlyGroup.ProjectId != nil {
		_ = d.Set("project_id", ReadOnlyGroup.ProjectId)
	}

	if ReadOnlyGroup.VpcId != nil {
		_ = d.Set("vpc_id", ReadOnlyGroup.VpcId)
	}

	if ReadOnlyGroup.SubnetId != nil {
		_ = d.Set("subnet_id", ReadOnlyGroup.SubnetId)
	}

	if ReadOnlyGroup.ReplayLagEliminate != nil {
		_ = d.Set("replay_lag_eliminate", ReadOnlyGroup.ReplayLagEliminate)
	}

	if ReadOnlyGroup.ReplayLatencyEliminate != nil {
		_ = d.Set("replay_latency_eliminate", ReadOnlyGroup.ReplayLatencyEliminate)
	}

	if ReadOnlyGroup.MaxReplayLag != nil {
		_ = d.Set("max_replay_lag", ReadOnlyGroup.MaxReplayLag)
	}

	if ReadOnlyGroup.MaxReplayLatency != nil {
		_ = d.Set("max_replay_latency", ReadOnlyGroup.MaxReplayLatency)
	}

	if ReadOnlyGroup.MinDelayEliminateReserve != nil {
		_ = d.Set("min_delay_eliminate_reserve", ReadOnlyGroup.MinDelayEliminateReserve)
	}

	if ReadOnlyGroup.SecurityGroupIds != nil {
		_ = d.Set("security_group_ids", ReadOnlyGroup.SecurityGroupIds)
	}

	return nil
}

func resourceTencentCloudPostgresReadOnlyGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgres_read_only_group.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := postgres.NewModifyReadOnlyGroupConfigRequest()

	readOnlyGroupId := d.Id()

	request.ReadOnlyGroupId = &readOnlyGroupId

	immutableArgs := []string{"master_d_b_instance_id", "name", "project_id", "vpc_id", "subnet_id", "replay_lag_eliminate", "replay_latency_eliminate", "max_replay_lag", "max_replay_latency", "min_delay_eliminate_reserve", "security_group_ids"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("replay_lag_eliminate") {
		if v, ok := d.GetOkExists("replay_lag_eliminate"); ok {
			request.ReplayLagEliminate = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("replay_latency_eliminate") {
		if v, ok := d.GetOkExists("replay_latency_eliminate"); ok {
			request.ReplayLatencyEliminate = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("max_replay_lag") {
		if v, ok := d.GetOkExists("max_replay_lag"); ok {
			request.MaxReplayLag = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("max_replay_latency") {
		if v, ok := d.GetOkExists("max_replay_latency"); ok {
			request.MaxReplayLatency = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("min_delay_eliminate_reserve") {
		if v, ok := d.GetOkExists("min_delay_eliminate_reserve"); ok {
			request.MinDelayEliminateReserve = helper.IntUint64(v.(int))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UsePostgresClient().ModifyReadOnlyGroupConfig(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update postgres ReadOnlyGroup failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudPostgresReadOnlyGroupRead(d, meta)
}

func resourceTencentCloudPostgresReadOnlyGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgres_read_only_group.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := PostgresService{client: meta.(*TencentCloudClient).apiV3Conn}
	readOnlyGroupId := d.Id()

	if err := service.DeletePostgresReadOnlyGroupById(ctx, readOnlyGroupId); err != nil {
		return err
	}

	return nil
}
