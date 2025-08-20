package dlc

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dlcv20210125 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dlc/v20210125"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudDlcUpdateStandardEngineResourceGroupResourceInformationOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDlcUpdateStandardEngineResourceGroupResourceInformationOperationCreate,
		Read:   resourceTencentCloudDlcUpdateStandardEngineResourceGroupResourceInformationOperationRead,
		Delete: resourceTencentCloudDlcUpdateStandardEngineResourceGroupResourceInformationOperationDelete,
		Schema: map[string]*schema.Schema{
			"engine_resource_group_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Engine resource group name.",
			},

			"driver_cu_spec": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Driver CU specifications:\nCurrently supported: small (default, 1 CU), medium (2 CU), large (4 CU), xlarge (8 CU). Memory CUs are CPUs with a ratio of 1:8, m.small (1 CU memory), m.medium (2 CU memory), m.large (4 CU memory), and m.xlarge (8 CU memory).",
			},

			"executor_cu_spec": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Executor CU specifications:\nCurrently supported: small (default, 1 CU), medium (2 CU), large (4 CU), xlarge (8 CU). Memory CUs are CPUs with a ratio of 1:8, m.small (1 CU memory), m.medium (2 CU memory), m.large (4 CU memory), and m.xlarge (8 CU memory).",
			},

			"min_executor_nums": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "Minimum number of executors.",
			},

			"max_executor_nums": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "Maximum number of executors.",
			},

			"size": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "AI resource group resource limit.",
			},

			"image_type": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Image type, built-in image: built-in, custom image: custom.",
			},

			"image_name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Image name.",
			},

			"image_version": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Image version, image id.",
			},

			"frame_type": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Framework Type.",
			},

			"public_domain": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Customized mirror domain name.",
			},

			"registry_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Custom image instance id.",
			},

			"region_name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Customize the image region.",
			},

			"python_cu_spec": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The resource limit for a Python stand-alone node in a Python resource group must be smaller than the resource limit for the resource group. Small: 1cu Medium: 2cu Large: 4cu Xlarge: 8cu 4xlarge: 16cu 8xlarge: 32cu 16xlarge: 64cu. If the resource type is high memory, add m before the type.",
			},

			"spark_spec_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Only SQL resource group resource configuration mode, fast: fast mode, custom: custom mode.",
			},

			"spark_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "SQL resource group resource limit only, only used in fast mode.",
			},
		},
	}
}

func resourceTencentCloudDlcUpdateStandardEngineResourceGroupResourceInformationOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_update_standard_engine_resource_group_resource_information_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId                   = tccommon.GetLogId(tccommon.ContextNil)
		ctx                     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request                 = dlcv20210125.NewUpdateStandardEngineResourceGroupResourceInfoRequest()
		engineResourceGroupName string
	)

	if v, ok := d.GetOk("engine_resource_group_name"); ok {
		request.EngineResourceGroupName = helper.String(v.(string))
		engineResourceGroupName = v.(string)
	}

	if v, ok := d.GetOk("driver_cu_spec"); ok {
		request.DriverCuSpec = helper.String(v.(string))
	}

	if v, ok := d.GetOk("executor_cu_spec"); ok {
		request.ExecutorCuSpec = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("min_executor_nums"); ok {
		request.MinExecutorNums = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("max_executor_nums"); ok {
		request.MaxExecutorNums = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("size"); ok {
		request.Size = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("image_type"); ok {
		request.ImageType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("image_name"); ok {
		request.ImageName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("image_version"); ok {
		request.ImageVersion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("frame_type"); ok {
		request.FrameType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("public_domain"); ok {
		request.PublicDomain = helper.String(v.(string))
	}

	if v, ok := d.GetOk("registry_id"); ok {
		request.RegistryId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("region_name"); ok {
		request.RegionName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("python_cu_spec"); ok {
		request.PythonCuSpec = helper.String(v.(string))
	}

	if v, ok := d.GetOk("spark_spec_mode"); ok {
		request.SparkSpecMode = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("spark_size"); ok {
		request.SparkSize = helper.IntInt64(v.(int))
	}

	request.IsEffectiveNow = helper.IntInt64(0)
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().UpdateStandardEngineResourceGroupResourceInfoWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create dlc update standard engine resource group resource information operation failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	d.SetId(engineResourceGroupName)

	// wait
	waitReq := dlcv20210125.NewDescribeStandardEngineResourceGroupsRequest()
	waitReq.Filters = []*dlcv20210125.Filter{
		{
			Name:   helper.String("engine-resource-group-name-unique"),
			Values: helper.Strings([]string{engineResourceGroupName}),
		},
	}
	reqErr = resource.Retry(tccommon.ReadRetryTimeout*7, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().DescribeStandardEngineResourceGroupsWithContext(ctx, waitReq)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe standard engine resource groups failed, Response is nil."))
		}

		if result.Response.UserEngineResourceGroupInfos == nil || len(result.Response.UserEngineResourceGroupInfos) == 0 {
			return resource.NonRetryableError(fmt.Errorf("UserEngineResourceGroupInfos is nil."))
		}

		if result.Response.UserEngineResourceGroupInfos[0].ResourceGroupState == nil {
			return resource.NonRetryableError(fmt.Errorf("ResourceGroupState is nil."))
		}

		resourceGroupState := result.Response.UserEngineResourceGroupInfos[0].ResourceGroupState
		if *resourceGroupState == 2 {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("ResourceGroupState is not running, current state is %d...", *resourceGroupState))
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s describe standard engine resource groups failed,, reason:%+v", logId, reqErr)
		return reqErr
	}

	return resourceTencentCloudDlcUpdateStandardEngineResourceGroupResourceInformationOperationRead(d, meta)
}

func resourceTencentCloudDlcUpdateStandardEngineResourceGroupResourceInformationOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_update_standard_engine_resource_group_resource_information_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudDlcUpdateStandardEngineResourceGroupResourceInformationOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_update_standard_engine_resource_group_resource_information_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
