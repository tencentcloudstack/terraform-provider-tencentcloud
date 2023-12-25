package dlc

import (
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dlc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dlc/v20210125"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudDlcDetachWorkGroupPolicyOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDlcDetachWorkGroupPolicyOperationCreate,
		Read:   resourceTencentCloudDlcDetachWorkGroupPolicyOperationRead,
		Delete: resourceTencentCloudDlcDetachWorkGroupPolicyOperationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"work_group_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Work group id.",
			},

			"policy_set": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				Description: "The set of policies to be bound.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"database": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Database name that requires authorization, fill in * to represent all databases under the current catalog. When the authorization type is administrator level, only * is allowed to be filled in. when the authorization type is data connection level, only blanks are allowed to be filled in. For other types, the database can be specified arbitrarily.",
						},
						"catalog": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "For the data source name that requires authorization, only * (representing all resources at this level) is supported under the administrator level; in the case of data source level and database level authentication, only COSDataCatalog or * is supported; in data table level authentication, it is possible Fill in the user-defined data source. If left blank, it defaults to DataLakeCatalog. note: If a user-defined data source is authenticated, the permissions that dlc can manage are a subset of the accounts provided by the user when accessing the data source.",
						},
						"table": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "For the table name that requires authorization, fill in * to represent all tables under the current database. when the authorization type is administrator level, only * is allowed to be filled in. when the authorization type is data connection level or database level, only blanks are allowed to be filled in. For other types, data tables can be specified arbitrarily.",
						},
						"operation": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Authorized permission operations provide different operations for different levels of authentication. administrator permissions: ALL, default is ALL if left blank; data connection level authentication: CREATE; database level authentication: ALL, CREATE, ALTER, DROP; data table permissions: ALL, SELECT, INSERT, ALTER, DELETE, DROP, UPDATE. note: under data table permissions, only SELECT operations are supported when the specified data source is not COSDataCatalog.",
						},
						"policy_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Authorization type, currently supports eight authorization types: ADMIN: Administrator level authentication DATASOURCE: data connection level authentication DATABASE: database level authentication TABLE: Table level authentication VIEW: view level authentication FUNCTION: Function level authentication COLUMN: Column level authentication ENGINE: Data engine authentication. if left blank, the default is administrator level authentication.",
						},
						"function": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "For the function name that requires authorization, fill in * to represent all functions under the current catalog. when the authorization type is administrator level, only * is allowed to be filled in. When the authorization type is data connection level, only blanks are allowed to be filled in. in other types, functions can be specified arbitrarily.",
						},
						"view": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "For views that require authorization, fill in * to represent all views under the current database. When the authorization type is administrator level, only * is allowed to be filled in. when the authorization type is data connection level or database level, only blanks are allowed to be filled in. for other types, the view can be specified arbitrarily.",
						},
						"column": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "For columns that require authorization, fill in * to represent all current columns. When the authorization type is administrator level, only * is allowed.",
						},
						"data_engine": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Data engines that require authorization, fill in * to represent all current engines. when the authorization type is administrator level, only * is allowed.",
						},
						"re_auth": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether the user can perform secondary authorization. when it is true, the authorized user can re-authorize the permissions obtained this time to other sub-users. default is false.",
						},
						"source": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Permission source, please leave it blank. USER: permissions come from the user itself; WORKGROUP: permissions come from the bound workgroup.",
						},
						"mode": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Authorization mode, please leave this parameter blank. COMMON: normal mode; SENIOR: advanced mode.",
						},
						"operator": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Operator, do not fill in the input parameters.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The time when the permission was created. Leave the input parameter blank.",
						},
						"source_id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The id of the workgroup to which the permission belongs. this value only exists when the source of the permission is a workgroup. that is, this field has a value only when the value of the Source field is WORKGROUP.",
						},
						"source_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The name of the workgroup to which the permission belongs. this value only exists when the source of the permission is a workgroup. that is, this field has a value only when the value of the source field is WORKGROUP.",
						},
						"id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Policy id.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudDlcDetachWorkGroupPolicyOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_detach_work_group_policy_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request     = dlc.NewDetachWorkGroupPolicyRequest()
		workGroupId string
	)
	if v, _ := d.GetOkExists("work_group_id"); v != nil {
		workGroupId = helper.IntToStr(v.(int))
		request.WorkGroupId = helper.IntInt64(v.(int))
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
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().DetachWorkGroupPolicy(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate dlc detachWorkGroupPolicyOperation failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(workGroupId)

	return resourceTencentCloudDlcDetachWorkGroupPolicyOperationRead(d, meta)
}

func resourceTencentCloudDlcDetachWorkGroupPolicyOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_detach_work_group_policy_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudDlcDetachWorkGroupPolicyOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_detach_work_group_policy_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
