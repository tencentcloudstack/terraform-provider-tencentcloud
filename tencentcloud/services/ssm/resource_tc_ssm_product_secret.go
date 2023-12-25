package ssm

import (
	"context"
	"fmt"
	"log"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ssm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssm/v20190923"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudSsmProductSecret() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSsmProductSecretCreate,
		Read:   resourceTencentCloudSsmProductSecretRead,
		Update: resourceTencentCloudSsmProductSecretUpdate,
		Delete: resourceTencentCloudSsmProductSecretDelete,
		Schema: map[string]*schema.Schema{
			"secret_name": {
				Required:    true,
				Type:        schema.TypeString,
				ForceNew:    true,
				Description: "Credential name, which must be unique in the same region. It can contain 128 bytes of letters, digits, hyphens, and underscores and must begin with a letter or digit.",
			},
			"user_name_prefix": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Prefix of the user account name, which is specified by you and can contain up to 8 characters.Supported character sets include:Digits: [0, 9].Lowercase letters: [a, z].Uppercase letters: [A, Z].Special symbols: underscore.The prefix must begin with a letter.",
			},
			"product_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Name of the Tencent Cloud service bound to the credential, such as `Mysql`, `Tdsql-mysql`, `Tdsql_C_Mysql`. you can use dataSource `tencentcloud_ssm_products` to query supported products.",
			},
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Tencent Cloud service instance ID.",
			},
			"description": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Description, which is used to describe the purpose in detail and can contain up to 2,048 bytes.",
			},
			"kms_key_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Specifies the KMS CMK that encrypts the credential. If this parameter is left empty, the CMK created by Secrets Manager by default will be used for encryption.You can also specify a custom KMS CMK created in the same region for encryption.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tags of secret.",
			},
			"domains": {
				Required:    true,
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Domain name of the account in the form of IP. You can enter `%`.",
			},
			"privileges_list": {
				Required:    true,
				Type:        schema.TypeList,
				Description: "List of permissions that need to be granted when the credential is bound to a Tencent Cloud service.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"privilege_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Permission name. Valid values: `GlobalPrivileges`, `DatabasePrivileges`, `TablePrivileges`, `ColumnPrivileges`. When the permission is `DatabasePrivileges`, the database name must be specified by the `Database` parameter; When the permission is `TablePrivileges`, the database name and the table name in the database must be specified by the `Database` and `TableName` parameters; When the permission is `ColumnPrivileges`, the database name, table name in the database, and column name in the table must be specified by the `Database`, `TableName`, and `ColumnName` parameters.",
						},
						"privileges": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Required:    true,
							Description: "Permission list. For the `Mysql` service, optional permission values are: 1. Valid values of `GlobalPrivileges`: SELECT,INSERT,UPDATE,DELETE,CREATE, PROCESS, DROP,REFERENCES,INDEX,ALTER,SHOW DATABASES,CREATE TEMPORARY TABLES,LOCK TABLES,EXECUTE,CREATE VIEW,SHOW VIEW,CREATE ROUTINE,ALTER ROUTINE,EVENT,TRIGGER. Note: if this parameter is not passed in, it means to clear the permission. 2. Valid values of `DatabasePrivileges`: SELECT,INSERT,UPDATE,DELETE,CREATE, DROP,REFERENCES,INDEX,ALTER,CREATE TEMPORARY TABLES,LOCK TABLES,EXECUTE,CREATE VIEW,SHOW VIEW,CREATE ROUTINE,ALTER ROUTINE,EVENT,TRIGGER. Note: if this parameter is not passed in, it means to clear the permission. 3. Valid values of `TablePrivileges`: SELECT,INSERT,UPDATE,DELETE,CREATE, DROP,REFERENCES,INDEX,ALTER,CREATE VIEW,SHOW VIEW, TRIGGER. Note: if this parameter is not passed in, it means to clear the permission. 4. Valid values of `ColumnPrivileges`: SELECT,INSERT,UPDATE,REFERENCES.Note: if this parameter is not passed in, it means to clear the permission.",
						},
						"database": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "This value takes effect only when `PrivilegeName` is `DatabasePrivileges`.",
						},
						"table_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "This value takes effect only when `PrivilegeName` is `TablePrivileges`, and the `Database` parameter is required in this case to explicitly indicate the database instance.",
						},
						"column_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "This value takes effect only when `PrivilegeName` is `ColumnPrivileges`, and the following parameters are required in this case:Database: explicitly indicate the database instance.TableName: explicitly indicate the table.",
						},
					},
				},
			},
			"rotation_begin_time": {
				Optional:    true,
				Type:        schema.TypeString,
				Computed:    true,
				Description: "User-Defined rotation start time in the format of 2006-01-02 15:04:05.When `EnableRotation` is `True`, this parameter is required.",
			},
			"enable_rotation": {
				Optional:    true,
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Specifies whether to enable rotation, when secret status is `Disabled`, rotation will be disabled. `True` - enable, `False` - do not enable. If this parameter is not specified, `False` will be used by default.",
			},
			"rotation_frequency": {
				Optional:    true,
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Rotation frequency in days. Default value: 1 day.",
			},
			"status": {
				Optional:     true,
				Type:         schema.TypeString,
				Computed:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{"Enabled", "Disabled"}),
				Description:  "Enable or Disable Secret. Valid values is `Enabled` or `Disabled`. Default is `Enabled`.",
			},
			"create_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Credential creation time in UNIX timestamp format.",
			},
			"secret_type": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "`0`: user-defined secret. `1`: Tencent Cloud services secret. `2`: SSH key secret. `3`: Tencent Cloud API key secret. Note: this field may return `null`, indicating that no valid values can be obtained.",
			},
		},
	}
}

func resourceTencentCloudSsmProductSecretCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ssm_product_secret.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = SsmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		request    = ssm.NewCreateProductSecretRequest()
		response   = ssm.NewCreateProductSecretResponse()
		secretInfo *SecretInfo
		secretName string
	)

	if v, ok := d.GetOk("secret_name"); ok {
		request.SecretName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("user_name_prefix"); ok {
		request.UserNamePrefix = helper.String(v.(string))
	}

	if v, ok := d.GetOk("product_name"); ok {
		request.ProductName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceID = helper.String(v.(string))
	}

	if v, ok := d.GetOk("domains"); ok {
		domainsSet := v.(*schema.Set).List()
		for i := range domainsSet {
			domains := domainsSet[i].(string)
			request.Domains = append(request.Domains, &domains)
		}
	}

	if v, ok := d.GetOk("privileges_list"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			productPrivilegeUnit := ssm.ProductPrivilegeUnit{}
			if v, ok := dMap["privilege_name"]; ok {
				productPrivilegeUnit.PrivilegeName = helper.String(v.(string))
			}
			if v, ok := dMap["privileges"]; ok {
				privilegesSet := v.(*schema.Set).List()
				for i := range privilegesSet {
					privileges := privilegesSet[i].(string)
					productPrivilegeUnit.Privileges = append(productPrivilegeUnit.Privileges, &privileges)
				}
			}
			if v, ok := dMap["database"]; ok {
				productPrivilegeUnit.Database = helper.String(v.(string))
			}
			if v, ok := dMap["table_name"]; ok {
				productPrivilegeUnit.TableName = helper.String(v.(string))
			}
			if v, ok := dMap["column_name"]; ok {
				productPrivilegeUnit.ColumnName = helper.String(v.(string))
			}
			request.PrivilegesList = append(request.PrivilegesList, &productPrivilegeUnit)
		}
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOk("kms_key_id"); ok {
		request.KmsKeyId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("rotation_begin_time"); ok {
		request.RotationBeginTime = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("enable_rotation"); ok {
		request.EnableRotation = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("rotation_frequency"); ok {
		request.RotationFrequency = helper.IntInt64(v.(int))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseSsmClient().CreateProductSecret(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create ssm productSecret failed, reason:%+v", logId, err)
		return err
	}

	secretName = *response.Response.SecretName
	d.SetId(secretName)
	flowId := *response.Response.FlowID
	conf := tccommon.BuildStateChangeConf([]string{}, []string{"1"}, tccommon.ReadRetryTimeout, time.Second, service.SsmProductSecretStateRefreshFunc(flowId, []string{"0"}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	// update status if disabled
	if v, ok := d.GetOk("status"); ok {
		status := v.(string)
		if status == "Disabled" {
			ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
			err = service.DisableSecret(ctx, secretName)
			if err != nil {
				return err
			}
		}
	}

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			secretInfo, err = service.DescribeSecretByName(ctx, secretName)
			if err != nil {
				return tccommon.RetryError(err)
			}

			return nil
		})

		if err != nil {
			return err
		}

		tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		tagService := svctag.NewTagService(tcClient)
		resourceName := tccommon.BuildTagResourceName("ssm", "secret", tcClient.Region, secretInfo.resourceId)
		if err = tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudSsmProductSecretRead(d, meta)
}

func resourceTencentCloudSsmProductSecretRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ssm_product_secret.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = SsmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		secretInfo *SecretInfo
		secretName = d.Id()
	)

	productSecret, err := service.DescribeSecretById(ctx, secretName, 1)
	if err != nil {
		return err
	}

	if productSecret == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `SsmProductSecret` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if productSecret.SecretName != nil {
		_ = d.Set("secret_name", productSecret.SecretName)
	}

	if productSecret.ProductName != nil {
		_ = d.Set("product_name", productSecret.ProductName)
	}

	if productSecret.ResourceID != nil {
		_ = d.Set("instance_id", productSecret.ResourceID)
	}

	if productSecret.Description != nil {
		_ = d.Set("description", productSecret.Description)
	}

	if productSecret.KmsKeyId != nil {
		_ = d.Set("kms_key_id", productSecret.KmsKeyId)
	}

	if productSecret.RotationBeginTime != nil {
		_ = d.Set("rotation_begin_time", productSecret.RotationBeginTime)
	}

	if productSecret.RotationStatus != nil {
		_ = d.Set("enable_rotation", helper.Bool(true))
		if *productSecret.RotationStatus == 0 {
			_ = d.Set("enable_rotation", helper.Bool(false))
		}
	}

	if productSecret.RotationFrequency != nil {
		_ = d.Set("rotation_frequency", productSecret.RotationFrequency)
	}

	if productSecret.CreateTime != nil {
		_ = d.Set("create_time", productSecret.CreateTime)
	}

	if productSecret.SecretType != nil {
		_ = d.Set("secret_type", productSecret.SecretType)
	}

	outErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		secretInfo, err = service.DescribeSecretByName(ctx, secretName)
		if err != nil {
			return tccommon.RetryError(err)
		}

		return nil
	})

	if outErr != nil {
		return outErr
	}

	tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	tagService := svctag.NewTagService(tcClient)
	tags, err := tagService.DescribeResourceTags(ctx, "ssm", "secret", tcClient.Region, secretInfo.resourceId)
	if err != nil {
		return err
	}

	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudSsmProductSecretUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ssm_product_secret.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		ssmService = SsmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		secretName = d.Id()
	)

	immutableArgs := []string{
		"user_name_prefix", "product_name", "instance_id",
		"domains", "privileges_list", "kms_key_id",
	}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("description") {
		request := ssm.NewUpdateDescriptionRequest()
		request.SecretName = &secretName

		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseSsmClient().UpdateDescription(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update ssm productSecret failed, reason:%+v", logId, err)
			return err
		}
	}

	if d.HasChange("status") {
		service := SsmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

		if v, ok := d.GetOk("status"); ok {
			status := v.(string)
			if status == "Disabled" {
				err := service.DisableSecret(ctx, secretName)
				if err != nil {
					return err
				}
			} else {
				err := service.EnableSecret(ctx, secretName)
				if err != nil {
					return err
				}
			}
		}
	}

	if d.HasChange("enable_rotation") || d.HasChange("rotation_begin_time") || d.HasChange("rotation_frequency") {
		request := ssm.NewUpdateRotationStatusRequest()
		request.SecretName = &secretName

		if v, ok := d.GetOk("rotation_begin_time"); ok {
			request.RotationBeginTime = helper.String(v.(string))
		}

		if v, ok := d.GetOkExists("enable_rotation"); ok {
			request.EnableRotation = helper.Bool(v.(bool))
		}

		if v, ok := d.GetOkExists("rotation_frequency"); ok {
			request.Frequency = helper.IntInt64(v.(int))
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseSsmClient().UpdateRotationStatus(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update ssm productSecret failed, reason:%+v", logId, err)
			return err
		}
	}

	if d.HasChange("tags") {
		tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		tagService := svctag.NewTagService(tcClient)

		oldValue, newValue := d.GetChange("tags")
		replaceTags, deleteTags := svctag.DiffTags(oldValue.(map[string]interface{}), newValue.(map[string]interface{}))
		secretInfo, err := ssmService.DescribeSecretByName(ctx, secretName)
		if err != nil {
			return err
		}

		resourceName := tccommon.BuildTagResourceName("ssm", "secret", tcClient.Region, secretInfo.resourceId)
		if err = tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}

	}

	return resourceTencentCloudSsmProductSecretRead(d, meta)
}

func resourceTencentCloudSsmProductSecretDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ssm_product_secret.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = SsmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		secretName = d.Id()
	)

	// disable before destroy
	err := service.DisableSecret(ctx, secretName)
	if err != nil {
		return err
	}

	if err = service.DeleteSsmProductSecretById(ctx, secretName); err != nil {
		return err
	}

	return nil
}
