package vpc

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudHaVipEipAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudHaVipEipAttachmentCreate,
		Read:   resourceTencentCloudHaVipEipAttachmentRead,
		Delete: resourceTencentCloudHaVipEipAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"havip_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the attached HA VIP.",
			},
			"address_ip": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: tccommon.ValidateIp,
				Description:  "Public address of the EIP.",
			},
		},
	}
}

func resourceTencentCloudHaVipEipAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ha_vip_eip_attachment.create")()

	haVipId := d.Get("havip_id").(string)
	addressIp := d.Get("address_ip").(string)

	bindErr := haVipAssociateEip(meta, haVipId, addressIp)
	if bindErr != nil {
		return bindErr
	}

	d.SetId(haVipId + "#" + addressIp)

	return resourceTencentCloudHaVipEipAttachmentRead(d, meta)
}

func resourceTencentCloudHaVipEipAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ha_vip_eip_attachment.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	haVipEipAttachmentId := d.Id()

	eip := ""
	haVip := ""
	has := false
	vpcService := VpcService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		eipId, haVipId, flag, e := vpcService.DescribeHaVipEipById(ctx, haVipEipAttachmentId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		has = flag
		eip = eipId
		haVip = haVipId
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read HA VIP EIP attachment failed, reason:%s\n", logId, err)
		return err
	}

	if !has {
		d.SetId("")
		return nil
	}
	_ = d.Set("havip_id", haVip)
	_ = d.Set("address_ip", eip)
	d.SetId(haVipEipAttachmentId)

	return nil
}

func resourceTencentCloudHaVipEipAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ha_vip_eip_attachment.delete")()

	haVipEipAttachmentId := d.Id()
	items := strings.Split(haVipEipAttachmentId, "#")
	if len(items) != 2 {
		return fmt.Errorf("decode HA VIP EIP attachment id error")
	}
	haVipId := items[0]
	addressIp := items[1]

	unBindErr := haVipDisassociateEip(meta, haVipId, addressIp)
	if unBindErr != nil {
		return unBindErr
	}

	return nil
}

func haVipAssociateEip(meta interface{}, havipId string, eip string) error {
	//associate eip
	logId := tccommon.GetLogId(tccommon.ContextNil)
	bindRequest := vpc.NewHaVipAssociateAddressIpRequest()
	bindRequest.HaVipId = helper.String(havipId)
	bindRequest.AddressIp = helper.String(eip)
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		_, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().HaVipAssociateAddressIp(bindRequest)
		if e != nil {
			return tccommon.RetryError(errors.WithStack(e))
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create HA VIP EIP attachment failed, reason:%+v", logId, err)
		return err
	}

	statRequest := vpc.NewDescribeHaVipsRequest()
	statRequest.HaVipIds = []*string{&havipId}
	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().DescribeHaVips(statRequest)
		if e != nil {
			return tccommon.RetryError(errors.WithStack(e), VPCUnsupportedOperation)
		} else {
			if len(result.Response.HaVipSet) > 0 {
				if *result.Response.HaVipSet[0].AddressIp == "" {
					return resource.RetryableError(fmt.Errorf("Not binded yet, retry describing"))
				} else {
					return nil
				}
			}
			return resource.NonRetryableError(fmt.Errorf("describe error"))
		}
	})
	if err != nil {
		log.Printf("[CRITAL]%s describe HA VIP failed, reason:%+v", logId, err)
		return err
	}

	return nil
}

func haVipDisassociateEip(meta interface{}, havipId string, eip string) error {
	//associate eip
	logId := tccommon.GetLogId(tccommon.ContextNil)
	bindRequest := vpc.NewHaVipDisassociateAddressIpRequest()
	bindRequest.HaVipId = helper.String(havipId)
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		_, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().HaVipDisassociateAddressIp(bindRequest)
		if e != nil {
			return tccommon.RetryError(errors.WithStack(e))
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create HA VIP attachment failed, reason:%+v", logId, err)
		return err
	}

	statRequest := vpc.NewDescribeHaVipsRequest()
	statRequest.HaVipIds = []*string{&havipId}
	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().DescribeHaVips(statRequest)
		if e != nil {
			//when associated eip is in deleting process, delete ha vip may return unsupported operation error
			return tccommon.RetryError(errors.WithStack(e), VPCUnsupportedOperation)

		} else {
			//if not, quit
			if len(result.Response.HaVipSet) > 0 {
				if *result.Response.HaVipSet[0].AddressIp != "" {
					return resource.RetryableError(fmt.Errorf("Not unbinded yet, retry describing"))
				} else {
					return nil
				}
			}
			return resource.NonRetryableError(fmt.Errorf("describe error"))
		}
	})
	if err != nil {
		log.Printf("[CRITAL]%s describe HA VIP failed, reason:%+v", logId, err)
		return err
	}

	return nil
}
