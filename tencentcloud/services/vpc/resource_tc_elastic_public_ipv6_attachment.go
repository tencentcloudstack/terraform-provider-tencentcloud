package vpc

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudElasticPublicIpv6Attachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudElasticPublicIpv6AttachmentCreate,
		Read:   resourceTencentCloudElasticPublicIpv6AttachmentRead,
		Delete: resourceTencentCloudElasticPublicIpv6AttachmentDelete,
		Update: resourceTencentCloudElasticPublicIpv6AttachmentUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"ipv6_address_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Elastic IPv6 unique ID, EIPv6 unique ID is like eipv6-11112222.",
			},

			"network_interface_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Elastic Network Interface ID to bind. Elastic Network Interface ID is like eni-11112222. NetworkInterfaceId and InstanceId cannot be specified simultaneously. The Elastic Network Interface ID can be queried by logging in to the console, or obtained through the networkInterfaceId in the return value of the DescribeNetworkInterfaces interface.",
			},

			"private_ipv6_address": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The intranet IPv6 to bind. If NetworkInterfaceId is specified, PrivateIPv6Address must also be specified, which means that the EIP is bound to the specified private network IP of the specified Elastic Network Interface. Also ensure that the specified PrivateIPv6Address is an intranet IPv6 on the specified NetworkInterfaceId. The intranet IPv6 of the specified Elastic Network Interface can be queried by logging in to the console, or obtained through the Ipv6AddressSet.Address in the return value of the DescribeNetworkInterfaces interface.",
			},

			"keep_bind_with_eni": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to keep the Elastic Network Interface bound when unbinding.",
			},
		},
	}
}

func resourceTencentCloudElasticPublicIpv6AttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_elastic_public_ipv6_attachment.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	var (
		ipId string
	)
	var (
		request  = vpc.NewAssociateIPv6AddressRequest()
		response = vpc.NewAssociateIPv6AddressResponse()
	)

	if v, ok := d.GetOk("ipv6_address_id"); ok {
		ipId = v.(string)
		request.IPv6AddressId = helper.String(ipId)
	}

	if v, ok := d.GetOk("network_interface_id"); ok {
		request.NetworkInterfaceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("private_ipv6_address"); ok {
		request.PrivateIPv6Address = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().AssociateIPv6AddressWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create elastic public ipv6 attachment failed, reason:%+v", logId, err)
		return err
	}

	_ = response

	d.SetId(ipId)

	return resourceTencentCloudElasticPublicIpv6AttachmentRead(d, meta)
}

func resourceTencentCloudElasticPublicIpv6AttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_elastic_public_ipv6_attachment.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	ipId := d.Id()

	respData, err := service.DescribeElasticPublicIpv6AttachmentById(ctx, ipId)
	if err != nil {
		return err
	}

	if respData == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `elastic_public_ipv6_attachment` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if len(respData.AddressSet) > 0 {
		address := respData.AddressSet[0]
		_ = d.Set("ipv6_address_id", address.AddressId)
		_ = d.Set("network_interface_id", address.NetworkInterfaceId)
		_ = d.Set("private_ipv6_address", address.PrivateAddressIp)
	}
	return nil
}

func resourceTencentCloudElasticPublicIpv6AttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_elastic_public_ipv6_attachment.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	ipId := d.Id()

	var (
		request  = vpc.NewDisassociateIPv6AddressRequest()
		response = vpc.NewDisassociateIPv6AddressResponse()
	)

	request.IPv6AddressId = helper.String(ipId)

	if v, ok := d.GetOkExists("keep_bind_with_eni"); ok {
		request.KeepBindWithEni = helper.Bool(v.(bool))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().DisassociateIPv6AddressWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete elastic public ipv6 attachment failed, reason:%+v", logId, err)
		return err
	}

	_ = response
	return nil
}

func resourceTencentCloudElasticPublicIpv6AttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_elastic_public_ipv6_attachment.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
