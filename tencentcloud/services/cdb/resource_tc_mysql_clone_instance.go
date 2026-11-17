package cdb

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	sdkError "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudMysqlCloneInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMysqlCloneInstanceCreate,
		Read:   resourceTencentCloudMysqlCloneInstanceRead,
		Update: resourceTencentCloudMysqlCloneInstanceUpdate,
		Delete: resourceTencentCloudMysqlCloneInstanceDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Source MySQL (CDB) instance ID to clone from.",
			},

			"specified_rollback_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Rollback time (yyyy-mm-dd hh:mm:ss). Mutually exclusive with `specified_backup_id`.",
			},

			"specified_backup_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Backup file ID to clone from. Mutually exclusive with `specified_rollback_time`.",
			},

			"uniq_vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "VPC ID.",
			},

			"uniq_subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Subnet ID. If `uniq_vpc_id` is set, this value is required.",
			},

			"memory": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Instance memory in MB. Updatable via `UpgradeDBInstance`.",
			},

			"volume": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Instance disk size in GB. Updatable via `UpgradeDBInstance`.",
			},

			"instance_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Cloned instance name.",
			},

			"security_group": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Security group IDs.",
			},

			"resource_tags": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "Instance tags.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Tag key.",
						},
						"value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Tag value.",
						},
					},
				},
			},

			"cpu": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "CPU cores. Updatable via `UpgradeDBInstance`.",
			},

			"protect_mode": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Data replication mode. 0 - async, 1 - semisync, 2 - strongsync. Updatable via `UpgradeDBInstance`.",
			},

			"deploy_mode": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Deploy mode. 0 - single zone, 1 - multi zone. Updatable via `UpgradeDBInstance`.",
			},

			"slave_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Slave 1 zone. Updatable via `UpgradeDBInstance`.",
			},

			"backup_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Slave 2 zone. Updatable via `UpgradeDBInstance`.",
			},

			"device_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Instance type. Updatable via `UpgradeDBInstance`.",
			},

			"instance_nodes": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Instance node count.",
			},

			"deploy_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Placement group ID.",
			},

			"dry_run": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Dry-run flag. true: only pre-check, false: send normal request.",
			},

			"cage_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Financial cage ID.",
			},

			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Project ID. Default is 0.",
			},

			"pay_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Payment type. PRE_PAID - prepaid, USED_PAID - postpaid. Default is postpaid.",
			},

			"period": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Instance duration in months. Required when `pay_type` is PRE_PAID.",
			},

			"cluster_topology": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: "Cloud disk node topology. Updatable via `UpgradeDBInstance`.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"read_write_node": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "RW node topology.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"zone": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "RW node zone.",
									},
									"node_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Node ID.",
									},
								},
							},
						},
						"read_only_nodes": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "RO node topology list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"zone": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "RO node zone.",
									},
									"node_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Node ID.",
									},
									"is_random_zone": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Whether distributed in a random zone. YES - random zone.",
									},
								},
							},
						},
					},
				},
			},

			"src_region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Source instance region for cross-region clone.",
			},

			"specified_sub_backup_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Cross-region backup ID.",
			},

			"master_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Master zone.",
			},

			"zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Instance zone.",
			},

			"fourth_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Slave 3 zone. Updatable via `UpgradeDBInstance`.",
			},
		},
	}
}

func resourceTencentCloudMysqlCloneInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_clone_instance.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		request = cdb.NewCreateCloneInstanceRequest()
	)

	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("specified_rollback_time"); ok {
		request.SpecifiedRollbackTime = helper.String(v.(string))
	}

	if v, ok := d.GetOk("specified_backup_id"); ok {
		request.SpecifiedBackupId = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("uniq_vpc_id"); ok {
		request.UniqVpcId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("uniq_subnet_id"); ok {
		request.UniqSubnetId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("memory"); ok {
		request.Memory = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("volume"); ok {
		request.Volume = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("instance_name"); ok {
		request.InstanceName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("security_group"); ok {
		securityGroups := v.([]interface{})
		requestSecurityGroup := make([]*string, 0, len(securityGroups))
		for _, item := range securityGroups {
			str := item.(string)
			requestSecurityGroup = append(requestSecurityGroup, &str)
		}
		request.SecurityGroup = requestSecurityGroup
	}

	if v, ok := d.GetOk("resource_tags"); ok {
		for _, item := range v.([]interface{}) {
			tagMap := item.(map[string]interface{})
			tagInfo := cdb.TagInfo{}
			if v, ok := tagMap["key"].(string); ok && v != "" {
				tagInfo.TagKey = helper.String(v)
			}
			if v, ok := tagMap["value"].(string); ok && v != "" {
				tagInfo.TagValue = []*string{helper.String(v)}
			}
			request.ResourceTags = append(request.ResourceTags, &tagInfo)
		}
	}

	if v, ok := d.GetOk("cpu"); ok {
		request.Cpu = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("protect_mode"); ok {
		request.ProtectMode = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("deploy_mode"); ok {
		request.DeployMode = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("slave_zone"); ok {
		request.SlaveZone = helper.String(v.(string))
	}

	if v, ok := d.GetOk("backup_zone"); ok {
		request.BackupZone = helper.String(v.(string))
	}

	if v, ok := d.GetOk("device_type"); ok {
		request.DeviceType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_nodes"); ok {
		request.InstanceNodes = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("deploy_group_id"); ok {
		request.DeployGroupId = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("dry_run"); ok {
		request.DryRun = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("cage_id"); ok {
		request.CageId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("project_id"); ok {
		request.ProjectId = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("pay_type"); ok {
		request.PayType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("period"); ok {
		request.Period = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("cluster_topology"); ok {
		for _, item := range v.([]interface{}) {
			clusterTopologyMap := item.(map[string]interface{})
			clusterTopology := cdb.ClusterTopology{}
			if v, ok := clusterTopologyMap["read_write_node"]; ok {
				for _, rwItem := range v.([]interface{}) {
					rwMap := rwItem.(map[string]interface{})
					readWriteNode := cdb.ReadWriteNode{}
					if v, ok := rwMap["zone"].(string); ok && v != "" {
						readWriteNode.Zone = helper.String(v)
					}
					if v, ok := rwMap["node_id"].(string); ok && v != "" {
						readWriteNode.NodeId = helper.String(v)
					}
					clusterTopology.ReadWriteNode = &readWriteNode
				}
			}
			if v, ok := clusterTopologyMap["read_only_nodes"]; ok {
				for _, roItem := range v.([]interface{}) {
					roMap := roItem.(map[string]interface{})
					readonlyNode := cdb.ReadonlyNode{}
					if v, ok := roMap["zone"].(string); ok && v != "" {
						readonlyNode.Zone = helper.String(v)
					}
					if v, ok := roMap["node_id"].(string); ok && v != "" {
						readonlyNode.NodeId = helper.String(v)
					}
					if v, ok := roMap["is_random_zone"].(string); ok && v != "" {
						readonlyNode.IsRandomZone = helper.String(v)
					}
					clusterTopology.ReadOnlyNodes = append(clusterTopology.ReadOnlyNodes, &readonlyNode)
				}
			}
			request.ClusterTopology = &clusterTopology
		}
	}

	if v, ok := d.GetOk("src_region"); ok {
		request.SrcRegion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("specified_sub_backup_id"); ok {
		request.SpecifiedSubBackupId = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("master_zone"); ok {
		request.MasterZone = helper.String(v.(string))
	}

	if v, ok := d.GetOk("zone"); ok {
		request.Zone = helper.String(v.(string))
	}

	if v, ok := d.GetOk("fourth_zone"); ok {
		request.FourthZone = helper.String(v.(string))
	}

	var asyncRequestId string
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMysqlClient().CreateCloneInstanceWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create mysql_clone_instance failed, Response is nil."))
		}
		if result.Response.AsyncRequestId == nil || *result.Response.AsyncRequestId == "" {
			log.Printf("[CRITAL]%s create mysql_clone_instance, logId=%s, id=%s\n", logId, logId, d.Id())
			return resource.NonRetryableError(fmt.Errorf("Create mysql_clone_instance failed, AsyncRequestId is nil or empty."))
		}
		asyncRequestId = *result.Response.AsyncRequestId
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create mysql_clone_instance failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		taskStatus, message, err := service.DescribeAsyncRequestInfo(ctx, asyncRequestId)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if taskStatus == MYSQL_TASK_STATUS_SUCCESS {
			return nil
		}
		if taskStatus == MYSQL_TASK_STATUS_INITIAL || taskStatus == MYSQL_TASK_STATUS_RUNNING {
			return resource.RetryableError(fmt.Errorf("mysql_clone_instance async task status is %s", taskStatus))
		}
		return resource.NonRetryableError(fmt.Errorf("mysql_clone_instance async task status is %s, message:%s", taskStatus, message))
	})

	if err != nil {
		log.Printf("[CRITAL]%s create mysql_clone_instance async task fail, reason:%s\n", logId, err.Error())
		return err
	}

	sourceInstanceId := d.Get("instance_id").(string)
	var dstInstanceId string
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		cloneList, e := service.DescribeMysqlCloneListByFilter(ctx, map[string]interface{}{
			"InstanceId": helper.String(sourceInstanceId),
		})
		if e != nil {
			return resource.RetryableError(e)
		}
		for _, clone := range cloneList {
			if clone.DstInstanceId != nil && *clone.DstInstanceId != "" {
				dstInstanceId = *clone.DstInstanceId
			}
		}
		if dstInstanceId == "" {
			return resource.RetryableError(fmt.Errorf("mysql_clone_instance DstInstanceId not found in DescribeCloneList"))
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create mysql_clone_instance, get DstInstanceId fail, reason:%s\n", logId, err.Error())
		return err
	}

	log.Printf("[DEBUG]%s create mysql_clone_instance, dstInstanceId=%s, logId=%s, id=%s\n", logId, dstInstanceId, logId, d.Id())
	d.SetId(dstInstanceId)

	return resourceTencentCloudMysqlCloneInstanceRead(d, meta)
}

func resourceTencentCloudMysqlCloneInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_clone_instance.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	respData, err := service.DescribeDBInstanceById(ctx, d.Id())
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[CRUD] mysql_clone_instance id=%s", d.Id())
		d.SetId("")
		return nil
	}

	if respData.Memory != nil {
		_ = d.Set("memory", respData.Memory)
	}

	if respData.Volume != nil {
		_ = d.Set("volume", respData.Volume)
	}

	if respData.Cpu != nil {
		_ = d.Set("cpu", respData.Cpu)
	}

	if respData.ProtectMode != nil {
		_ = d.Set("protect_mode", respData.ProtectMode)
	}

	if respData.DeployMode != nil {
		_ = d.Set("deploy_mode", respData.DeployMode)
	}

	if respData.DeviceType != nil {
		_ = d.Set("device_type", respData.DeviceType)
	}

	if respData.InstanceName != nil {
		_ = d.Set("instance_name", respData.InstanceName)
	}

	if respData.Zone != nil {
		_ = d.Set("zone", respData.Zone)
	}

	if respData.ProjectId != nil {
		_ = d.Set("project_id", respData.ProjectId)
	}

	if respData.DeployGroupId != nil {
		_ = d.Set("deploy_group_id", respData.DeployGroupId)
	}

	if respData.UniqVpcId != nil {
		_ = d.Set("uniq_vpc_id", respData.UniqVpcId)
	}

	if respData.UniqSubnetId != nil {
		_ = d.Set("uniq_subnet_id", respData.UniqSubnetId)
	}

	if respData.SlaveInfo != nil {
		if respData.SlaveInfo.First != nil && respData.SlaveInfo.First.Zone != nil {
			_ = d.Set("slave_zone", respData.SlaveInfo.First.Zone)
		}
		if respData.SlaveInfo.Second != nil && respData.SlaveInfo.Second.Zone != nil {
			_ = d.Set("backup_zone", respData.SlaveInfo.Second.Zone)
		}
	}

	return nil
}

func resourceTencentCloudMysqlCloneInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_clone_instance.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	immutableArgs := []string{
		"instance_id",
		"specified_rollback_time",
		"specified_backup_id",
		"uniq_vpc_id",
		"uniq_subnet_id",
		"instance_name",
		"security_group",
		"resource_tags",
		"instance_nodes",
		"deploy_group_id",
		"dry_run",
		"cage_id",
		"project_id",
		"pay_type",
		"period",
		"src_region",
		"specified_sub_backup_id",
		"master_zone",
		"zone",
	}

	for _, arg := range immutableArgs {
		if d.HasChange(arg) {
			return fmt.Errorf("mysql_clone_instance argument `%s` cannot be modified", arg)
		}
	}

	needChange := false
	mutableArgs := []string{"memory", "volume", "cpu", "protect_mode", "deploy_mode", "slave_zone", "backup_zone", "device_type", "cluster_topology", "fourth_zone"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := cdb.NewUpgradeDBInstanceRequest()
		request.InstanceId = helper.String(d.Id())

		if v, ok := d.GetOk("memory"); ok {
			request.Memory = helper.IntInt64(v.(int))
		}

		if v, ok := d.GetOk("volume"); ok {
			request.Volume = helper.IntInt64(v.(int))
		}

		if v, ok := d.GetOk("cpu"); ok {
			request.Cpu = helper.IntInt64(v.(int))
		}

		if v, ok := d.GetOk("protect_mode"); ok {
			request.ProtectMode = helper.IntInt64(v.(int))
		}

		if v, ok := d.GetOk("deploy_mode"); ok {
			request.DeployMode = helper.IntInt64(v.(int))
		}

		if v, ok := d.GetOk("slave_zone"); ok {
			request.SlaveZone = helper.String(v.(string))
		}

		if v, ok := d.GetOk("backup_zone"); ok {
			request.BackupZone = helper.String(v.(string))
		}

		if v, ok := d.GetOk("device_type"); ok {
			request.DeviceType = helper.String(v.(string))
		}

		if v, ok := d.GetOk("fourth_zone"); ok {
			request.FourthZone = helper.String(v.(string))
		}

		if v, ok := d.GetOk("cluster_topology"); ok {
			for _, item := range v.([]interface{}) {
				clusterTopologyMap := item.(map[string]interface{})
				clusterTopology := cdb.ClusterTopology{}
				if v, ok := clusterTopologyMap["read_write_node"]; ok {
					for _, rwItem := range v.([]interface{}) {
						rwMap := rwItem.(map[string]interface{})
						readWriteNode := cdb.ReadWriteNode{}
						if v, ok := rwMap["zone"].(string); ok && v != "" {
							readWriteNode.Zone = helper.String(v)
						}
						if v, ok := rwMap["node_id"].(string); ok && v != "" {
							readWriteNode.NodeId = helper.String(v)
						}
						clusterTopology.ReadWriteNode = &readWriteNode
					}
				}
				if v, ok := clusterTopologyMap["read_only_nodes"]; ok {
					for _, roItem := range v.([]interface{}) {
						roMap := roItem.(map[string]interface{})
						readonlyNode := cdb.ReadonlyNode{}
						if v, ok := roMap["zone"].(string); ok && v != "" {
							readonlyNode.Zone = helper.String(v)
						}
						if v, ok := roMap["node_id"].(string); ok && v != "" {
							readonlyNode.NodeId = helper.String(v)
						}
						if v, ok := roMap["is_random_zone"].(string); ok && v != "" {
							readonlyNode.IsRandomZone = helper.String(v)
						}
						clusterTopology.ReadOnlyNodes = append(clusterTopology.ReadOnlyNodes, &readonlyNode)
					}
				}
				request.ClusterTopology = &clusterTopology
			}
		}

		var asyncRequestId string
		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMysqlClient().UpgradeDBInstanceWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Upgrade mysql_clone_instance failed, Response is nil."))
			}
			if result.Response.AsyncRequestId == nil || *result.Response.AsyncRequestId == "" {
				log.Printf("[CRITAL]%s update mysql_clone_instance, logId=%s, id=%s\n", logId, logId, d.Id())
				return resource.NonRetryableError(fmt.Errorf("Upgrade mysql_clone_instance failed, AsyncRequestId is nil or empty."))
			}
			asyncRequestId = *result.Response.AsyncRequestId
			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update mysql_clone_instance failed, reason:%+v", logId, reqErr)
			return reqErr
		}

		err := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			taskStatus, message, err := service.DescribeAsyncRequestInfo(ctx, asyncRequestId)
			if err != nil {
				return resource.NonRetryableError(err)
			}
			if taskStatus == MYSQL_TASK_STATUS_SUCCESS {
				return nil
			}
			if taskStatus == MYSQL_TASK_STATUS_INITIAL || taskStatus == MYSQL_TASK_STATUS_RUNNING {
				return resource.RetryableError(fmt.Errorf("mysql_clone_instance upgrade async task status is %s", taskStatus))
			}
			return resource.NonRetryableError(fmt.Errorf("mysql_clone_instance upgrade async task status is %s, message:%s", taskStatus, message))
		})

		if err != nil {
			log.Printf("[CRITAL]%s update mysql_clone_instance async task fail, reason:%s\n", logId, err.Error())
			return err
		}
	}

	return resourceTencentCloudMysqlCloneInstanceRead(d, meta)
}

func resourceTencentCloudMysqlCloneInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_clone_instance.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service    = MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		instanceId = d.Id()
	)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		_, e := service.IsolateDBInstance(ctx, instanceId)
		if e != nil {
			return tccommon.RetryError(e, tccommon.InternalError)
		}
		return nil
	})

	if err != nil {
		return err
	}

	var hasDeleted = false
	err = resource.Retry(7*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		mysqlInfo, err := service.DescribeDBInstanceById(ctx, instanceId)

		if err != nil {
			if _, ok := err.(*sdkError.TencentCloudSDKError); !ok {
				return resource.RetryableError(err)
			} else {
				return resource.NonRetryableError(err)
			}
		}

		if mysqlInfo == nil {
			hasDeleted = true
			return nil
		}
		if *mysqlInfo.Status == MYSQL_STATUS_ISOLATING || *mysqlInfo.Status == MYSQL_STATUS_RUNNING {
			return resource.RetryableError(fmt.Errorf("mysql_clone_instance isolating."))
		}
		if *mysqlInfo.Status == MYSQL_STATUS_ISOLATED {
			return nil
		}
		return resource.NonRetryableError(fmt.Errorf("after IsolateDBInstance mysql_clone_instance status is %d", *mysqlInfo.Status))
	})

	if hasDeleted {
		return nil
	}
	if err != nil {
		return err
	}

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		e := service.OfflineIsolatedInstances(ctx, instanceId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		return nil
	})

	if err != nil {
		return err
	}

	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		mysqlInfo, err := service.DescribeIsolatedDBInstanceById(ctx, instanceId)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if mysqlInfo == nil {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("mysql_clone_instance still in isolated list."))
	})

	if err != nil {
		log.Printf("[CRITAL]%s delete mysql_clone_instance fail, reason:%s\n", logId, err.Error())
		return err
	}

	return nil
}
