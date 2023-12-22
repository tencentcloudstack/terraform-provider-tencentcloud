package vpc

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudVpcClassicLinkInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudVpcClassicLinkInstancesRead,
		Schema: map[string]*schema.Schema{
			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Filter conditions.`vpc-id` - String - (Filter condition) The VPC instance ID. `vm-ip` - String - (Filter condition) The IP address of the CVM on the basic network.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The attribute name. If more than one Filter exists, the logical relation between these Filters is `AND`.",
						},
						"values": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "The attribute value. If there are multiple Values for one Filter, the logical relation between these Values under the same Filter is `OR`.",
						},
					},
				},
			},

			"classic_link_instance_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Classiclink instance.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VPC instance ID.",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique ID of the CVM instance.",
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

func dataSourceTencentCloudVpcClassicLinkInstancesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_vpc_classic_link_instances.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*vpc.FilterObject, 0, len(filtersSet))

		for _, item := range filtersSet {
			filterObject := vpc.FilterObject{}
			filterObjectMap := item.(map[string]interface{})

			if v, ok := filterObjectMap["name"]; ok {
				filterObject.Name = helper.String(v.(string))
			}
			if v, ok := filterObjectMap["values"]; ok {
				valuesSet := v.(*schema.Set).List()
				filterObject.Values = helper.InterfacesStringsPoint(valuesSet)
			}
			tmpSet = append(tmpSet, &filterObject)
		}
		paramMap["Filters"] = tmpSet
	}

	service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var classicLinkInstanceSet []*vpc.ClassicLinkInstance

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeVpcClassicLinkInstancesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		classicLinkInstanceSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(classicLinkInstanceSet))
	tmpList := make([]map[string]interface{}, 0, len(classicLinkInstanceSet))

	if classicLinkInstanceSet != nil {
		for _, classicLinkInstance := range classicLinkInstanceSet {
			classicLinkInstanceMap := map[string]interface{}{}

			if classicLinkInstance.VpcId != nil {
				classicLinkInstanceMap["vpc_id"] = classicLinkInstance.VpcId
			}

			if classicLinkInstance.InstanceId != nil {
				classicLinkInstanceMap["instance_id"] = classicLinkInstance.InstanceId
			}

			ids = append(ids, *classicLinkInstance.InstanceId)
			tmpList = append(tmpList, classicLinkInstanceMap)
		}

		_ = d.Set("classic_link_instance_set", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
