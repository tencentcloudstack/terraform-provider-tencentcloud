package vpc

import (
	"context"
	"fmt"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceTencentCloudVpc() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "This data source has been deprecated in Terraform TencentCloud provider version 1.10.0. Please use 'tencentcloud_vpc_instances' instead.",
		Read:               dataSourceTencentCloudVpcRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The ID of the specific VPC to retrieve.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the specific VPC to retrieve.",
			},
			"cidr_block": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The CIDR block of the VPC.",
			},
			"is_default": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether or not the default VPC.",
			},
			"is_multicast": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether or not the VPC has Multicast support.",
			},
		},
	}
}

func dataSourceTencentCloudVpcRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_vpc.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var (
		vpcId string
		name  string
	)
	if temp, ok := d.GetOk("id"); ok {
		tempStr := temp.(string)
		if tempStr != "" {
			vpcId = tempStr
		}
	}
	if temp, ok := d.GetOk("name"); ok {
		tempStr := temp.(string)
		if tempStr != "" {
			name = tempStr
		}
	}

	var vpcInfos, err = service.DescribeVpcs(ctx, vpcId, name, map[string]string{}, nil, "", "")
	if err != nil {
		return err
	}

	if len(vpcInfos) == 0 {
		return fmt.Errorf("vpc id=%s, name=%s not found", vpcId, name)
	}

	vpc := vpcInfos[0]
	d.SetId(vpc.vpcId)
	_ = d.Set("name", vpc.name)
	_ = d.Set("cidr_block", vpc.cidr)
	_ = d.Set("is_default", vpc.isDefault)
	_ = d.Set("is_multicast", vpc.isMulticast)

	return nil
}
