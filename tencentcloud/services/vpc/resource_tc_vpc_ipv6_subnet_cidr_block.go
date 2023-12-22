package vpc

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudVpcIpv6SubnetCidrBlock() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVpcIpv6SubnetCidrBlockCreate,
		Read:   resourceTencentCloudVpcIpv6SubnetCidrBlockRead,
		Delete: resourceTencentCloudVpcIpv6SubnetCidrBlockDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Required:    true,
				Type:        schema.TypeString,
				ForceNew:    true,
				Description: "The private network `ID` where the subnet is located. Such as:`vpc-f49l6u0z`.",
			},

			"ipv6_subnet_cidr_blocks": {
				Required:    true,
				Type:        schema.TypeList,
				ForceNew:    true,
				MaxItems:    1,
				Description: "Allocate a list of `IPv6` subnets.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subnet_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Subnet instance `ID`. Such as:`subnet-pxir56ns`.",
						},
						"ipv6_cidr_block": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "`IPv6` subnet segment. Such as: `3402:4e00:20:1001::/64`.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudVpcIpv6SubnetCidrBlockCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_ipv6_subnet_cidr_block.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request  = vpc.NewAssignIpv6SubnetCidrBlockRequest()
		vpcId    string
		subnetId string
	)
	if v, ok := d.GetOk("vpc_id"); ok {
		vpcId = v.(string)
		request.VpcId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("ipv6_subnet_cidr_blocks"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			ipv6SubnetCidrBlock := vpc.Ipv6SubnetCidrBlock{}
			if v, ok := dMap["subnet_id"]; ok {
				subnetId = v.(string)
				ipv6SubnetCidrBlock.SubnetId = helper.String(v.(string))
			}
			if v, ok := dMap["ipv6_cidr_block"]; ok {
				ipv6SubnetCidrBlock.Ipv6CidrBlock = helper.String(v.(string))
			}
			request.Ipv6SubnetCidrBlocks = append(request.Ipv6SubnetCidrBlocks, &ipv6SubnetCidrBlock)
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().AssignIpv6SubnetCidrBlock(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create vpc ipv6SubnetCidrBlock failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(vpcId + tccommon.FILED_SP + subnetId)

	return resourceTencentCloudVpcIpv6SubnetCidrBlockRead(d, meta)
}

func resourceTencentCloudVpcIpv6SubnetCidrBlockRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_ipv6_subnet_cidr_block.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	subnetId := idSplit[1]

	ipv6SubnetCidrBlock, err := service.DescribeSubnetById(ctx, subnetId)
	if err != nil {
		return err
	}

	if ipv6SubnetCidrBlock == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `VpcIpv6SubnetCidrBlock` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if ipv6SubnetCidrBlock.VpcId != nil {
		_ = d.Set("vpc_id", ipv6SubnetCidrBlock.VpcId)
	}

	if ipv6SubnetCidrBlock.Ipv6CidrBlock != nil {
		ipv6SubnetCidrBlocksList := []interface{}{}
		ipv6SubnetCidrBlocksMap := map[string]interface{}{}
		ipv6SubnetCidrBlocksMap["subnet_id"] = &subnetId
		ipv6SubnetCidrBlocksMap["ipv6_cidr_block"] = ipv6SubnetCidrBlock.Ipv6CidrBlock
		ipv6SubnetCidrBlocksList = append(ipv6SubnetCidrBlocksList, ipv6SubnetCidrBlocksMap)
		_ = d.Set("ipv6_subnet_cidr_blocks", ipv6SubnetCidrBlocksList)
	}

	return nil
}

func resourceTencentCloudVpcIpv6SubnetCidrBlockDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_ipv6_subnet_cidr_block.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	vpcId := idSplit[0]
	subnetId := idSplit[1]

	if err := service.DeleteVpcIpv6SubnetCidrBlockById(ctx, vpcId, subnetId); err != nil {
		return err
	}

	return nil
}
