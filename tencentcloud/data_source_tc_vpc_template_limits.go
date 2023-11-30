package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudVpcTemplateLimits() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudVpcTemplateLimitsRead,
		Schema: map[string]*schema.Schema{
			"template_limit": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "template limit.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address_template_member_limit": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "address template member limit.",
						},
						"address_template_group_member_limit": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "address template group member limit.",
						},
						"service_template_member_limit": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "service template member limit.",
						},
						"service_template_group_member_limit": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "service template group member limit.",
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

func dataSourceTencentCloudVpcTemplateLimitsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_vpc_template_limits.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	var templateLimit *vpc.TemplateLimit

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeVpcTemplateLimits(ctx)
		if e != nil {
			return retryError(e)
		}
		templateLimit = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0)
	templateLimitMap := map[string]interface{}{}

	if templateLimit != nil {

		if templateLimit.AddressTemplateMemberLimit != nil {
			templateLimitMap["address_template_member_limit"] = templateLimit.AddressTemplateMemberLimit
		}

		if templateLimit.AddressTemplateGroupMemberLimit != nil {
			templateLimitMap["address_template_group_member_limit"] = templateLimit.AddressTemplateGroupMemberLimit
		}

		if templateLimit.ServiceTemplateMemberLimit != nil {
			templateLimitMap["service_template_member_limit"] = templateLimit.ServiceTemplateMemberLimit
		}

		if templateLimit.ServiceTemplateGroupMemberLimit != nil {
			templateLimitMap["service_template_group_member_limit"] = templateLimit.ServiceTemplateGroupMemberLimit
		}

		ids = append(ids, helper.UInt64ToStr(*templateLimit.AddressTemplateMemberLimit))
		_ = d.Set("template_limit", []interface{}{templateLimitMap})
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), templateLimitMap); e != nil {
			return e
		}
	}
	return nil
}
