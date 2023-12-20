package cdh

import (
	"context"
	"errors"
	"fmt"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCdhInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCdhInstanceCreate,
		Read:   resourceTencentCloudCdhInstanceRead,
		Update: resourceTencentCloudCdhInstanceUpdate,
		Delete: resourceTencentCloudCdhInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"availability_zone": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The available zone for the CDH instance.",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "The project the instance belongs to, default to 0.",
			},
			"host_type": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The type of the CDH instance.",
			},
			"host_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the CDH instance. The max length of host_name is 60.",
			},
			// payment
			"charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      CDH_CHARGE_TYPE_PREPAID,
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{CDH_CHARGE_TYPE_PREPAID}),
				Description:  "The charge type of instance. Valid values are `PREPAID`. The default is `PREPAID`.",
			},
			"prepaid_period": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: tccommon.ValidateAllowedIntValue(CDH_PREPAID_PERIOD),
				Description:  "The tenancy (time unit is month) of the prepaid instance, NOTE: it only works when charge_type is set to `PREPAID`. Valid values are `1`, `2`, `3`, `4`, `5`, `6`, `7`, `8`, `9`, `10`, `11`, `12`, `24`, `36`.",
			},
			"prepaid_renew_flag": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(CDH_PREPAID_RENEW_FLAG),
				Description:  "Auto renewal flag. Valid values: `NOTIFY_AND_AUTO_RENEW`: notify upon expiration and renew automatically, `NOTIFY_AND_MANUAL_RENEW`: notify upon expiration but do not renew automatically, `DISABLE_NOTIFY_AND_MANUAL_RENEW`: neither notify upon expiration nor renew automatically. Default value: `NOTIFY_AND_MANUAL_RENEW`. If this parameter is specified as `NOTIFY_AND_AUTO_RENEW`, the instance will be automatically renewed on a monthly basis if the account balance is sufficient. NOTE: it only works when charge_type is set to `PREPAID`.",
			},
			//computed
			"host_state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "State of the CDH instance.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Create time of the instance.",
			},
			"expired_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Expired time of the instance.",
			},
			"cvm_instance_ids": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Id of CVM instances that have been created on the CDH instance.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"host_resource": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An information list of host resource. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cpu_total_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of total CPU cores of the instance.",
						},
						"cpu_available_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of available CPU cores of the instance.",
						},
						"memory_total_size": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Instance memory total capacity, unit in GB.",
						},
						"memory_available_size": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Instance memory available capacity, unit in GB.",
						},
						"disk_total_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Instance disk total capacity, unit in GB.",
						},
						"disk_available_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Instance disk available capacity, unit in GB.",
						},
						"disk_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of the disk.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudCdhInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cdh_instance.create")()
	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		ctx                  = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		cdhService           = CdhService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		hostType, chargeType string
		placement            cvm.Placement
		hostChargePrepaid    cvm.ChargePrepaid
		hostId               string
		hostInstance         *cvm.HostItem
		outErr, inErr        error
	)

	if v, ok := d.GetOk("host_type"); ok {
		hostType = v.(string)
	}
	if v, ok := d.GetOk("charge_type"); ok {
		chargeType = v.(string)
	}
	placement.Zone = helper.String(d.Get("availability_zone").(string))
	if v, ok := d.GetOk("project_id"); ok {
		placement.ProjectId = helper.IntInt64(v.(int))
	}
	if v, ok := d.GetOk("prepaid_period"); ok {
		hostChargePrepaid.Period = helper.IntUint64(v.(int))
	}
	if v, ok := d.GetOk("prepaid_renew_flag"); ok {
		hostChargePrepaid.RenewFlag = helper.String(v.(string))
	}

	outErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		hostId, inErr = cdhService.CreateCdhInstance(ctx, &placement, &hostChargePrepaid, chargeType, hostType)
		if inErr != nil {
			if sdkErr, ok := inErr.(*sdkErrors.TencentCloudSDKError); ok && sdkErr.Code == CDH_ZONE_SOLD_OUT_FOR_SPECIFIED_INSTANCE_ERROR {
				return resource.NonRetryableError(inErr)
			}
			return tccommon.RetryError(inErr)
		}
		return nil
	})
	if outErr != nil {
		return outErr
	}
	d.SetId(hostId)

	outErr = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		hostInstance, inErr = cdhService.DescribeCdhInstanceById(ctx, d.Id())
		if inErr != nil {
			return tccommon.RetryError(inErr)
		}
		if *hostInstance.HostState == CDH_HOST_STATE_PENDING {
			return resource.RetryableError(errors.New("cdh instance is pending"))
		}
		if *hostInstance.HostState == CDH_HOST_STATE_LAUNCH_FAILURE {
			return resource.NonRetryableError(errors.New("cdh instance launch failure"))
		}
		return nil
	})
	if outErr != nil {
		return outErr
	}

	if v, ok := d.GetOk("host_name"); ok {
		hostName := v.(string)
		outErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			inErr = cdhService.ModifyHostName(ctx, d.Id(), hostName)
			if inErr != nil {
				return tccommon.RetryError(inErr)
			}
			return nil
		})
		if outErr != nil {
			return outErr
		}
	}
	return resourceTencentCloudCdhInstanceRead(d, meta)
}

func resourceTencentCloudCdhInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cdh_instance.read")()
	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		ctx           = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		cdhService    = CdhService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		hostInstance  *cvm.HostItem
		outErr, inErr error
	)

	outErr = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		hostInstance, inErr = cdhService.DescribeCdhInstanceById(ctx, d.Id())
		if inErr != nil {
			return tccommon.RetryError(inErr)
		}
		return nil
	})
	if outErr != nil {
		return outErr
	}

	if hostInstance == nil {
		d.SetId("")
		return nil
	}

	_ = d.Set("availability_zone", hostInstance.Placement.Zone)
	_ = d.Set("project_id", hostInstance.Placement.ProjectId)
	_ = d.Set("host_type", hostInstance.HostType)
	_ = d.Set("host_name", hostInstance.HostName)
	_ = d.Set("charge_type", hostInstance.HostChargeType)
	_ = d.Set("prepaid_renew_flag", hostInstance.RenewFlag)
	_ = d.Set("host_state", hostInstance.HostState)
	_ = d.Set("create_time", hostInstance.CreatedTime)
	_ = d.Set("expired_time", hostInstance.ExpiredTime)
	_ = d.Set("cvm_instance_ids", hostInstance.InstanceIds)

	hostResource := map[string]interface{}{
		"cpu_total_num":         hostInstance.HostResource.CpuTotal,
		"cpu_available_num":     hostInstance.HostResource.CpuAvailable,
		"memory_total_size":     hostInstance.HostResource.MemTotal,
		"memory_available_size": hostInstance.HostResource.MemAvailable,
		"disk_total_size":       hostInstance.HostResource.DiskTotal,
		"disk_available_size":   hostInstance.HostResource.DiskAvailable,
		"disk_type":             hostInstance.HostResource.DiskType,
	}
	_ = d.Set("host_resource", []map[string]interface{}{hostResource})

	return nil
}

func resourceTencentCloudCdhInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cdh_instance.update")()
	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		ctx           = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		cdhService    = CdhService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		outErr, inErr error
	)

	d.Partial(true)

	unsupportedUpdateFields := []string{
		"prepaid_period",
	}
	for _, field := range unsupportedUpdateFields {
		if d.HasChange(field) {
			return fmt.Errorf("tencentcloud_cdh_instance update on %s is not support yet", field)
		}
	}

	if d.HasChange("project_id") {
		outErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			inErr = cdhService.ModifyProject(ctx, d.Id(), d.Get("project_id").(int))
			if inErr != nil {
				return tccommon.RetryError(inErr)
			}
			return nil
		})
		if outErr != nil {
			return outErr
		}

	}

	if d.HasChange("host_name") {
		outErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			inErr = cdhService.ModifyHostName(ctx, d.Id(), d.Get("host_name").(string))
			if inErr != nil {
				return tccommon.RetryError(inErr)
			}
			return nil
		})
		if outErr != nil {
			return outErr
		}

	}

	if d.HasChange("prepaid_renew_flag") {
		outErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			inErr = cdhService.ModifyPrepaidRenewFlag(ctx, d.Id(), d.Get("prepaid_renew_flag").(string))
			if inErr != nil {
				return tccommon.RetryError(inErr)
			}
			return nil
		})
		if outErr != nil {
			return outErr
		}

	}

	d.Partial(false)

	return resourceTencentCloudCdhInstanceRead(d, meta)
}

func resourceTencentCloudCdhInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cdh_instance.delete")()

	return fmt.Errorf("PREPAID CDH instance do not support delete operation with terraform")
}
