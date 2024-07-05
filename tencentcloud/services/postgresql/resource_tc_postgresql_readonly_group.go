package postgresql

import (
	"context"
	"fmt"
	"log"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	postgresql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudPostgresqlReadonlyGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPostgresqlReadOnlyGroupCreate,
		Read:   resourceTencentCloudPostgresqlReadOnlyGroupRead,
		Update: resourceTencentCloudPostgresqlReadOnlyGroupUpdate,
		Delete: resourceTencentCLoudPostgresqlReadOnlyGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
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
			"net_info_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of db instance net info.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Ip address of the net info.",
						},
						"port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Port of the net info.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudPostgresqlReadOnlyGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_readonly_group.create")()

	var (
		logId              = tccommon.GetLogId(tccommon.ContextNil)
		ctx                = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		postgresqlService  = PostgresqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		request            = postgresql.NewCreateReadOnlyGroupRequest()
		response           *postgresql.CreateReadOnlyGroupResponse
		masterDbInstanceId string
	)

	if v, ok := d.GetOk("master_db_instance_id"); ok {
		request.MasterDBInstanceId = helper.String(v.(string))
		masterDbInstanceId = v.(string)
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

	if v, ok := d.GetOk("security_groups_ids"); ok {
		securityGroupsIds := v.(*schema.Set).List()
		request.SecurityGroupIds = make([]*string, 0, len(securityGroupsIds))
		for _, item := range securityGroupsIds {
			request.SecurityGroupIds = append(request.SecurityGroupIds, helper.String(item.(string)))
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePostgresqlClient().CreateReadOnlyGroup(request)
		if e != nil {
			return tccommon.RetryError(e)
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

	readOnlyGroupId := *response.Response.ReadOnlyGroupId
	d.SetId(readOnlyGroupId)

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		groups, e := postgresqlService.DescribePostgresqlReadOnlyGroupById(ctx, masterDbInstanceId)
		if e != nil {
			return tccommon.RetryError(e)
		}

		var status string
		for _, gg := range groups {
			if *gg.ReadOnlyGroupId == readOnlyGroupId {
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

	return resourceTencentCloudPostgresqlReadOnlyGroupRead(d, meta)
}

func resourceTencentCloudPostgresqlReadOnlyGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_readonly_group.read")()

	var (
		logId             = tccommon.GetLogId(tccommon.ContextNil)
		ctx               = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		postgresqlService = PostgresqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		readOnlyGroupId   = d.Id()
	)

	readOnlyGroupInfo, err := postgresqlService.DescribePostgresqlReadonlyGroupsById(ctx, readOnlyGroupId)
	if err != nil {
		return err
	}

	if readOnlyGroupInfo.MasterDBInstanceId != nil {
		_ = d.Set("master_db_instance_id", readOnlyGroupInfo.MasterDBInstanceId)
	}

	if readOnlyGroupInfo.ReadOnlyGroupName != nil {
		_ = d.Set("name", readOnlyGroupInfo.ReadOnlyGroupName)
	}

	if readOnlyGroupInfo.ProjectId != nil {
		_ = d.Set("project_id", readOnlyGroupInfo.ProjectId)
	}

	if readOnlyGroupInfo.VpcId != nil {
		_ = d.Set("vpc_id", readOnlyGroupInfo.VpcId)
	}

	if readOnlyGroupInfo.SubnetId != nil {
		_ = d.Set("subnet_id", readOnlyGroupInfo.SubnetId)
	}

	if readOnlyGroupInfo.ReplayLagEliminate != nil {
		_ = d.Set("replay_lag_eliminate", readOnlyGroupInfo.ReplayLagEliminate)
	}

	if readOnlyGroupInfo.ReplayLatencyEliminate != nil {
		_ = d.Set("replay_latency_eliminate", readOnlyGroupInfo.ReplayLatencyEliminate)
	}

	if readOnlyGroupInfo.MaxReplayLag != nil {
		_ = d.Set("max_replay_lag", readOnlyGroupInfo.MaxReplayLag)
	}

	if readOnlyGroupInfo.MaxReplayLatency != nil {
		_ = d.Set("max_replay_latency", readOnlyGroupInfo.MaxReplayLatency)
	}

	if readOnlyGroupInfo.MinDelayEliminateReserve != nil {
		_ = d.Set("min_delay_eliminate_reserve", readOnlyGroupInfo.MinDelayEliminateReserve)
	}

	if readOnlyGroupInfo.DBInstanceNetInfo != nil {
		netInfoList := []interface{}{}
		for _, netInfo := range readOnlyGroupInfo.DBInstanceNetInfo {
			netInfoMap := map[string]interface{}{}
			if netInfo.Ip != nil {
				netInfoMap["ip"] = *netInfo.Ip
			}

			if netInfo.Port != nil {
				netInfoMap["port"] = helper.UInt64Int64(*netInfo.Port)
			}

			netInfoList = append(netInfoList, netInfoMap)
		}

		_ = d.Set("net_info_list", netInfoList)
	}

	// security groups
	sg, err := postgresqlService.DescribeDBInstanceSecurityGroupsByGroupId(ctx, readOnlyGroupId)
	if err != nil {
		return err
	}

	if len(sg) > 0 {
		_ = d.Set("security_groups_ids", sg)
	}

	return nil
}

func resourceTencentCloudPostgresqlReadOnlyGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_readonly_group.update")()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = PostgresqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	// update vpc and subnet
	if d.HasChange("vpc_id") || d.HasChange("subnet_id") {
		var (
			vpcOld             string
			vpcNew             string
			subnetOld          string
			subnetNew          string
			vipOld             string
			vipNew             string
			masterDbInstanceId string
		)

		if v, ok := d.GetOk("master_db_instance_id"); ok {
			masterDbInstanceId = v.(string)
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

		// get the old ip before creating
		netInfos, err := service.DescribePostgresqlReadonlyGroupNetInfosById(ctx, masterDbInstanceId, d.Id())
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

		err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePostgresqlClient().CreateReadOnlyGroupNetworkAccess(request)
			if e != nil {
				return tccommon.RetryError(e)
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
		conf := tccommon.BuildStateChangeConf([]string{}, []string{"opened"}, 3*tccommon.ReadRetryTimeout, time.Second, service.PostgresqlReadonlyGroupNetworkAccessStateRefreshFunc(masterDbInstanceId, d.Id(), vpcNew, subnetNew, vipOld, "", []string{}))

		if object, e := conf.WaitForState(); e != nil {
			return e
		} else {
			// find the vip assiged by system
			ret := object.(*postgresql.DBInstanceNetInfo)
			vipNew = *ret.Ip
		}

		log.Printf("[DEBUG]%s resourceTencentCloudPostgresqlReadOnlyGroupUpdate, msaterDbInstanceId:[%s], roGroupId:[%s], vpcOld:[%s], vpcNew:[%s], subnetOld:[%s], subnetNew:[%s], vipOld:[%s], vipNew:[%s]\n",
			logId, masterDbInstanceId, d.Id(), vpcOld, vpcNew, subnetOld, subnetNew, vipOld, vipNew)

		// wait unit network changing operation of ro group done
		conf = tccommon.BuildStateChangeConf([]string{}, []string{"ok"}, 3*tccommon.ReadRetryTimeout, time.Second, service.PostgresqlReadonlyGroupStateRefreshFunc(masterDbInstanceId, d.Id(), []string{}))
		if _, e := conf.WaitForState(); e != nil {
			return e
		}

		// delete the old one
		if err := service.DeletePostgresqlReadonlyGroupNetworkAccessById(ctx, d.Id(), vpcOld, subnetOld, vipOld); err != nil {
			return err
		}

		// wait for old network removed
		conf = tccommon.BuildStateChangeConf([]string{}, []string{"closed"}, 3*tccommon.ReadRetryTimeout, time.Second, service.PostgresqlReadonlyGroupNetworkAccessStateRefreshFunc(masterDbInstanceId, d.Id(), vpcOld, subnetOld, vipNew, vipOld, []string{}))
		if _, e := conf.WaitForState(); e != nil {
			return e
		}

		// wait unit network changing operation of ro group done
		conf = tccommon.BuildStateChangeConf([]string{}, []string{"ok"}, 3*tccommon.ReadRetryTimeout, time.Second, service.PostgresqlReadonlyGroupStateRefreshFunc(masterDbInstanceId, d.Id(), []string{}))
		if _, e := conf.WaitForState(); e != nil {
			return e
		}

		// refresh the private ip with new one
	}

	// required attributes
	request := postgresql.NewModifyReadOnlyGroupConfigRequest()
	request.ReadOnlyGroupId = helper.String(d.Id())
	request.ReadOnlyGroupName = helper.String(d.Get("name").(string))
	request.ReplayLagEliminate = helper.IntUint64(d.Get("replay_lag_eliminate").(int))
	request.ReplayLatencyEliminate = helper.IntUint64(d.Get("replay_latency_eliminate").(int))
	request.MaxReplayLag = helper.IntUint64(d.Get("max_replay_lag").(int))
	request.MaxReplayLatency = helper.IntUint64(d.Get("max_replay_latency").(int))
	request.MinDelayEliminateReserve = helper.IntUint64(d.Get("min_delay_eliminate_reserve").(int))

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePostgresqlClient().ModifyReadOnlyGroupConfig(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		return err
	}

	if d.HasChange("security_groups_ids") {
		ids := d.Get("security_groups_ids").(*schema.Set).List()
		var sgIds []*string
		for _, id := range ids {
			sgIds = append(sgIds, helper.String(id.(string)))
		}

		err = service.ModifyDBInstanceSecurityGroupsByGroupId(ctx, d.Id(), sgIds)
		if err != nil {
			return err
		}
	}

	return resourceTencentCloudPostgresqlReadOnlyGroupRead(d, meta)
}

func resourceTencentCLoudPostgresqlReadOnlyGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_readonly_group.delete")()

	var (
		logId             = tccommon.GetLogId(tccommon.ContextNil)
		ctx               = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		postgresqlService = PostgresqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		groupId           = d.Id()
	)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		e := postgresqlService.DeletePostgresqlReadOnlyGroupById(ctx, groupId)
		if e != nil {
			return tccommon.RetryError(e)
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
