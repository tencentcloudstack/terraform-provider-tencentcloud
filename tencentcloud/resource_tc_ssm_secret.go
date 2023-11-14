/*
Provides a resource to create a ssm secret

Example Usage

```hcl
resource "tencentcloud_ssm_secret" "secret" {
  secret_name = "test_name"
  user_name_prefix = "test_prefix"
  enable_rotation = False
  rotation_begin_time = "2006-01-02 15:04:05"
  rotation_frequency = 1
  instance_i_d = "cdb-xxxxxxxx"
  description = ""
  kms_key_id = ""
  product_name = "Mysql"
  domains =
  privileges_list {
		privilege_name = "GlobalPrivileges"
		privileges =
		database = ""
		table_name = ""
		column_name = ""

  }
  project_id =
  s_s_h_key_name = ""
  version_id = "v1.0"
  secret_type = 0
  secret_binary = ""
  secret_string = "test"
  additional_config = "{}"
  tags = {
    "createdBy" = "terraform"
  }
}
```

Import

ssm secret can be imported using the id, e.g.

```
terraform import tencentcloud_ssm_secret.secret secret_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ssm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssm/v20190923"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudSsmSecret() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSsmSecretCreate,
		Read:   resourceTencentCloudSsmSecretRead,
		Update: resourceTencentCloudSsmSecretUpdate,
		Delete: resourceTencentCloudSsmSecretDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"secret_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Secret name.",
			},

			"user_name_prefix": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The prefix of the user account name is specified by the user, and the length is limited to 8 characters. The optional character set includes: numeric characters: [0, 9], lowercase characters: [a, z], uppercase characters: [A, Z] , Special characters (full English symbols): underscore (_), the prefix must start with an uppercase or lowercase letter.",
			},

			"enable_rotation": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "Whether to allow rotation.",
			},

			"rotation_begin_time": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "User-defined start rotation time, format: 2006-01-02 15:04:05. This parameter is required when EnableRotation is True.",
			},

			"rotation_frequency": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "The rotation period, in days, defaults to 1 day.",
			},

			"instance_i_d": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Cloud product instance ID.",
			},

			"description": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Descriptive information, used to describe the purpose in detail, supports up to 2048 bytes.",
			},

			"kms_key_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Specifies the KMS CMK that encrypts the credentials. If it is empty, it means to use Secrets Manager to encrypt the CMK you created by default. You can also specify a KMS CMK created in the same region for encryption.",
			},

			"product_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The name of the cloud product to which the credentials are bound.",
			},

			"domains": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The domain name of the account, IP format, support to fill in %.",
			},

			"privileges_list": {
				Required:    true,
				Type:        schema.TypeList,
				Description: "List of permissions that need to be granted when binding credentials to cloud product instances.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"privilege_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Privilege name, currently optional: GlobalPrivileges DatabasePrivileges TablePrivileges ColumnPrivilegesWhen the authority is DatabasePrivileges, the database name must be specified through the parameter Database;When the permission is TablePrivileges, the database name and the table name in the database must be specified through the parameters Database and TableName;When the permission is ColumnPrivileges, you must specify the database, the table name in the database, and the column name in the table through the parameters Database, TableName, and CoulmnName.",
						},
						"privileges": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "List of permissions. For Mysql products, the optional permission values ​​are:The optional values ​​​​of the permissions in GlobalPrivileges are: SELECT,INSERT,UPDATE,DELETE,CREATE,PROCESS,DROP,REFERENCES,INDEX,ALTER, SHOW DATABASES, CREATE TEMPORARY TABLES, LOCK TABLES, EXECUTE, CREATE VIEW, SHOW VIEW, CREATE ROUTINE, ALTER ROUTINE, EVENT, TRIGGER. Note that not passing this parameter means clearing the permission.The optional values ​​​​of the DatabasePrivileges permission are: SELECT, INSERT, UPDATE, DELETE, CREATE, DROP, REFERENCES, INDEX, ALTER, CREATE TEMPORARY TABLES, LOCK TABLES, EXECUTE, CREATE VIEW, SHOW VIEW, CREATE ROUTINE, ALTER ROUTINE, EVENT, TRIGGER. Note that not passing this parameter means clearing the permission.The optional values ​​of the TablePrivileges permissions are: The optional values ​​of the permissions are: SELECT,INSERT,UPDATE,DELETE,CREATE,DROP,REFERENCES,INDEX,ALTER , CREATE VIEW, SHOW VIEW, TRIGGER. Note that not passing this parameter means clearing the permission.The optional values ​​​​of the ColumnPrivileges permission are: SELECT,INSERT,UPDATE,REFERENCES. Note that not passing this parameter means clearing the permission.",
						},
						"database": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "This value is valid only when PrivilegeName is DatabasePrivileges.",
						},
						"table_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "This value is valid only when PrivilegeName is TablePrivileges, and at this time, you need to fill in Database to explicitly indicate the database instance where it is located.",
						},
						"column_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "This value only takes effect when PrivilegeName is ColumnPrivileges, and must be filled in at this time: Database - explicitly indicates the database instance where it is located. TableName - explicitly indicates the table.",
						},
					},
				},
			},

			"project_id": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "The ID of the project to which the key pair belongs after it is created.",
			},

			"s_s_h_key_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The name of the SSH key pair entered by the user can be composed of numbers, letters, and underscores. It can only start with numbers and letters, and the length cannot exceed 25 characters.",
			},

			"version_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Credential version. When querying credential information, it needs to be queried based on SecretName and VersionId. The maximum length is 64 bytes. It uses a combination of letters, numbers or - _ . and starts with a letter or number. If empty, the default initial credential version number is used. Optional, if it is empty or the credential is a cloud product credential, the version number defaults to SSM_Current.",
			},

			"secret_type": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Credential type, the default is custom credentials.",
			},

			"secret_binary": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The plaintext of the base64-encoded binary credential information. Only one of SecretBinary and SecretString must be set, and the maximum supported size is 4096 bytes.",
			},

			"secret_string": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Text type credential information in plain text (base64 encoding is not required). Only one SecretBinary and SecretString must be set, and the maximum supported size is 4096 bytes.",
			},

			"additional_config": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "A JSON formatted string to specify additional configuration for a particular credential type.",
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag description list.",
			},
		},
	}
}

func resourceTencentCloudSsmSecretCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ssm_secret.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		createSecretRequest  = ssm.NewCreateSecretRequest()
		createSecretResponse = ssm.NewCreateSecretResponse()
	)
	if v, ok := d.GetOk("secret_name"); ok {
		secretName = v.(string)
		request.SecretName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("user_name_prefix"); ok {
		request.UserNamePrefix = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("enable_rotation"); ok {
		request.EnableRotation = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("rotation_begin_time"); ok {
		request.RotationBeginTime = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("rotation_frequency"); ok {
		request.RotationFrequency = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("instance_i_d"); ok {
		request.InstanceID = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOk("kms_key_id"); ok {
		request.KmsKeyId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("product_name"); ok {
		request.ProductName = helper.String(v.(string))
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

	if v, ok := d.GetOkExists("project_id"); ok {
		request.ProjectId = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("s_s_h_key_name"); ok {
		request.SSHKeyName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("version_id"); ok {
		request.VersionId = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("secret_type"); ok {
		request.SecretType = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("secret_binary"); ok {
		request.SecretBinary = helper.String(v.(string))
	}

	if v, ok := d.GetOk("secret_string"); ok {
		request.SecretString = helper.String(v.(string))
	}

	if v, ok := d.GetOk("additional_config"); ok {
		request.AdditionalConfig = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSsmClient().CreateSecret(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create ssm secret failed, reason:%+v", logId, err)
		return err
	}

	secretName = *response.Response.secretName
	d.SetId(secretName)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := TagService{client: meta.(*TencentCloudClient).apiV3Conn}
		region := meta.(*TencentCloudClient).apiV3Conn.Region
		resourceName := fmt.Sprintf("qcs::ssm:%s:uin/:secret/%s", region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudSsmSecretRead(d, meta)
}

func resourceTencentCloudSsmSecretRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ssm_secret.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := SsmService{client: meta.(*TencentCloudClient).apiV3Conn}

	secretId := d.Id()

	secret, err := service.DescribeSsmSecretById(ctx, secretName)
	if err != nil {
		return err
	}

	if secret == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `SsmSecret` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if secret.SecretName != nil {
		_ = d.Set("secret_name", secret.SecretName)
	}

	if secret.UserNamePrefix != nil {
		_ = d.Set("user_name_prefix", secret.UserNamePrefix)
	}

	if secret.EnableRotation != nil {
		_ = d.Set("enable_rotation", secret.EnableRotation)
	}

	if secret.RotationBeginTime != nil {
		_ = d.Set("rotation_begin_time", secret.RotationBeginTime)
	}

	if secret.RotationFrequency != nil {
		_ = d.Set("rotation_frequency", secret.RotationFrequency)
	}

	if secret.InstanceID != nil {
		_ = d.Set("instance_i_d", secret.InstanceID)
	}

	if secret.Description != nil {
		_ = d.Set("description", secret.Description)
	}

	if secret.KmsKeyId != nil {
		_ = d.Set("kms_key_id", secret.KmsKeyId)
	}

	if secret.ProductName != nil {
		_ = d.Set("product_name", secret.ProductName)
	}

	if secret.Domains != nil {
		_ = d.Set("domains", secret.Domains)
	}

	if secret.PrivilegesList != nil {
		privilegesListList := []interface{}{}
		for _, privilegesList := range secret.PrivilegesList {
			privilegesListMap := map[string]interface{}{}

			if secret.PrivilegesList.PrivilegeName != nil {
				privilegesListMap["privilege_name"] = secret.PrivilegesList.PrivilegeName
			}

			if secret.PrivilegesList.Privileges != nil {
				privilegesListMap["privileges"] = secret.PrivilegesList.Privileges
			}

			if secret.PrivilegesList.Database != nil {
				privilegesListMap["database"] = secret.PrivilegesList.Database
			}

			if secret.PrivilegesList.TableName != nil {
				privilegesListMap["table_name"] = secret.PrivilegesList.TableName
			}

			if secret.PrivilegesList.ColumnName != nil {
				privilegesListMap["column_name"] = secret.PrivilegesList.ColumnName
			}

			privilegesListList = append(privilegesListList, privilegesListMap)
		}

		_ = d.Set("privileges_list", privilegesListList)

	}

	if secret.ProjectId != nil {
		_ = d.Set("project_id", secret.ProjectId)
	}

	if secret.SSHKeyName != nil {
		_ = d.Set("s_s_h_key_name", secret.SSHKeyName)
	}

	if secret.VersionId != nil {
		_ = d.Set("version_id", secret.VersionId)
	}

	if secret.SecretType != nil {
		_ = d.Set("secret_type", secret.SecretType)
	}

	if secret.SecretBinary != nil {
		_ = d.Set("secret_binary", secret.SecretBinary)
	}

	if secret.SecretString != nil {
		_ = d.Set("secret_string", secret.SecretString)
	}

	if secret.AdditionalConfig != nil {
		_ = d.Set("additional_config", secret.AdditionalConfig)
	}

	tcClient := meta.(*TencentCloudClient).apiV3Conn
	tagService := &TagService{client: tcClient}
	tags, err := tagService.DescribeResourceTags(ctx, "ssm", "secret", tcClient.Region, d.Id())
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudSsmSecretUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ssm_secret.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		enableSecretRequest  = ssm.NewEnableSecretRequest()
		enableSecretResponse = ssm.NewEnableSecretResponse()
	)

	secretId := d.Id()

	request.SecretName = &secretName

	immutableArgs := []string{"secret_name", "user_name_prefix", "enable_rotation", "rotation_begin_time", "rotation_frequency", "instance_i_d", "description", "kms_key_id", "product_name", "domains", "privileges_list", "project_id", "s_s_h_key_name", "version_id", "secret_type", "secret_binary", "secret_string", "additional_config"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("secret_name") {
		if v, ok := d.GetOk("secret_name"); ok {
			request.SecretName = helper.String(v.(string))
		}
	}

	if d.HasChange("enable_rotation") {
		if v, ok := d.GetOkExists("enable_rotation"); ok {
			request.EnableRotation = helper.Bool(v.(bool))
		}
	}

	if d.HasChange("rotation_begin_time") {
		if v, ok := d.GetOk("rotation_begin_time"); ok {
			request.RotationBeginTime = helper.String(v.(string))
		}
	}

	if d.HasChange("description") {
		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}
	}

	if d.HasChange("version_id") {
		if v, ok := d.GetOk("version_id"); ok {
			request.VersionId = helper.String(v.(string))
		}
	}

	if d.HasChange("secret_binary") {
		if v, ok := d.GetOk("secret_binary"); ok {
			request.SecretBinary = helper.String(v.(string))
		}
	}

	if d.HasChange("secret_string") {
		if v, ok := d.GetOk("secret_string"); ok {
			request.SecretString = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSsmClient().EnableSecret(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update ssm secret failed, reason:%+v", logId, err)
		return err
	}

	if d.HasChange("tags") {
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		tcClient := meta.(*TencentCloudClient).apiV3Conn
		tagService := &TagService{client: tcClient}
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		resourceName := BuildTagResourceName("ssm", "secret", tcClient.Region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
	}

	return resourceTencentCloudSsmSecretRead(d, meta)
}

func resourceTencentCloudSsmSecretDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ssm_secret.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := SsmService{client: meta.(*TencentCloudClient).apiV3Conn}
	secretId := d.Id()

	if err := service.DeleteSsmSecretById(ctx, secretName); err != nil {
		return err
	}

	return nil
}
