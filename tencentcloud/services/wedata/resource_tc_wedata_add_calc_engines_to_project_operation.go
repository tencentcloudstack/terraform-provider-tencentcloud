package wedata

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wedatav20250806 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/wedata/v20250806"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudWedataAddCalcEnginesToProjectOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudWedataAddCalcEnginesToProjectOperationCreate,
		Read:   resourceTencentCloudWedataAddCalcEnginesToProjectOperationRead,
		Delete: resourceTencentCloudWedataAddCalcEnginesToProjectOperationDelete,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Project ID to be modified.",
			},

			"dlc_info": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: "DLC cluster information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"compute_resources": {
							Type:        schema.TypeSet,
							Required:    true,
							Description: "DLC resource names (need to add role Uin to DLC, otherwise resources may not be available).",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"region": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "DLC region.",
						},
						"default_database": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Specify the default database for the DLC cluster.",
						},
						"standard_mode_env_tag": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Cluster configuration tag (only effective for standard mode projects and required for standard mode). Enum values:\n- Prod  (Production environment)\n- Dev  (Development environment).",
						},
						"access_account": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Access account (only effective for standard mode projects and required for standard mode), used to submit DLC tasks.\nIt is recommended to use a specified sub-account and set corresponding database table permissions for the sub-account; task runner mode may cause task failures when the responsible person leaves; main account mode is not easy for permission control when multiple projects have different permissions.\n\nEnum values:\n- TASK_RUNNER (Task Runner)\n- OWNER (Main Account Mode)\n- SUB (Sub-Account Mode).",
						},
						"sub_account_uin": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Sub-account ID (only effective for standard mode projects), when AccessAccount is in sub-account mode, the sub-account ID information needs to be specified, other modes do not need to be specified.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudWedataAddCalcEnginesToProjectOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_add_calc_engines_to_project_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request   = wedatav20250806.NewAddCalcEnginesToProjectRequest()
		projectId string
	)

	if v, ok := d.GetOk("project_id"); ok {
		request.ProjectId = helper.String(v.(string))
		projectId = v.(string)
	}

	if v, ok := d.GetOk("dlc_info"); ok {
		for _, item := range v.([]interface{}) {
			dLCInfoMap := item.(map[string]interface{})
			dLCClusterInfo := wedatav20250806.DLCClusterInfo{}
			if v, ok := dLCInfoMap["compute_resources"]; ok {
				computeResourcesSet := v.(*schema.Set).List()
				for i := range computeResourcesSet {
					computeResources := computeResourcesSet[i].(string)
					dLCClusterInfo.ComputeResources = append(dLCClusterInfo.ComputeResources, helper.String(computeResources))
				}
			}

			if v, ok := dLCInfoMap["region"].(string); ok && v != "" {
				dLCClusterInfo.Region = helper.String(v)
			}

			if v, ok := dLCInfoMap["default_database"].(string); ok && v != "" {
				dLCClusterInfo.DefaultDatabase = helper.String(v)
			}

			if v, ok := dLCInfoMap["standard_mode_env_tag"].(string); ok && v != "" {
				dLCClusterInfo.StandardModeEnvTag = helper.String(v)
			}

			if v, ok := dLCInfoMap["access_account"].(string); ok && v != "" {
				dLCClusterInfo.AccessAccount = helper.String(v)
			}

			if v, ok := dLCInfoMap["sub_account_uin"].(string); ok && v != "" {
				dLCClusterInfo.SubAccountUin = helper.String(v)
			}

			request.DLCInfo = append(request.DLCInfo, &dLCClusterInfo)
		}
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWedataV20250806Client().AddCalcEnginesToProjectWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create wedata add calc engines to project operation failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	d.SetId(projectId)
	return resourceTencentCloudWedataAddCalcEnginesToProjectOperationRead(d, meta)
}

func resourceTencentCloudWedataAddCalcEnginesToProjectOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_add_calc_engines_to_project_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudWedataAddCalcEnginesToProjectOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_add_calc_engines_to_project_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
