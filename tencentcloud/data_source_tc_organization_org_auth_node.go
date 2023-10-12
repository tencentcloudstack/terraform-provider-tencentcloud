/*
Use this data source to query detailed information of organization org_auth_node

Example Usage

```hcl
data "tencentcloud_organization_org_auth_node" "org_auth_node" {
  }
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	organization "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/organization/v20210331"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudOrganizationOrgAuthNode() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudOrganizationOrgAuthNodeRead,
		Schema: map[string]*schema.Schema{
			"auth_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Verified company name.",
			},

			"items": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Organization auth node list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"relation_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Relationship Id.",
						},
						"auth_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Verified company name.",
						},
						"manager": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Organization auth manager.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"member_uin": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Member uin.",
									},
									"member_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Member name.",
									},
								},
							},
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

func dataSourceTencentCloudOrganizationOrgAuthNodeRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_organization_org_auth_node.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("auth_name"); ok {
		paramMap["AuthName"] = helper.String(v.(string))
	}

	service := OrganizationService{client: meta.(*TencentCloudClient).apiV3Conn}

	var items []*organization.AuthNode

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeOrganizationOrgAuthNodeByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
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
		for _, authNode := range items {
			authNodeMap := map[string]interface{}{}

			if authNode.RelationId != nil {
				authNodeMap["relation_id"] = authNode.RelationId
			}

			if authNode.AuthName != nil {
				authNodeMap["auth_name"] = authNode.AuthName
			}

			if authNode.Manager != nil {
				managerMap := map[string]interface{}{}

				if authNode.Manager.MemberUin != nil {
					managerMap["member_uin"] = authNode.Manager.MemberUin
				}

				if authNode.Manager.MemberName != nil {
					managerMap["member_name"] = authNode.Manager.MemberName
				}

				authNodeMap["manager"] = []interface{}{managerMap}
			}

			ids = append(ids, *authNode.AuthName)
			tmpList = append(tmpList, authNodeMap)
		}

		_ = d.Set("items", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output3, ok := d.GetOk("result_output_file")
	if ok && output3.(string) != "" {
		if e := writeToFile(output3.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
