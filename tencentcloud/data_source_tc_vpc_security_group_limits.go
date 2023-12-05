package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudVpcSecurityGroupLimits() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudVpcSecurityGroupLimitsRead,
		Schema: map[string]*schema.Schema{
			"security_group_limit_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "sg limit set.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"security_group_limit": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "number of sg can be created.",
						},
						"security_group_policy_limit": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "number of sg polciy can be created.",
						},
						"referred_security_group_limit": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "number of sg can be referred.",
						},
						"security_group_instance_limit": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "number of sg associated instances.",
						},
						"instance_security_group_limit": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "number of instances associated sg.",
						},
						"security_group_extended_policy_limit": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "number of sg extended policy.",
						},
						"security_group_referred_cvm_and_eni_limit": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "number of eni and cvm can be referred.",
						},
						"security_group_referred_svc_limit": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "number of svc can be referred.",
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

func dataSourceTencentCloudVpcSecurityGroupLimitsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_vpc_security_group_limits.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	var securityGroupLimitSet *vpc.SecurityGroupLimitSet

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeVpcSecurityGroupLimits(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		securityGroupLimitSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0)
	securityGroupLimitSetMap := map[string]interface{}{}

	if securityGroupLimitSet != nil {
		if securityGroupLimitSet.SecurityGroupLimit != nil {
			securityGroupLimitSetMap["security_group_limit"] = securityGroupLimitSet.SecurityGroupLimit
		}

		if securityGroupLimitSet.SecurityGroupPolicyLimit != nil {
			securityGroupLimitSetMap["security_group_policy_limit"] = securityGroupLimitSet.SecurityGroupPolicyLimit
		}

		if securityGroupLimitSet.ReferedSecurityGroupLimit != nil {
			securityGroupLimitSetMap["referred_security_group_limit"] = securityGroupLimitSet.ReferedSecurityGroupLimit
		}

		if securityGroupLimitSet.SecurityGroupInstanceLimit != nil {
			securityGroupLimitSetMap["security_group_instance_limit"] = securityGroupLimitSet.SecurityGroupInstanceLimit
		}

		if securityGroupLimitSet.InstanceSecurityGroupLimit != nil {
			securityGroupLimitSetMap["instance_security_group_limit"] = securityGroupLimitSet.InstanceSecurityGroupLimit
		}

		if securityGroupLimitSet.SecurityGroupExtendedPolicyLimit != nil {
			securityGroupLimitSetMap["security_group_extended_policy_limit"] = securityGroupLimitSet.SecurityGroupExtendedPolicyLimit
		}

		if securityGroupLimitSet.SecurityGroupReferedCvmAndEniLimit != nil {
			securityGroupLimitSetMap["security_group_referred_cvm_and_eni_limit"] = securityGroupLimitSet.SecurityGroupReferedCvmAndEniLimit
		}

		if securityGroupLimitSet.SecurityGroupReferedSvcLimit != nil {
			securityGroupLimitSetMap["security_group_referred_svc_limit"] = securityGroupLimitSet.SecurityGroupReferedSvcLimit
		}

		ids = append(ids, helper.UInt64ToStr(*securityGroupLimitSet.SecurityGroupLimit))
		_ = d.Set("security_group_limit_set", []interface{}{securityGroupLimitSetMap})
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), securityGroupLimitSetMap); e != nil {
			return e
		}
	}
	return nil
}
