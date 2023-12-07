package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudVpcSecurityGroupReferences() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudVpcSecurityGroupReferencesRead,
		Schema: map[string]*schema.Schema{
			"security_group_ids": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "A set of security group instance IDs, e.g. [sg-12345678].",
			},

			"referred_security_group_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Referred security groups.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"security_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Security group instance ID.",
						},
						"referred_security_group_ids": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "IDs of all referred security group instances.",
						},
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudVpcSecurityGroupReferencesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_vpc_security_group_references.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("security_group_ids"); ok {
		securityGroupIdsSet := v.(*schema.Set).List()
		paramMap["SecurityGroupIds"] = helper.InterfacesStringsPoint(securityGroupIdsSet)
	}

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	var referredSecurityGroupSet []*vpc.ReferredSecurityGroup

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeVpcSecurityGroupReferences(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		referredSecurityGroupSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(referredSecurityGroupSet))
	tmpList := make([]map[string]interface{}, 0, len(referredSecurityGroupSet))

	if referredSecurityGroupSet != nil {
		for _, referredSecurityGroup := range referredSecurityGroupSet {
			referredSecurityGroupMap := map[string]interface{}{}

			if referredSecurityGroup.SecurityGroupId != nil {
				referredSecurityGroupMap["security_group_id"] = referredSecurityGroup.SecurityGroupId
			}

			if referredSecurityGroup.ReferredSecurityGroupIds != nil {
				referredSecurityGroupMap["referred_security_group_ids"] = referredSecurityGroup.ReferredSecurityGroupIds
			}

			ids = append(ids, *referredSecurityGroup.SecurityGroupId)
			tmpList = append(tmpList, referredSecurityGroupMap)
		}

		_ = d.Set("referred_security_group_set", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
