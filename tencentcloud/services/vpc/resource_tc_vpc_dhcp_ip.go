package vpc

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudVpcDhcpIp() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVpcDhcpIpCreate,
		Read:   resourceTencentCloudVpcDhcpIpRead,
		Update: resourceTencentCloudVpcDhcpIpUpdate,
		Delete: resourceTencentCloudVpcDhcpIpDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The private network `ID`.",
			},

			"subnet_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Subnet `ID`.",
			},

			"dhcp_ip_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "`DhcpIp` name.",
			},
		},
	}
}

func resourceTencentCloudVpcDhcpIpCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_dhcp_ip.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request  = vpc.NewCreateDhcpIpRequest()
		response = vpc.NewCreateDhcpIpResponse()
		dhcpIpId string
	)
	if v, ok := d.GetOk("vpc_id"); ok {
		request.VpcId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("subnet_id"); ok {
		request.SubnetId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("dhcp_ip_name"); ok {
		request.DhcpIpName = helper.String(v.(string))
	}

	// 默认1
	request.SecondaryPrivateIpAddressCount = helper.IntUint64(1)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().CreateDhcpIp(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create vpc dhcpIp failed, reason:%+v", logId, err)
		return err
	}

	dhcpIpSet := response.Response.DhcpIpSet
	if len(dhcpIpSet) < 1 {
		return fmt.Errorf("create vpc dhcpIp failed.")
	}

	dhcpIpId = *dhcpIpSet[0].DhcpIpId
	d.SetId(dhcpIpId)

	return resourceTencentCloudVpcDhcpIpRead(d, meta)
}

func resourceTencentCloudVpcDhcpIpRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_dhcp_ip.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	dhcpIpId := d.Id()

	dhcpIp, err := service.DescribeVpcDhcpIpById(ctx, dhcpIpId)
	if err != nil {
		return err
	}

	if dhcpIp == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `VpcDhcpIp` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if dhcpIp.VpcId != nil {
		_ = d.Set("vpc_id", dhcpIp.VpcId)
	}

	if dhcpIp.SubnetId != nil {
		_ = d.Set("subnet_id", dhcpIp.SubnetId)
	}

	if dhcpIp.DhcpIpName != nil {
		_ = d.Set("dhcp_ip_name", dhcpIp.DhcpIpName)
	}

	return nil
}

func resourceTencentCloudVpcDhcpIpUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_dhcp_ip.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := vpc.NewModifyDhcpIpAttributeRequest()

	dhcpIpId := d.Id()

	request.DhcpIpId = &dhcpIpId

	immutableArgs := []string{"vpc_id", "subnet_id"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("dhcp_ip_name") {
		if v, ok := d.GetOk("dhcp_ip_name"); ok {
			request.DhcpIpName = helper.String(v.(string))
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().ModifyDhcpIpAttribute(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update vpc dhcpIp failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudVpcDhcpIpRead(d, meta)
}

func resourceTencentCloudVpcDhcpIpDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_dhcp_ip.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	dhcpIpId := d.Id()

	if err := service.DeleteVpcDhcpIpById(ctx, dhcpIpId); err != nil {
		return err
	}

	return nil
}
