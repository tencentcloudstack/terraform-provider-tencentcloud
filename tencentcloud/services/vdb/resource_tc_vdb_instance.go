package vdb

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vdb/v20230616"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudVdbInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVdbInstanceCreate,
		Read:   resourceTencentCloudVdbInstanceRead,
		Update: resourceTencentCloudVdbInstanceUpdate,
		Delete: resourceTencentCloudVdbInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			// --- Create input parameters (non-deprecated) ---
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "VPC ID.",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Subnet ID.",
			},
			"pay_mode": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Billing mode. 0: pay-as-you-go, 1: monthly subscription.",
			},
			"instance_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Instance name. Supports up to 60 characters.",
			},
			"security_group_ids": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Security group IDs.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"pay_period": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Monthly subscription period in months. Valid values: 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 24, 36. Default is 1.",
			},
			"auto_renew": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Auto-renew flag. 0: disabled, 1: enabled.",
			},
			"params": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Instance extra parameters, submitted via JSON.",
			},
			"resource_tags": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Tag list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Tag key.",
						},
						"tag_value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Tag value.",
						},
					},
				},
			},
			"instance_type": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Instance type. Valid values: base, single, cluster.",
			},
			"mode": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Availability zone mode for cluster type. Valid values: two, three.",
			},
			"product_type": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "Product version. 0: standard, 1: capacity-enhanced.",
			},
			"node_type": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Node type. Valid values: compute, normal, store.",
			},
			"cpu": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "CPU cores.",
			},
			"memory": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Memory size in GB.",
			},
			"disk_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Disk size in GB.",
			},
			"worker_node_num": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Number of worker nodes.",
			},
			"force_delete": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether to force delete (destroy) the instance. If false, only isolate to recycle bin. If true, isolate then destroy. Default is false.",
			},
			// --- Computed fields from InstanceInfo response ---
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Instance status.",
			},
			"region": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Region.",
			},
			"zone": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Availability zone.",
			},
			"product": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Product.",
			},
			"shard_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Shard number.",
			},
			"api_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "API version.",
			},
			"extend": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Extended information in JSON format.",
			},
			"expired_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Expiration time.",
			},
			"is_no_expired": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the instance never expires.",
			},
			"wan_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Public network address.",
			},
			"isolate_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Isolation time.",
			},
			"task_status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Task status. 0: no task, 1: pending, 2-11: various operations in progress.",
			},
			"networks": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Network information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VPC ID.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Subnet ID.",
						},
						"vip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Internal IP.",
						},
						"port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Internal port.",
						},
						"preserve_duration": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Old IP preservation duration in days.",
						},
						"expire_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Old IP expiration time.",
						},
					},
				},
			},
			"nodes": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Instance node list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Pod name.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Pod status.",
						},
					},
				},
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time.",
			},
			"engine_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Engine name.",
			},
			"engine_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Engine version.",
			},
		},
	}
}

func resourceTencentCloudVdbInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vdb_instance.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request  = vdb.NewCreateInstanceRequest()
		response = vdb.NewCreateInstanceResponse()
	)

	if v, ok := d.GetOk("vpc_id"); ok {
		request.VpcId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("subnet_id"); ok {
		request.SubnetId = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("pay_mode"); ok {
		request.PayMode = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("instance_name"); ok {
		request.InstanceName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("security_group_ids"); ok {
		for _, item := range v.([]interface{}) {
			request.SecurityGroupIds = append(request.SecurityGroupIds, helper.String(item.(string)))
		}
	}

	if v, ok := d.GetOkExists("pay_period"); ok {
		request.PayPeriod = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("auto_renew"); ok {
		request.AutoRenew = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("params"); ok {
		request.Params = helper.String(v.(string))
	}

	if v, ok := d.GetOk("resource_tags"); ok {
		for _, item := range v.([]interface{}) {
			tagMap := item.(map[string]interface{})
			tag := vdb.Tag{}
			if v, ok := tagMap["tag_key"].(string); ok && v != "" {
				tag.TagKey = helper.String(v)
			}
			if v, ok := tagMap["tag_value"].(string); ok && v != "" {
				tag.TagValue = helper.String(v)
			}
			request.ResourceTags = append(request.ResourceTags, &tag)
		}
	}

	if v, ok := d.GetOk("instance_type"); ok {
		request.InstanceType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("mode"); ok {
		request.Mode = helper.String(v.(string))
	}

	request.GoodsNum = helper.IntInt64(1)

	if v, ok := d.GetOkExists("product_type"); ok {
		request.ProductType = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("node_type"); ok {
		request.NodeType = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("cpu"); ok {
		request.Cpu = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("memory"); ok {
		request.Memory = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("disk_size"); ok {
		request.DiskSize = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("worker_node_num"); ok {
		request.WorkerNodeNum = helper.IntUint64(v.(int))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVdbV20230616Client().CreateInstanceWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create vdb instance failed, Response is nil"))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create vdb instance failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.InstanceIds == nil || len(response.Response.InstanceIds) == 0 {
		return fmt.Errorf("InstanceIds is empty")
	}

	instanceId := *response.Response.InstanceIds[0]
	d.SetId(instanceId)

	service := VdbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	if err := service.WaitForInstanceStatus(ctx, instanceId, "online", d.Timeout(schema.TimeoutCreate)); err != nil {
		return err
	}

	return resourceTencentCloudVdbInstanceRead(d, meta)
}

func resourceTencentCloudVdbInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vdb_instance.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service    = VdbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		instanceId = d.Id()
	)

	instance, err := service.DescribeVdbInstanceById(ctx, instanceId)
	if err != nil {
		return err
	}

	if instance == nil {
		log.Printf("[WARN]%s resource `tencentcloud_vdb_instance` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if instance.Name != nil {
		_ = d.Set("instance_name", instance.Name)
	}

	if instance.Status != nil {
		_ = d.Set("status", instance.Status)
	}

	if instance.PayMode != nil {
		_ = d.Set("pay_mode", instance.PayMode)
	}

	if instance.Cpu != nil {
		_ = d.Set("cpu", int(*instance.Cpu))
	}

	if instance.Memory != nil {
		_ = d.Set("memory", int(*instance.Memory))
	}

	if instance.Disk != nil {
		_ = d.Set("disk_size", instance.Disk)
	}

	if instance.InstanceType != nil {
		_ = d.Set("instance_type", instance.InstanceType)
	}

	if instance.NodeType != nil {
		_ = d.Set("node_type", instance.NodeType)
	}

	if instance.ProductType != nil {
		_ = d.Set("product_type", instance.ProductType)
	}

	if instance.AutoRenew != nil {
		_ = d.Set("auto_renew", instance.AutoRenew)
	}

	if instance.CreatedAt != nil {
		_ = d.Set("created_at", instance.CreatedAt)
	}

	if instance.EngineName != nil {
		_ = d.Set("engine_name", instance.EngineName)
	}

	if instance.EngineVersion != nil {
		_ = d.Set("engine_version", instance.EngineVersion)
	}

	if instance.Region != nil {
		_ = d.Set("region", instance.Region)
	}

	if instance.Zone != nil {
		_ = d.Set("zone", instance.Zone)
	}

	if instance.Product != nil {
		_ = d.Set("product", instance.Product)
	}

	if instance.ShardNum != nil {
		_ = d.Set("shard_num", instance.ShardNum)
	}

	if instance.ApiVersion != nil {
		_ = d.Set("api_version", instance.ApiVersion)
	}

	if instance.Extend != nil {
		_ = d.Set("extend", instance.Extend)
	}

	if instance.ExpiredAt != nil {
		_ = d.Set("expired_at", instance.ExpiredAt)
	}

	if instance.IsNoExpired != nil {
		_ = d.Set("is_no_expired", instance.IsNoExpired)
	}

	if instance.WanAddress != nil {
		_ = d.Set("wan_address", instance.WanAddress)
	}

	if instance.IsolateAt != nil {
		_ = d.Set("isolate_at", instance.IsolateAt)
	}

	if instance.TaskStatus != nil {
		_ = d.Set("task_status", instance.TaskStatus)
	}

	if instance.Networks != nil {
		networksList := make([]map[string]interface{}, 0, len(instance.Networks))
		for _, n := range instance.Networks {
			networkMap := map[string]interface{}{}
			if n.VpcId != nil {
				networkMap["vpc_id"] = n.VpcId
			}
			if n.SubnetId != nil {
				networkMap["subnet_id"] = n.SubnetId
			}
			if n.Vip != nil {
				networkMap["vip"] = n.Vip
			}
			if n.Port != nil {
				networkMap["port"] = n.Port
			}
			if n.PreserveDuration != nil {
				networkMap["preserve_duration"] = n.PreserveDuration
			}
			if n.ExpireTime != nil {
				networkMap["expire_time"] = n.ExpireTime
			}
			networksList = append(networksList, networkMap)
		}
		_ = d.Set("networks", networksList)

		if len(instance.Networks) > 0 {
			if instance.Networks[0].VpcId != nil {
				_ = d.Set("vpc_id", instance.Networks[0].VpcId)
			}
			if instance.Networks[0].SubnetId != nil {
				_ = d.Set("subnet_id", instance.Networks[0].SubnetId)
			}
		}
	}

	if instance.ResourceTags != nil {
		tagsList := make([]map[string]interface{}, 0, len(instance.ResourceTags))
		for _, tag := range instance.ResourceTags {
			tagMap := map[string]interface{}{}
			if tag.TagKey != nil {
				tagMap["tag_key"] = tag.TagKey
			}
			if tag.TagValue != nil {
				tagMap["tag_value"] = tag.TagValue
			}
			tagsList = append(tagsList, tagMap)
		}
		_ = d.Set("resource_tags", tagsList)
	}

	if instance.SecurityGroupIds != nil {
		_ = d.Set("security_group_ids", instance.SecurityGroupIds)
	}

	if instance.ReplicaNum != nil {
		_ = d.Set("worker_node_num", instance.ReplicaNum)
	}

	// Read instance nodes
	nodes, err := service.DescribeVdbInstanceNodesById(ctx, instanceId)
	if err != nil {
		log.Printf("[WARN]%s read vdb instance nodes failed, reason:%+v", logId, err)
	} else if nodes != nil {
		nodesList := make([]map[string]interface{}, 0, len(nodes))
		for _, node := range nodes {
			nodeMap := map[string]interface{}{}
			if node.Name != nil {
				nodeMap["name"] = node.Name
			}
			if node.Status != nil {
				nodeMap["status"] = node.Status
			}
			nodesList = append(nodesList, nodeMap)
		}
		_ = d.Set("nodes", nodesList)
	}

	return nil
}

func resourceTencentCloudVdbInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vdb_instance.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		instanceId = d.Id()
		service    = VdbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	// These fields cannot be modified in-place. Any change to them will be rejected.
	immutableFields := []string{
		"pay_mode", "instance_name",
		"pay_period", "auto_renew", "params", "resource_tags",
		"instance_type", "mode", "product_type", "node_type",
	}
	for _, field := range immutableFields {
		if d.HasChange(field) {
			return fmt.Errorf("argument `%s` cannot be changed", field)
		}
	}

	// Serial execution: ScaleUp first, then ScaleOut
	// Step 1: Scale up (CPU, Memory, DiskSize)
	if d.HasChange("cpu") || d.HasChange("memory") || d.HasChange("disk_size") {
		request := vdb.NewScaleUpInstanceRequest()
		request.InstanceId = helper.String(instanceId)

		if v, ok := d.GetOkExists("cpu"); ok {
			cpu := float64(v.(int))
			request.Cpu = &cpu
		}

		if v, ok := d.GetOkExists("memory"); ok {
			mem := float64(v.(int))
			request.Memory = &mem
		}

		if v, ok := d.GetOkExists("disk_size"); ok {
			request.StorageSize = helper.IntUint64(v.(int))
		}

		runNow := true
		request.RunNow = &runNow

		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVdbV20230616Client().ScaleUpInstanceWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			}
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s scale up vdb instance failed, reason:%+v", logId, reqErr)
			return reqErr
		}

		targetCpu := float64(d.Get("cpu").(int))
		targetMemory := float64(d.Get("memory").(int))
		targetDiskSize := uint64(d.Get("disk_size").(int))
		if err := service.WaitForInstanceScaleUp(ctx, instanceId, targetCpu, targetMemory, targetDiskSize, d.Timeout(schema.TimeoutUpdate)); err != nil {
			return err
		}
	}

	// Step 2: Scale out (WorkerNodeNum) — executes after ScaleUp is complete
	if d.HasChange("worker_node_num") {
		request := vdb.NewScaleOutInstanceRequest()
		request.InstanceId = helper.String(instanceId)

		if v, ok := d.GetOkExists("worker_node_num"); ok {
			request.ReplicaNum = helper.IntUint64(v.(int))
		}

		runNow := true
		request.RunNow = &runNow

		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVdbV20230616Client().ScaleOutInstanceWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			}
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s scale out vdb instance failed, reason:%+v", logId, reqErr)
			return reqErr
		}

		targetReplicaNum := uint64(d.Get("worker_node_num").(int))
		if err := service.WaitForInstanceScaleOut(ctx, instanceId, targetReplicaNum, d.Timeout(schema.TimeoutUpdate)); err != nil {
			return err
		}
	}

	// Step 3: Update security groups
	if d.HasChange("security_group_ids") {
		sgIds := make([]*string, 0)
		targetSgIds := make([]string, 0)
		for _, item := range d.Get("security_group_ids").([]interface{}) {
			id := item.(string)
			sgIds = append(sgIds, helper.String(id))
			targetSgIds = append(targetSgIds, id)
		}

		request := vdb.NewModifyDBInstanceSecurityGroupsRequest()
		request.InstanceIds = []*string{helper.String(instanceId)}
		request.SecurityGroupIds = sgIds

		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVdbV20230616Client().ModifyDBInstanceSecurityGroupsWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			}
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s modify vdb instance security groups failed, reason:%+v", logId, reqErr)
			return reqErr
		}

		if err := service.WaitForSecurityGroupsMatch(ctx, instanceId, targetSgIds, d.Timeout(schema.TimeoutUpdate)); err != nil {
			return err
		}
	}

	return resourceTencentCloudVdbInstanceRead(d, meta)
}

func resourceTencentCloudVdbInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vdb_instance.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		instanceId = d.Id()
		service    = VdbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	// Step 1: Isolate instance
	isolateRequest := vdb.NewIsolateInstanceRequest()
	isolateRequest.InstanceId = helper.String(instanceId)

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVdbV20230616Client().IsolateInstanceWithContext(ctx, isolateRequest)
		if e != nil {
			return tccommon.RetryError(e)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, isolateRequest.GetAction(), isolateRequest.ToJsonString(), result.ToJsonString())
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s isolate vdb instance failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	// Wait for isolated status
	if err := service.WaitForInstanceStatus(ctx, instanceId, "isolated", d.Timeout(schema.TimeoutDelete)); err != nil {
		return err
	}

	// Step 2: If force_delete, destroy the instance
	forceDelete := d.Get("force_delete").(bool)
	if forceDelete {
		destroyRequest := vdb.NewDestroyInstancesRequest()
		destroyRequest.InstanceIds = []*string{helper.String(instanceId)}

		reqErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVdbV20230616Client().DestroyInstancesWithContext(ctx, destroyRequest)
			if e != nil {
				return tccommon.RetryError(e)
			}
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, destroyRequest.GetAction(), destroyRequest.ToJsonString(), result.ToJsonString())
			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s destroy vdb instance failed, reason:%+v", logId, reqErr)
			return reqErr
		}

		// Wait until instance is completely gone
		if err := service.WaitForInstanceNotFound(ctx, instanceId, d.Timeout(schema.TimeoutDelete)); err != nil {
			return err
		}
	}

	return nil
}
