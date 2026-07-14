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

func ResourceTencentCloudDlcAttachUserPolicyrAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDlcAttachUserPolicyrAttachmentCreate,
		Read:   resourceTencentCloudDlcAttachUserPolicyrAttachmentRead,
		Update: resourceTencentCloudDlcAttachUserPolicyrAttachmentUpdate,
		Delete: resourceTencentCloudDlcAttachUserPolicyrAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "User ID, which is the same as the sub-user UIN. The CreateUser API is needed to create a user at first. The DescribeUsers API can be used for viewing.",
			},

			"policy_set": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: "Collection of authentication policies.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"database": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The name of the target database. `*` represents all databases in the current catalog. To grant admin permissions, it must be `*`; to grant data connection permissions, it must be null; to grant other permissions, it can be any database.",
						},
						"catalog": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The name of the target data source. To grant admin permission, it must be `*` (all resources at this level); to grant data source and database permissions, it must be `COSDataCatalog` or `*`; to grant table permissions, it can be a custom data source; if it is left empty, `DataLakeCatalog` is used.",
						},
						"table": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The name of the target table. `*` represents all tables in the current database. To grant admin permissions, it must be `*`; to grant data connection and database permissions, it must be null; to grant other permissions, it can be any table.",
						},
						"operation": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The target permissions, which vary by permission level. Admin: `ALL` (default); data connection: `CREATE`; database: `ALL`, `CREATE`, `ALTER`, and `DROP`; table: `ALL`, `SELECT`, `INSERT`, `ALTER`, `DELETE`, `DROP`, and `UPDATE`.",
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
						"engine_generation": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Engine type.",
						},
						"model": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The name of the target Model. `*` represents all tables in the current database. To grant admin permissions, it must be `*`; to grant data connection and database permissions, it must be null; to grant other permissions, it can be any table.",
						},
						"policy_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The deterministic string PolicyId corresponding to the user and workgroup.",
						},
						"source": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The permission source, which is not required when input parameters are passed in. Valid values: `USER` (from the user) and `WORKGROUP` (from one or more associated work groups).",
						},
						"mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The grant mode, which is not required as an input parameter. Valid values: `COMMON` and `SENIOR`.",
						},
						"operator": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The operator, which is not required as an input parameter.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The permission policy creation time, which is not required as an input parameter.",
						},
						"source_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The ID of the work group, which applies only when the value of the `Source` field is `WORKGROUP`.",
						},
						"source_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the work group, which applies only when the value of the `Source` field is `WORKGROUP`.",
						},
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The policy ID.",
						},
						"is_admin_policy": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether the permission source is admin.",
						},
					},
				},
			},

			"account_type": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "User source type. Valid values: `TencentAccount` (common Tencent Cloud user) and `EntraAccount` (Microsoft user).",
			},
		},
	}
}

func resourceTencentCloudDlcAttachUserPolicyrAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_attach_user_policyr_attachment.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = dlc.NewAttachUserPolicyRequest()
		userId  string
	)

	if v, ok := d.GetOk("user_id"); ok {
		userId = v.(string)
		request.UserId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("account_type"); ok {
		request.AccountType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("policy_set"); ok {
		for _, item := range v.([]interface{}) {
			policy := dlc.Policy{}
			dMap := item.(map[string]interface{})
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

			if v, ok := dMap["engine_generation"]; ok {
				policy.EngineGeneration = helper.String(v.(string))
			}

			if v, ok := dMap["model"]; ok {
				policy.Model = helper.String(v.(string))
			}

			if v, ok := dMap["policy_id"]; ok {
				policy.PolicyId = helper.String(v.(string))
			}

			request.PolicySet = append(request.PolicySet, &policy)
		}
	}

	var response *dlc.AttachUserPolicyResponse
	if err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().AttachUserPolicyWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create dlc attach_user_policyr_attachment failed, Response is nil."))
		}

		if result.Response.PolicySet == nil {
			log.Printf("[CRITAL]%s dlc attach_user_policyr_attachment id=%s, response PolicySet is nil.", logId, userId)
			return resource.NonRetryableError(fmt.Errorf("Create dlc attach_user_policyr_attachment failed, PolicySet is nil."))
		}

		for _, policy := range result.Response.PolicySet {
			if policy == nil {
				continue
			}
			if policy.PolicyId == nil || *policy.PolicyId == "" {
				log.Printf("[CRITAL]%s dlc attach_user_policyr_attachment user_id=%s, policy PolicyId is empty.", logId, userId)
				return resource.NonRetryableError(fmt.Errorf("Create dlc attach_user_policyr_attachment failed, policy PolicyId is empty."))
			}
		}

		response = result
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s create dlc attach_user_policyr_attachment failed, reason:%+v", logId, err)
		return err
	}

	if response != nil && response.Response != nil && response.Response.PolicySet != nil {
		_ = d.Set("policy_set", flattenDlcAttachUserPolicyrAttachmentPolicySet(response.Response.PolicySet))
	}

	accountType := d.Get("account_type").(string)
	d.SetId(strings.Join([]string{userId, accountType}, tccommon.FILED_SP))
	return resourceTencentCloudDlcAttachUserPolicyrAttachmentRead(d, meta)
}

