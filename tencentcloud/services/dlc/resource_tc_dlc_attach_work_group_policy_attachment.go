package dlc

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dlc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dlc/v20210125"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudDlcAttachWorkGroupPolicyAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDlcAttachWorkGroupPolicyAttachmentCreate,
		Read:   resourceTencentCloudDlcAttachWorkGroupPolicyAttachmentRead,
		Delete: resourceTencentCloudDlcAttachWorkGroupPolicyAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"work_group_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Work group ID.",
			},

			"policy_set": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Collection of policies to be bound.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"database": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name of the target database. `*` represents all databases in the current catalog. To grant admin permissions, it must be `*`; to grant data connection permissions, it must be null; to grant other permissions, it can be any database.",
						},
						"catalog": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name of the target data source. To grant admin permission, it must be `*` (all resources at this level); to grant data source and database permissions, it must be `COSDataCatalog` or `*`; to grant table permissions, it can be a custom data source; if it is left empty, `DataLakeCatalog` is used. Note: To grant permissions on a custom data source, the permissions that can be managed in the Data Lake Compute console are subsets of the account permissions granted when you connect the data source to the console.",
						},
						"table": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name of the target table. `*` represents all tables in the current database. To grant admin permissions, it must be `*`; to grant data connection and database permissions, it must be null; to grant other permissions, it can be any table.",
						},
						"operation": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The target permissions, which vary by permission level. Admin: `ALL` (default); data connection: `CREATE`; database: `ALL`, `CREATE`, `ALTER`, and `DROP`; table: `ALL`, `SELECT`, `INSERT`, `ALTER`, `DELETE`, `DROP`, and `UPDATE`. Note: For table permissions, if a data source other than `COSDataCatalog` is specified, only the `SELECT` permission can be granted here.",
						},
						"policy_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The permission type. Valid values: `ADMIN`, `DATASOURCE`, `DATABASE`, `TABLE`, `VIEW`, `FUNCTION`, `COLUMN`, and `ENGINE`. Note: If it is left empty, `ADMIN` is used.",
						},
						"function": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The name of the target function. `*` represents all functions in the current catalog. To grant admin permissions, it must be `*`; to grant data connection permissions, it must be null; to grant other permissions, it can be any function.",
						},
						"view": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The name of the target view. `*` represents all views in the current database. To grant admin permissions, it must be `*`; to grant data connection and database permissions, it must be null; to grant other permissions, it can be any view.",
						},
						"column": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The name of the target column. `*` represents all columns. To grant admin permissions, it must be `*`.",
						},
						"data_engine": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The name of the target data engine. `*` represents all engines. To grant admin permissions, it must be `*`.",
						},
						"re_auth": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether the grantee is allowed to further grant the permissions. Valid values: `false` (default) and `true` (the grantee can grant permissions gained here to other sub-users).",
						},
						"source": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The permission source, which is not required when input parameters are passed in. Valid values: `USER` (from the user) and `WORKGROUP` (from one or more associated work groups).",
						},
						"mode": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The grant mode, which is not required as an input parameter. Valid values: `COMMON` and `SENIOR`.",
						},
						"operator": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The operator, which is not required as an input parameter.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The permission policy creation time, which is not required as an input parameter.",
						},
						"source_id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The ID of the work group, which applies only when the value of the `Source` field is `WORKGROUP`.",
						},
						"source_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The name of the work group, which applies only when the value of the `Source` field is `WORKGROUP`.",
						},
						"id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The policy ID.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudDlcAttachWorkGroupPolicyAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_attach_work_group_policy_attachment.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId       = tccommon.GetLogId(tccommon.ContextNil)
		request     = dlc.NewAttachWorkGroupPolicyRequest()
		workGroupId string
		policyId    string
	)

	if v, _ := d.GetOkExists("work_group_id"); v != nil {
		workGroupId = helper.IntToStr(v.(int))
		request.WorkGroupId = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("policy_set"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			policy := dlc.Policy{}
			if v, ok := dMap["database"]; ok {
				policy.Database = helper.String(v.(string))
			}

			if v, ok := dMap["catalog"]; ok {
				policy.Catalog = helper.String(v.(string))
			}

			if v, ok := dMap["table"]; ok {
				policy.Table = helper.String(v.(string))
			}

			if v, ok := dMap["operation"]; ok {
				policy.Operation = helper.String(v.(string))
			}

			if v, ok := dMap["policy_type"]; ok {
				policy.PolicyType = helper.String(v.(string))
			}

			if v, ok := dMap["function"]; ok {
				policy.Function = helper.String(v.(string))
			}

			if v, ok := dMap["view"]; ok {
				policy.View = helper.String(v.(string))
			}

			if v, ok := dMap["column"]; ok {
				policy.Column = helper.String(v.(string))
			}

			if v, ok := dMap["data_engine"]; ok {
				policy.DataEngine = helper.String(v.(string))
			}

			if v, ok := dMap["re_auth"]; ok {
				policy.ReAuth = helper.Bool(v.(bool))
			}

			if v, ok := dMap["source"]; ok {
				policy.Source = helper.String(v.(string))
			}

			if v, ok := dMap["mode"]; ok {
				policy.Mode = helper.String(v.(string))
			}

			if v, ok := dMap["operator"]; ok {
				policy.Operator = helper.String(v.(string))
			}

			if v, ok := dMap["create_time"]; ok {
				policy.CreateTime = helper.String(v.(string))
			}

			if v, ok := dMap["source_id"]; ok {
				policy.SourceId = helper.IntInt64(v.(int))
			}

			if v, ok := dMap["source_name"]; ok {
				policy.SourceName = helper.String(v.(string))
			}

			if v, ok := dMap["id"]; ok {
				policy.Id = helper.IntInt64(v.(int))
			}

			request.PolicySet = append(request.PolicySet, &policy)
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().AttachWorkGroupPolicy(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.PolicySet == nil || len(result.Response.PolicySet) == 0 {
			return resource.NonRetryableError(fmt.Errorf("create dlc attach work group policy attachment failed, Response is nil."))
		}

		if result.Response.PolicySet[0].PolicyId == nil || *result.Response.PolicySet[0].PolicyId == "" {
			log.Printf("[CRITAL]%s create dlc attach_work_group_policy_attachment, logId=%s, id=%s", logId, logId, d.Id())
			return resource.NonRetryableError(fmt.Errorf("create dlc attach_work_group_policy_attachment failed, PolicyId is empty."))
		}

		policyId = *result.Response.PolicySet[0].PolicyId
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create dlc attach_work_group_policy_attachment failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(workGroupId + tccommon.FILED_SP + policyId)
	return resourceTencentCloudDlcAttachWorkGroupPolicyAttachmentRead(d, meta)
}

func resourceTencentCloudDlcAttachWorkGroupPolicyAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_attach_work_group_policy_attachment.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = DlcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	workGroupId := idSplit[0]
	policyId := idSplit[1]

	policy, err := service.DescribeDlcWorkGroupPolicyAttachmentById(ctx, workGroupId, policyId)
	if err != nil {
		return err
	}

	if policy == nil {
		log.Printf("[CRUD] dlc_attach_work_group_policy_attachment id=%s", d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("work_group_id", helper.StrToInt64(workGroupId))

	if policy != nil {
		policySetMap := map[string]interface{}{}
		if policy.Database != nil {
			policySetMap["database"] = *policy.Database
		}

		if policy.Catalog != nil {
			policySetMap["catalog"] = *policy.Catalog
		}

		if policy.Table != nil {
			policySetMap["table"] = *policy.Table
		}

		if policy.Operation != nil {
			policySetMap["operation"] = *policy.Operation
		}

		if policy.PolicyType != nil {
			policySetMap["policy_type"] = *policy.PolicyType
		}

		if policy.Function != nil {
			policySetMap["function"] = *policy.Function
		}

		if policy.View != nil {
			policySetMap["view"] = *policy.View
		}

		if policy.Column != nil {
			policySetMap["column"] = *policy.Column
		}

		if policy.DataEngine != nil {
			policySetMap["data_engine"] = *policy.DataEngine
		}

		if policy.ReAuth != nil {
			policySetMap["re_auth"] = *policy.ReAuth
		}

		if policy.Source != nil {
			policySetMap["source"] = *policy.Source
		}

		if policy.Mode != nil {
			policySetMap["mode"] = *policy.Mode
		}

		if policy.Operator != nil {
			policySetMap["operator"] = *policy.Operator
		}

		if policy.CreateTime != nil {
			policySetMap["create_time"] = *policy.CreateTime
		}

		if policy.SourceId != nil {
			policySetMap["source_id"] = *policy.SourceId
		}

		if policy.SourceName != nil {
			policySetMap["source_name"] = *policy.SourceName
		}

		if policy.Id != nil {
			policySetMap["id"] = *policy.Id
		}

		_ = d.Set("policy_set", []interface{}{policySetMap})
	}

	return nil
}

func resourceTencentCloudDlcAttachWorkGroupPolicyAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_attach_work_group_policy_attachment.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = DlcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	workGroupId := idSplit[0]
	policyId := idSplit[1]

	if err := service.DeleteDlcAttachWorkGroupPolicyAttachmentByPolicyId(ctx, workGroupId, policyId); err != nil {
		return err
	}

	return nil
}
