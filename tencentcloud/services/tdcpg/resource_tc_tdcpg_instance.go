package tdcpg

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tdcpg "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdcpg/v20211118"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTdcpgInstance() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudTdcpgInstanceRead,
		Create: resourceTencentCloudTdcpgInstanceCreate,
		Update: resourceTencentCloudTdcpgInstanceUpdate,
		Delete: resourceTencentCloudTdcpgInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "cluster id.",
			},

			"cpu": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "cpu cores.",
			},

			"memory": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "memory size.",
			},

			"instance_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "instance name.",
			},

			"operation_timing": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "operation timing, optional value is IMMEDIATE or MAINTAIN_PERIOD.",
			},
		},
	}
}

func resourceTencentCloudTdcpgInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tdcpg_instance.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request    = tdcpg.NewCreateClusterInstancesRequest()
		response   *tdcpg.CreateClusterInstancesResponse
		service    = TdcpgService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		clusterId  string
		instanceId string
		dealNames  []*string
	)

	if v, ok := d.GetOk("cluster_id"); ok {
		request.ClusterId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cpu"); ok {
		request.CPU = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("memory"); ok {
		request.Memory = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("instance_name"); ok {
		request.InstanceName = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTdcpgClient().CreateClusterInstances(request)
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
		log.Printf("[CRITICAL]%s create tdcpg instance failed, reason:%+v", logId, err)
		return err
	}

	err = resource.Retry(3*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		dealNames = response.Response.DealNameSet
		resources, e := service.DescribeTdcpgResourceByDealName(ctx, dealNames)

		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s call api[%s] success, request body [%s], resources [%v]\n",
				logId, "DescribeTdcpgResourceByDealName", request.ToJsonString(), resources)
		}
		clusterId = *resources[0].ClusterId
		instanceId = *resources[0].InstanceIdSet[0]
		return nil
	})

	if err != nil {
		log.Printf("[CRITICAL]%s query tdcpg cluster resource by deal name:[%v] failed, reason:%+v", logId, dealNames, err)
		return err
	}

	d.SetId(clusterId + tccommon.FILED_SP + instanceId)
	return resourceTencentCloudTdcpgInstanceRead(d, meta)
}

func resourceTencentCloudTdcpgInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tdcpg_instance.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = TdcpgService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		instance   *tdcpg.Instance
		ids        = strings.Split(d.Id(), tccommon.FILED_SP)
		clusterId  = ids[0]
		instanceId = ids[1]
	)

	// query the instance of cluster
	err := resource.Retry(3*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		instances, e := service.DescribeTdcpgInstance(ctx, &clusterId, &instanceId)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if instances != nil && instances.InstanceSet != nil {
			status := *instances.InstanceSet[0].Status

			if status == "running" {
				instance = instances.InstanceSet[0]
				return nil
			}

			if status == "creating" || status == "recovering" {
				return resource.RetryableError(fmt.Errorf("tdcpg instance[%s] status is still creating or recovering, retry...", d.Id()))
			}
			return resource.NonRetryableError(fmt.Errorf("tdcpg instance[%s] status is invalid, exit!", d.Id()))
		}
		return resource.RetryableError(fmt.Errorf("can not get tdcpg instance[%s] status, retry...", d.Id()))
	})
	if err != nil {
		return err
	}

	if instance == nil {
		d.SetId("")
		return fmt.Errorf("resource `instance` %s does not exist", instanceId)
	}

	if instance.ClusterId != nil {
		_ = d.Set("cluster_id", instance.ClusterId)
	}

	if instance.CPU != nil {
		_ = d.Set("cpu", instance.CPU)
	}

	if instance.Memory != nil {
		_ = d.Set("memory", instance.Memory)
	}

	if instance.InstanceName != nil {
		_ = d.Set("instance_name", instance.InstanceName)
	}

	return nil
}

func resourceTencentCloudTdcpgInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tdcpg_instance.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = TdcpgService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		request    = tdcpg.NewModifyClusterInstancesSpecRequest()
		ids        = strings.Split(d.Id(), tccommon.FILED_SP)
		clusterId  = ids[0]
		instanceId = ids[1]
	)

	request.ClusterId = &clusterId
	request.InstanceIdSet = []*string{helper.String(instanceId)}

	if v, ok := d.GetOk("cpu"); ok {
		request.CPU = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("memory"); ok {
		request.Memory = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("operation_timing"); ok {
		request.OperationTiming = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTdcpgClient().ModifyClusterInstancesSpec(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITICAL]%s modify tdcpg instance failed, reason:%+v", logId, err)
		return err
	}

	// check the instance value to make sure modify successfully.
	err = resource.Retry(3*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		instances, e := service.DescribeTdcpgInstance(ctx, &clusterId, &instanceId)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if instances != nil && instances.InstanceSet != nil {
			instance := *instances.InstanceSet[0]

			if *instance.Status == "running" {
				if int(*instance.CPU) != d.Get("cpu").(int) || int(*instance.Memory) != d.Get("memory").(int) {
					return resource.RetryableError(fmt.Errorf("the modify instance[%s] operation still on going, retry...", d.Id()))
				}
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("tdcpg instance[%s] status is invalid, exit!", d.Id()))
		}
		return resource.RetryableError(fmt.Errorf("can not get tdcpg instance[%s] status, retry...", d.Id()))
	})

	if err != nil {
		return err
	}

	return resourceTencentCloudTdcpgInstanceRead(d, meta)
}

func resourceTencentCloudTdcpgInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tdcpg_instance.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = TdcpgService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		ids        = strings.Split(d.Id(), tccommon.FILED_SP)
		clusterId  = ids[0]
		instanceId = ids[1]
	)

	if err := service.DeleteTdcpgInstanceById(ctx, &clusterId, &instanceId); err != nil {
		return err
	}

	return nil
}
