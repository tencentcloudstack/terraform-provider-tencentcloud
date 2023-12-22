package vpc

import (
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudHaVip() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudHaVipCreate,
		Read:   resourceTencentCloudHaVipRead,
		Update: resourceTencentCloudHaVipUpdate,
		Delete: resourceTencentCloudHaVipDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(1, 60),
				Description:  "Name of the HA VIP. The length of character is limited to 1-60.",
			},
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
			"vip": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: tccommon.ValidateIp,
				Description:  "Virtual IP address, it must not be occupied and in this VPC network segment. If not set, it will be assigned after resource created automatically.",
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "State of the HA VIP. Valid value: `AVAILABLE`, `UNBIND`.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Instance ID that is associated.",
			},
			"network_interface_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Network interface ID that is associated.",
			},
			"address_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "EIP that is associated.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Create time of the HA VIP.",
			},
		},
	}
}

func resourceTencentCloudHaVipCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ha_vip.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := vpc.NewCreateHaVipRequest()
	request.VpcId = helper.String(d.Get("vpc_id").(string))
	request.SubnetId = helper.String(d.Get("subnet_id").(string))
	request.HaVipName = helper.String(d.Get("name").(string))
	//optional
	if v, ok := d.GetOk("vip"); ok {
		request.Vip = helper.String(v.(string))
	}
	var response *vpc.CreateHaVipResponse
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().CreateHaVip(request)
		if e != nil {
			return tccommon.RetryError(errors.WithStack(e))
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create HA VIP failed, reason:%+v", logId, err)
		return err
	}

	if response.Response.HaVip == nil {
		return fmt.Errorf("HA VIP id is nil")
	}
	haVipId := *response.Response.HaVip.HaVipId
	d.SetId(haVipId)

	return resourceTencentCloudHaVipRead(d, meta)
}

func resourceTencentCloudHaVipRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ha_vip.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	haVipId := d.Id()
	request := vpc.NewDescribeHaVipsRequest()
	request.HaVipIds = []*string{&haVipId}

	var response *vpc.DescribeHaVipsResponse
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().DescribeHaVips(request)
		if e != nil {
			return tccommon.RetryError(errors.WithStack(e))
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read HA VIP failed, reason:%+v", logId, err)
		return err
	}
	if len(response.Response.HaVipSet) < 1 {
		d.SetId("")
		return nil
	}

	haVip := response.Response.HaVipSet[0]

	_ = d.Set("name", *haVip.HaVipName)
	_ = d.Set("create_time", *haVip.CreatedTime)
	_ = d.Set("vip", *haVip.Vip)
	_ = d.Set("vpc_id", *haVip.VpcId)
	_ = d.Set("subnet_id", *haVip.SubnetId)
	_ = d.Set("address_ip", *haVip.AddressIp)
	_ = d.Set("state", *haVip.State)
	_ = d.Set("network_interface_id", *haVip.NetworkInterfaceId)
	_ = d.Set("instance_id", *haVip.InstanceId)

	return nil
}

func resourceTencentCloudHaVipUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ha_vip.update")()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	haVipId := d.Id()
	request := vpc.NewModifyHaVipAttributeRequest()
	request.HaVipId = &haVipId
	if d.HasChange("name") {
		request.HaVipName = helper.String(d.Get("name").(string))
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			_, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().ModifyHaVipAttribute(request)
			if e != nil {
				return tccommon.RetryError(errors.WithStack(e))
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s modify HA VIP failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudHaVipRead(d, meta)
}

func resourceTencentCloudHaVipDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ha_vip.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	haVipId := d.Id()

	request := vpc.NewDeleteHaVipRequest()
	request.HaVipId = &haVipId
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		_, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().DeleteHaVip(request)
		if e != nil {
			return tccommon.RetryError(errors.WithStack(e), VPCUnsupportedOperation)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete HA VIP failed, reason:%+v", logId, err)
		return err
	}
	//to get the status of haVip
	statRequest := vpc.NewDescribeHaVipsRequest()
	statRequest.HaVipIds = []*string{&haVipId}
	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().DescribeHaVips(statRequest)
		if e != nil {
			ee, ok := e.(*sdkErrors.TencentCloudSDKError)
			if !ok {
				return tccommon.RetryError(errors.WithStack(ee))
			}
			if ee.Code == VPCNotFound {
				log.Printf("[CRITAL]%s api[%s] success, request body [%s], reason[%s]\n",
					logId, statRequest.GetAction(), statRequest.ToJsonString(), e)
				return nil
			} else {
				//when associated eip is in deleting process, delete ha vip may return unsupported operation error
				return tccommon.RetryError(errors.WithStack(e), VPCUnsupportedOperation)
			}
		} else {
			//if not, quit
			if len(result.Response.HaVipSet) == 0 {
				return nil
			}
			//else consider delete fail
			return resource.RetryableError(fmt.Errorf("deleting retry"))
		}
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete HA VIP failed, reason:%+v", logId, err)
		return err
	}
	return nil
}
