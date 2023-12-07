package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"time"

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
		request            = postgresql.NewCreateReadOnlyGroupRequest()
		response           *postgresql.CreateReadOnlyGroupResponse
		msaterDbInstanceId string
	)
	if v, ok := d.GetOk("master_db_instance_id"); ok {
		request.MasterDBInstanceId = helper.String(v.(string))
		msaterDbInstanceId = v.(string)
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

	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	postgresqlService := PostgresqlService{client: meta.(*TencentCloudClient).apiV3Conn}

	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		groups, e := postgresqlService.DescribePostgresqlReadOnlyGroupById(ctx, msaterDbInstanceId)
		if e != nil {
			return retryError(e)
		}

		var status string
		for _, gg := range groups {
			if *gg.ReadOnlyGroupId == instanceId {
				status = *gg.Status
				if status == "ok" {
					return nil
				}
			}
		}
		return resource.RetryableError(fmt.Errorf("waiting status[%s] to running, retry... ", status))
	})
	if err != nil {
		return err
	}

	//return resourceTencentCloudPostgresqlReadOnlyGroupRead(d, meta)
	return nil
}

func resourceTencentCloudPostgresqlReadOnlyGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgresql_readonly_group.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	// for now, the id should be the master db instance id, cause the describe api only support this kind of filter.
	var id string
	if v, ok := d.GetOk("master_db_instance_id"); ok {
		id = v.(string)
	}

	postgresqlService := PostgresqlService{client: meta.(*TencentCloudClient).apiV3Conn}
	_, err := postgresqlService.DescribePostgresqlReadOnlyGroupById(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func resourceTencentCloudPostgresqlReadOnlyGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgresql_readonly_group.update")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	request := postgresql.NewModifyReadOnlyGroupConfigRequest()

	request.ReadOnlyGroupId = helper.String(d.Id())

	// update vpc and subnet
	if d.HasChange("vpc_id") || d.HasChange("subnet_id") {
		var (
			vpcOld             string
			vpcNew             string
			subnetOld          string
			subnetNew          string
			vipOld             string
			vipNew             string
			msaterDbInstanceId string
		)

		if v, ok := d.GetOk("master_db_instance_id"); ok {
			msaterDbInstanceId = v.(string)
		}

		old, new := d.GetChange("vpc_id")
		if old != nil {
			vpcOld = old.(string)
		}
		if new != nil {
			vpcNew = new.(string)
		}

		old, new = d.GetChange("subnet_id")
		if old != nil {
			subnetOld = old.(string)
		}
		if new != nil {
			subnetNew = new.(string)
		}

		service := PostgresqlService{client: meta.(*TencentCloudClient).apiV3Conn}
		// get the old ip before creating
		netInfos, err := service.DescribePostgresqlReadonlyGroupNetInfosById(ctx, msaterDbInstanceId, d.Id())
		if err != nil {
			return err
		}

		var oldNetInfo *postgresql.DBInstanceNetInfo
		for _, info := range netInfos {
			if *info.NetType == "private" {
				if *info.VpcId == vpcOld && *info.SubnetId == subnetOld {
					oldNetInfo = info
					break
				}
			}
		}

		if oldNetInfo != nil {
			vipOld = *oldNetInfo.Ip
		}

		// Create new network first, then delete the old one
		request := postgresql.NewCreateReadOnlyGroupNetworkAccessRequest()
		request.ReadOnlyGroupId = helper.String(d.Id())
		request.VpcId = helper.String(vpcNew)
		request.SubnetId = helper.String(subnetNew)
		// ip assigned by system
		request.IsAssignVip = helper.Bool(false)

		err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UsePostgresqlClient().CreateReadOnlyGroupNetworkAccess(request)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s create postgresql ReadOnlyGroup NetworkAccess failed, reason:%+v", logId, err)
			return err
		}

		// wait for new network enabled
		conf := BuildStateChangeConf([]string{}, []string{"opened"}, 3*readRetryTimeout, time.Second, service.PostgresqlReadonlyGroupNetworkAccessStateRefreshFunc(msaterDbInstanceId, d.Id(), vpcNew, subnetNew, vipOld, "", []string{}))

		if object, e := conf.WaitForState(); e != nil {
			return e
		} else {
			// find the vip assiged by system
			ret := object.(*postgresql.DBInstanceNetInfo)
			vipNew = *ret.Ip
		}

		log.Printf("[DEBUG]%s resourceTencentCloudPostgresqlReadOnlyGroupUpdate, msaterDbInstanceId:[%s], roGroupId:[%s], vpcOld:[%s], vpcNew:[%s], subnetOld:[%s], subnetNew:[%s], vipOld:[%s], vipNew:[%s]\n",
			logId, msaterDbInstanceId, d.Id(), vpcOld, vpcNew, subnetOld, subnetNew, vipOld, vipNew)

		// wait unit network changing operation of ro group done
		conf = BuildStateChangeConf([]string{}, []string{"ok"}, 3*readRetryTimeout, time.Second, service.PostgresqlReadonlyGroupStateRefreshFunc(msaterDbInstanceId, d.Id(), []string{}))
		if _, e := conf.WaitForState(); e != nil {
			return e
		}

		// delete the old one
		if err := service.DeletePostgresqlReadonlyGroupNetworkAccessById(ctx, d.Id(), vpcOld, subnetOld, vipOld); err != nil {
			return err
		}

		// wait for old network removed
		conf = BuildStateChangeConf([]string{}, []string{"closed"}, 3*readRetryTimeout, time.Second, service.PostgresqlReadonlyGroupNetworkAccessStateRefreshFunc(msaterDbInstanceId, d.Id(), vpcOld, subnetOld, vipNew, vipOld, []string{}))
		if _, e := conf.WaitForState(); e != nil {
			return e
		}

		// wait unit network changing operation of ro group done
		conf = BuildStateChangeConf([]string{}, []string{"ok"}, 3*readRetryTimeout, time.Second, service.PostgresqlReadonlyGroupStateRefreshFunc(msaterDbInstanceId, d.Id(), []string{}))
		if _, e := conf.WaitForState(); e != nil {
			return e
		}

		// refresh the private ip with new one
	}

	// required attributes
	request.ReadOnlyGroupName = helper.String(d.Get("name").(string))
	request.ReplayLagEliminate = helper.IntUint64(d.Get("replay_lag_eliminate").(int))
	request.ReplayLatencyEliminate = helper.IntUint64(d.Get("replay_latency_eliminate").(int))
	request.MaxReplayLag = helper.IntUint64(d.Get("max_replay_lag").(int))
	request.MaxReplayLatency = helper.IntUint64(d.Get("max_replay_latency").(int))
	request.MinDelayEliminateReserve = helper.IntUint64(d.Get("min_delay_eliminate_reserve").(int))

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