func resourceTencentCloudDlcAttachUserPolicyrAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_attach_user_policyr_attachment.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	userId := idSplit[0]
	accountType := idSplit[1]

	request := dlc.NewDescribeUserInfoRequest()
	request.UserId = helper.String(userId)
	request.Type = helper.String("DataAuth")
	if accountType != "" {
		request.AccountType = helper.String(accountType)
	}
	request.Limit = helper.IntInt64(100)

	var response *dlc.DescribeUserInfoResponse
	if err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().DescribeUserInfoWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe dlc attach_user_policyr_attachment failed, Response is nil."))
		}

		response = result
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s read dlc attach_user_policyr_attachment failed, reason:%+v", logId, err)
		return err
	}

	if response.Response.UserInfo == nil {
		log.Printf("[CRUD]%s dlc attach_user_policyr_attachment id=%s, UserInfo is nil.", logId, d.Id())
		d.SetId("")
		return nil
	}

	userInfo := response.Response.UserInfo
	if userInfo.DataPolicyInfo == nil || userInfo.DataPolicyInfo.PolicySet == nil || len(userInfo.DataPolicyInfo.PolicySet) == 0 {
		log.Printf("[CRUD]%s dlc attach_user_policyr_attachment id=%s, DataPolicyInfo is empty.", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("user_id", userId)
	if accountType != "" {
		_ = d.Set("account_type", accountType)
	}

	policySetList := flattenDlcAttachUserPolicyrAttachmentPolicySet(userInfo.DataPolicyInfo.PolicySet)
	_ = d.Set("policy_set", policySetList)

	return nil
}

func resourceTencentCloudDlcAttachUserPolicyrAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_attach_user_policyr_attachment.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	immutableArgs := []string{"user_id", "policy_set", "account_type"}
	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("the argument `%s` cannot be changed in-place for dlc attach_user_policyr_attachment, please recreate the resource.", v)
		}
	}

	return nil
}

func resourceTencentCloudDlcAttachUserPolicyrAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_attach_user_policyr_attachment.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	userId := idSplit[0]
	accountType := idSplit[1]

	request := dlc.NewDetachUserPolicyRequest()
	request.UserId = helper.String(userId)
	if accountType != "" {
		request.AccountType = helper.String(accountType)
	}

	if v, ok := d.GetOk("policy_set"); ok {
		for _, item := range v.([]interface{}) {
			policy := dlc.Policy{}
			dMap := item.(map[string]interface{})
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

			if v, ok := dMap["engine_generation"]; ok {
				policy.EngineGeneration = helper.String(v.(string))
			}

			if v, ok := dMap["model"]; ok {
				policy.Model = helper.String(v.(string))
			}

			if v, ok := dMap["policy_id"]; ok {
				policy.PolicyId = helper.String(v.(string))
			}

			request.PolicySet = append(request.PolicySet, &policy)
		}
	}

	if err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().DetachUserPolicyWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete dlc attach_user_policyr_attachment failed, Response is nil."))
		}

		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s delete dlc attach_user_policyr_attachment failed, reason:%+v", logId, err)
		return err
	}

	return nil
}

func flattenDlcAttachUserPolicyrAttachmentPolicySet(policySet []*dlc.Policy) []interface{} {
	policySetList := make([]interface{}, 0, len(policySet))
	for _, policy := range policySet {
		if policy == nil {
			continue
		}
		policyMap := map[string]interface{}{}
		if policy.Database != nil {
			policyMap["database"] = policy.Database
		}

		if policy.Catalog != nil {
			policyMap["catalog"] = policy.Catalog
		}

		if policy.Table != nil {
			policyMap["table"] = policy.Table
		}

		if policy.Operation != nil {
			policyMap["operation"] = policy.Operation
		}

		if policy.PolicyType != nil {
			policyMap["policy_type"] = policy.PolicyType
		}

		if policy.Function != nil {
			policyMap["function"] = policy.Function
		}

		if policy.View != nil {
			policyMap["view"] = policy.View
		}

		if policy.Column != nil {
			policyMap["column"] = policy.Column
		}

		if policy.DataEngine != nil {
			policyMap["data_engine"] = policy.DataEngine
		}

		if policy.ReAuth != nil {
			policyMap["re_auth"] = policy.ReAuth
		}

		if policy.EngineGeneration != nil {
			policyMap["engine_generation"] = policy.EngineGeneration
		}

		if policy.Model != nil {
			policyMap["model"] = policy.Model
		}

		if policy.PolicyId != nil {
			policyMap["policy_id"] = policy.PolicyId
		}

		if policy.Source != nil {
			policyMap["source"] = policy.Source
		}

		if policy.Mode != nil {
			policyMap["mode"] = policy.Mode
		}

		if policy.Operator != nil {
			policyMap["operator"] = policy.Operator
		}

		if policy.CreateTime != nil {
			policyMap["create_time"] = policy.CreateTime
		}

		if policy.SourceId != nil {
			policyMap["source_id"] = policy.SourceId
		}

		if policy.SourceName != nil {
			policyMap["source_name"] = policy.SourceName
		}

		if policy.Id != nil {
			policyMap["id"] = policy.Id
		}

		if policy.IsAdminPolicy != nil {
			policyMap["is_admin_policy"] = policy.IsAdminPolicy
		}

		policySetList = append(policySetList, policyMap)
	}

	return policySetList
}
