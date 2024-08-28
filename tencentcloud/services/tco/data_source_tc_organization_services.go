package tco

import (
	"context"
	"strconv"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	organization "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/organization/v20210331"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudOrganizationServices() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudOrganizationServicesRead,
		Schema: map[string]*schema.Schema{
			"search_key": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Keyword for search by name.",
			},
			// computed
			"items": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Organization service list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Organization service ID. Note: This field may return null, indicating that no valid values can be obtained.",
						},
						"product_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Organization service product name. Note: This field may return null, indicating that no valid values can be obtained.",
						},
						"is_assign": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether to support delegation. Valid values: 1 (yes), 2 (no). Note: This field may return null, indicating that no valid values can be obtained.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Organization service description. Note: This field may return null, indicating that no valid values can be obtained.",
						},
						"member_num": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Number of the current delegated admins. Note: This field may return null, indicating that no valid values can be obtained.",
						},
						"document": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Help documentation. Note: This field may return null, indicating that no valid values can be obtained.",
						},
						"console_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Console path of the organization service product. Note: This field may return null, indicating that no valid values can be obtained.",
						},
						"is_usage_status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether to access the usage status. Valid values: 1 (yes), 2 (no). Note: This field may return null, indicating that no valid values can be obtained.",
						},
						"can_assign_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Limit for the number of delegated admins. Note: This field may return null, indicating that no valid values can be obtained.",
						},
						"product": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Organization service product identifier. Note: This field may return null, indicating that no valid values can be obtained.",
						},
						"service_grant": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether to support organization service authorization. Valid values: 1 (yes), 2 (no). Note: This field may return null, indicating that no valid values can be obtained.",
						},
						"grant_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Enabling status of organization service authorization. This field is valid when ServiceGrant is 1. Valid values: Enabled, Disabled. Note: This field may return null, indicating that no valid values can be obtained.",
						},
						"is_set_management_scope": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether to support setting the delegated management scope. Valid values: 1 (yes), 2 (no).\nNote: This field may return null, indicating that no valid values can be obtained.",
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

func dataSourceTencentCloudOrganizationServicesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_organization_services.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = OrganizationService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		items   []*organization.OrganizationServiceAssign
	)

	paramMap := make(map[string]interface{})

	if v, ok := d.GetOk("search_key"); ok {
		paramMap["SearchKey"] = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeOrganizationServicesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		items = result
		return nil
	})

	if err != nil {
		return err
	}

	ids := make([]string, 0, len(items))
	tmpList := make([]map[string]interface{}, 0, len(items))

	if items != nil {
		for _, item := range items {
			orgServiceAssignMap := map[string]interface{}{}

			if item.ServiceId != nil {
				orgServiceAssignMap["service_id"] = item.ServiceId
				serviceIdStr := strconv.FormatUint(*item.ServiceId, 10)
				ids = append(ids, serviceIdStr)
			}

			if item.ProductName != nil {
				orgServiceAssignMap["product_name"] = item.ProductName
			}

			if item.IsAssign != nil {
				orgServiceAssignMap["is_assign"] = item.IsAssign
			}

			if item.Description != nil {
				orgServiceAssignMap["description"] = item.Description
			}

			if item.MemberNum != nil {
				orgServiceAssignMap["member_num"] = item.MemberNum
			}

			if item.ConsoleUrl != nil {
				orgServiceAssignMap["console_url"] = item.ConsoleUrl
			}

			if item.IsUsageStatus != nil {
				orgServiceAssignMap["is_usage_status"] = item.IsUsageStatus
			}

			if item.CanAssignCount != nil {
				orgServiceAssignMap["can_assign_count"] = item.CanAssignCount
			}

			if item.Product != nil {
				orgServiceAssignMap["product"] = item.Product
			}

			if item.ServiceGrant != nil {
				orgServiceAssignMap["service_grant"] = item.ServiceGrant
			}

			if item.GrantStatus != nil {
				orgServiceAssignMap["grant_status"] = item.GrantStatus
			}

			if item.IsSetManagementScope != nil {
				orgServiceAssignMap["is_set_management_scope"] = item.IsSetManagementScope
			}

			tmpList = append(tmpList, orgServiceAssignMap)
		}

		_ = d.Set("items", tmpList)
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
