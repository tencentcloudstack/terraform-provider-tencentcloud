package vpc

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudVpcIpv6CidrBlock() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVpcIpv6CidrBlockCreate,
		Read:   resourceTencentCloudVpcIpv6CidrBlockRead,
		Delete: resourceTencentCloudVpcIpv6CidrBlockDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Required:    true,
				Type:        schema.TypeString,
				ForceNew:    true,
				Description: "`VPC` instance `ID`, in the form of `vpc-f49l6u0z`.",
			},
			"ipv6_cidr_block": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ipv6 cidr block.",
			},
		},
	}
}

func resourceTencentCloudVpcIpv6CidrBlockCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_ipv6_cidr_block.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request = vpc.NewAssignIpv6CidrBlockRequest()
		vpcId   string
	)
	if v, ok := d.GetOk("vpc_id"); ok {
		vpcId = v.(string)
		request.VpcId = helper.String(vpcId)
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().AssignIpv6CidrBlock(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create vpc ipv6CidrBlock failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(vpcId)

	return resourceTencentCloudVpcIpv6CidrBlockRead(d, meta)
}

func resourceTencentCloudVpcIpv6CidrBlockRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_ipv6_cidr_block.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	vpcId := d.Id()

	instance, err := service.DescribeVpcById(ctx, vpcId)
	if err != nil {
		return err
	}

	if instance == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `VpcIpv6CidrBlock` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if instance.VpcId != nil {
		_ = d.Set("vpc_id", instance.VpcId)
	}

	if instance.Ipv6CidrBlock != nil {
		_ = d.Set("ipv6_cidr_block", instance.Ipv6CidrBlock)
	}

	return nil
}

func resourceTencentCloudVpcIpv6CidrBlockDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_ipv6_cidr_block.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	vpcId := d.Id()

	if err := service.DeleteVpcIpv6CidrBlockById(ctx, vpcId); err != nil {
		return err
	}

	return nil
}
