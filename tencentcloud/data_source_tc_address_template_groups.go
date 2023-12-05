package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudAddressTemplateGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudAddressTemplateGroupsRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the address template group to query.",
			},
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Id of the address template group to query.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},

			// Computed values
			"group_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Information list of the dedicated address template groups.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Id of the address template group.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of address template group.",
						},
						"template_ids": {
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "ID set of the address template.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudAddressTemplateGroupsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_address_template_groups.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var name, templateId string
	var filters = make([]*vpc.Filter, 0)
	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
		filters = append(filters, &vpc.Filter{Name: helper.String("address-template-group-name"), Values: []*string{&name}})
	}

	if v, ok := d.GetOk("id"); ok {
		templateId = v.(string)
		filters = append(filters, &vpc.Filter{Name: helper.String("address-template-group-id"), Values: []*string{&templateId}})
	}

	vpcService := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}
	var outErr, inErr error
	groups, outErr := vpcService.DescribeAddressTemplateGroups(ctx, filters)
	if outErr != nil {
		outErr = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			groups, inErr = vpcService.DescribeAddressTemplateGroups(ctx, filters)
			if inErr != nil {
				return retryError(inErr)
			}
			return nil
		})
	}

	if outErr != nil {
		return outErr
	}

	ids := make([]string, 0, len(groups))
	templateGroupList := make([]map[string]interface{}, 0, len(groups))
	for _, ins := range groups {
		mapping := map[string]interface{}{
			"id":           ins.AddressTemplateGroupId,
			"name":         ins.AddressTemplateGroupName,
			"template_ids": ins.AddressTemplateIdSet,
		}
		templateGroupList = append(templateGroupList, mapping)
		ids = append(ids, *ins.AddressTemplateGroupId)
	}
	d.SetId(helper.DataResourceIdsHash(ids))
	if e := d.Set("group_list", templateGroupList); e != nil {
		log.Printf("[CRITAL]%s provider set address template group list fail, reason:%s\n", logId, e)
		return e
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), templateGroupList); e != nil {
			return e
		}
	}

	return nil

}
