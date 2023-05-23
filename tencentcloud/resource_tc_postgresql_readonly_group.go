/*
Use this resource to create postgresql readonly group.

Example Usage

```hcl
resource "tencentcloud_postgresql_readonly_group" "group" {
  master_db_instance_id = "postgres-gzg9jb2n"
  name = "world"
  project_id = 0
  vpc_id = "vpc-86v957zb"
  subnet_id = "subnet-enm92y0m"
  replay_lag_eliminate = 1
  replay_latency_eliminate =  1
  max_replay_lag = 100
  max_replay_latency = 512
  min_delay_eliminate_reserve = 1
#  security_groups_ids = []
}
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	postgresql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudPostgresqlReadonlyGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPostgresqlReadOnlyGroupCreate,
		Read:   resourceTencentCloudPostgresqlReadOnlyGroupRead,
		Update: resourceTencentCloudPostgresqlReadOnlyGroupUpdate,
		Delete: resourceTencentCLoudPostgresqlReadOnlyGroupDelete,
		//Importer: &schema.ResourceImporter{
		//	State: schema.ImportStatePassthrough,
		//},

		Schema: map[string]*schema.Schema{
			"master_db_instance_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "Primary instance ID.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "RO group name.",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Project ID.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "VPC ID.",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "VPC subnet ID.",
			},
			"replay_lag_eliminate": {
				Type:     schema.TypeInt,
				Required: true,
				Description: "Whether to remove a read-only replica from an RO group if the delay between the read-only replica " +
					"and the primary instance exceeds the threshold. Valid values: 0 (no), 1 (yes).",
			},
			"replay_latency_eliminate": {
				Type:     schema.TypeInt,
				Required: true,
				Description: "Whether to remove a read-only replica from an RO group if the sync log size difference between " +
					"the read-only replica and the primary instance exceeds the threshold. Valid values: 0 (no), 1 (yes).",
			},
			"max_replay_lag": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Delay threshold in ms.",
			},
			"max_replay_latency": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Delayed log size threshold in MB.",
			},
			"min_delay_eliminate_reserve": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The minimum number of read-only replicas that must be retained in an RO group.",
			},
			"security_groups_ids": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "ID of security group. If both vpc_id and subnet_id are not set, this argument should not be set either.",
			},
			// Computed values
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Create time of the postgresql instance.",
			},
		},
	}
}

func resourceTencentCloudPostgresqlReadOnlyGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgresql_readonly_group.create")()

	logId := getLogId(contextNil)

	var (
		request  = postgresql.NewCreateReadOnlyGroupRequest()
		response *postgresql.CreateReadOnlyGroupResponse
	)
	if v, ok := d.GetOk("master_db_instance_id"); ok {
		request.MasterDBInstanceId = helper.String(v.(string))
	}
	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}
	if v, ok := d.GetOk("project_id"); ok {
		request.ProjectId = helper.IntUint64(v.(int))
	}
	if v, ok := d.GetOk("vpc_id"); ok {
		request.VpcId = helper.String(v.(string))
	}
	if v, ok := d.GetOk("subnet_id"); ok {
		request.SubnetId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("replay_lag_eliminate"); ok {
		request.ReplayLagEliminate = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("replay_latency_eliminate"); ok {
		request.ReplayLatencyEliminate = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("max_replay_lag"); ok {
		request.MaxReplayLag = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("max_replay_latency"); ok {
		request.MaxReplayLatency = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("min_delay_eliminate_reserve"); ok {
		request.MinDelayEliminateReserve = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("security_groups_ids"); ok {
		securityGroupsIds := v.(*schema.Set).List()
		request.SecurityGroupIds = make([]*string, 0, len(securityGroupsIds))
		for _, item := range securityGroupsIds {
			request.SecurityGroupIds = append(request.SecurityGroupIds, helper.String(item.(string)))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UsePostgresqlClient().CreateReadOnlyGroup(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		return err
	}
	instanceId := *response.Response.ReadOnlyGroupId
	d.SetId(instanceId)

	//return resourceTencentCloudPostgresqlReadOnlyGroupRead(d, meta)
	return nil
}

func resourceTencentCloudPostgresqlReadOnlyGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgresql_readonly_group.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	postgresqlService := PostgresqlService{client: meta.(*TencentCloudClient).apiV3Conn}
	_, err := postgresqlService.DescribePostgresqlReadOnlyGroupById(ctx, d.Id())
	if err != nil {
		return err
	}

	return nil
}

func resourceTencentCloudPostgresqlReadOnlyGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgresql_readonly_group.update")()

	logId := getLogId(contextNil)
	request := postgresql.NewModifyReadOnlyGroupConfigRequest()

	request.ReadOnlyGroupId = helper.String(d.Id())

	if d.HasChange("name") {
		request.ReadOnlyGroupName = helper.String(d.Get("name").(string))
	}
	if d.HasChange("replay_lag_eliminate") {
		request.ReplayLagEliminate = helper.IntUint64(d.Get("replay_lag_eliminate").(int))
	}
	if d.HasChange("replay_latency_eliminate") {
		request.ReplayLatencyEliminate = helper.IntUint64(d.Get("replay_latency_eliminate").(int))
	}
	if d.HasChange("max_replay_lag") {
		request.MaxReplayLag = helper.IntUint64(d.Get("max_replay_lag").(int))
	}
	if d.HasChange("max_replay_latency") {
		request.MaxReplayLatency = helper.IntUint64(d.Get("max_replay_latency").(int))
	}
	if d.HasChange("min_delay_eliminate_reserve") {
		request.MinDelayEliminateReserve = helper.IntUint64(d.Get("min_delay_eliminate_reserve").(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UsePostgresqlClient().ModifyReadOnlyGroupConfig(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

func resourceTencentCLoudPostgresqlReadOnlyGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgresql_readonly_group.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	groupId := d.Id()
	postgresqlService := PostgresqlService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		e := postgresqlService.DeletePostgresqlReadOnlyGroupById(ctx, groupId)
		if e != nil {
			return retryError(e)
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
