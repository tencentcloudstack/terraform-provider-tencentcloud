package trocket

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tdmqRocketmq "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq/v20200217"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudTdmqRocketmqRole() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTdmqRocketmqRoleRead,
		Schema: map[string]*schema.Schema{
			"role_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Fuzzy query by role name.",
			},

			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Cluster ID (required).",
			},

			"role_sets": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Array of roles.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"role_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Role name.",
						},
						"token": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Value of the role token.",
						},
						"remark": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Remarks.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time.",
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Update time.",
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

func dataSourceTencentCloudTdmqRocketmqRoleRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_tdmqRocketmq_role.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("role_name"); ok {
		paramMap["role_name"] = v.(string)
	}

	if v, ok := d.GetOk("cluster_id"); ok {
		paramMap["cluster_id"] = v.(string)
	}

	tdmqRocketmqService := TdmqRocketmqService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var roleSets []*tdmqRocketmq.Role
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		results, e := tdmqRocketmqService.DescribeTdmqRocketmqRoleByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		roleSets = results
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read TdmqRocketmq roleSets failed, reason:%+v", logId, err)
		return err
	}

	roleSetList := []interface{}{}
	ids := make([]string, 0)
	for _, roleSet := range roleSets {
		roleSetMap := map[string]interface{}{}
		ids = append(ids, *roleSet.RoleName)
		roleSetMap["role_name"] = roleSet.RoleName
		if roleSet.Token != nil {
			roleSetMap["token"] = roleSet.Token
		}
		if roleSet.Remark != nil {
			roleSetMap["remark"] = roleSet.Remark
		}
		if roleSet.CreateTime != nil {
			roleSetMap["create_time"] = roleSet.CreateTime
		}
		if roleSet.UpdateTime != nil {
			roleSetMap["update_time"] = roleSet.UpdateTime
		}

		roleSetList = append(roleSetList, roleSetMap)
	}
	d.SetId(helper.DataResourceIdsHash(ids))
	_ = d.Set("role_sets", roleSetList)

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), roleSetList); e != nil {
			return e
		}
	}

	return nil
}
