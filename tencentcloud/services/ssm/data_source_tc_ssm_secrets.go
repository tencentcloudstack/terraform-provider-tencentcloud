package ssm

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ssm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssm/v20190923"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudSsmSecrets() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudSsmSecretsRead,
		Schema: map[string]*schema.Schema{
			"order_type": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "The order to sort the create time of secret. `0` - desc, `1` - asc. Default value is `0`.",
			},
			"state": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Filter by state of secret. `0` - all secrets are queried, `1` - only Enabled secrets are queried, `2` - only Disabled secrets are queried, `3` - only PendingDelete secrets are queried.",
			},
			"secret_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Secret name used to filter result.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tags to filter secret.",
			},
			"secret_type": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "0- represents user-defined credentials, defaults to 0. 1- represents the user's cloud product credentials. 2- represents SSH key pair credentials. 3- represents cloud API key pair credentials.",
			},
			"product_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "This parameter only takes effect when the SecretType parameter value is 1. When the SecretType value is 1, if the Product Name value is empty, it means to query all types of cloud product credentials. If the Product Name value is MySQL, it means to query MySQL database credentials. If the Product Name value is Tdsql mysql, it means to query Tdsql (MySQL version) credentials.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			"secret_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of SSM secrets.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"secret_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of secret.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description of secret.",
						},
						"kms_key_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "KMS keyId used to encrypt secret.",
						},
						"create_uin": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Uin of Creator.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Status of secret.",
						},
						"delete_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Delete time of CMK.",
						},
						"create_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Create time of secret.",
						},
						"kms_key_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "KMS CMK type used to encrypt credentials, DEFAULT represents the default key created by SecretsManager, and CUSTOMER represents the user specified key.",
						},
						"rotation_status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "1: - Turn on the rotation; 0- No rotation Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"next_rotation_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Next rotation start time, uinx timestamp.",
						},
						"secret_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "0- User defined credentials; 1- Cloud product credentials; 2- SSH key pair credentials; 3- Cloud API key pair credentials.",
						},
						"product_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cloud product name, only effective when SecretType is 1, which means the credential type is cloud product credential.",
						},
						"resource_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "When the credential type is SSH key pair credential, this field is valid and is used to represent the name of the SSH key pair credential.",
						},
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "When the credential type is SSH key pair credential, this field is valid and represents the item ID to which the SSH key pair belongs.",
						},
						"associated_instance_ids": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "When the credential type is SSH key pair credential, this field is valid and is used to represent the CVM instance ID associated with the SSH key pair.",
						},
						"target_uin": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "When the credential type is a cloud API key pair credential, this field is valid and is used to represent the user UIN to which the cloud API key pair belongs.",
						},
						"rotation_frequency": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The frequency of rotation, in days, takes effect when rotation is on.",
						},
						"resource_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The cloud product instance ID number corresponding to the cloud product credentials.",
						},
						"rotation_begin_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user specified rotation start time.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudSsmSecretsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_ssm_secrets.read")()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		ssmService = SsmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		secrets    []*ssm.SecretMetadata
	)

	param := make(map[string]interface{})
	if v, ok := d.GetOk("order_type"); ok {
		param["order_type"] = v.(int)
	}

	if v, ok := d.GetOk("state"); ok {
		param["state"] = v.(int)
	}

	if v, ok := d.GetOk("secret_name"); ok {
		param["secret_name"] = v.(string)
	}

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		param["tag_filter"] = tags
	}

	if v, ok := d.GetOk("secret_type"); ok {
		param["secret_type"] = v.(string)
	}

	if v, ok := d.GetOk("product_name"); ok {
		param["product_name"] = v.(string)
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		results, e := ssmService.DescribeSecretsByFilter(ctx, param)
		if e != nil {
			return tccommon.RetryError(e)
		}

		secrets = results
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s read SSM secrets failed, reason:%+v", logId, err)
		return err
	}

	secretList := make([]map[string]interface{}, 0, len(secrets))
	secretNames := make([]string, 0, len(secrets))
	for _, secret := range secrets {
		mapping := map[string]interface{}{
			"secret_name":             secret.SecretName,
			"description":             secret.Description,
			"kms_key_id":              secret.KmsKeyId,
			"create_uin":              secret.CreateUin,
			"status":                  secret.Status,
			"delete_time":             secret.DeleteTime,
			"create_time":             secret.CreateTime,
			"kms_key_type":            secret.KmsKeyType,
			"rotation_status":         secret.RotationStatus,
			"next_rotation_time":      secret.NextRotationTime,
			"secret_type":             secret.SecretType,
			"product_name":            secret.ProductName,
			"resource_name":           secret.ResourceName,
			"project_id":              secret.ProjectID,
			"associated_instance_ids": secret.AssociatedInstanceIDs,
			"target_uin":              secret.TargetUin,
			"rotation_frequency":      secret.RotationFrequency,
			"resource_id":             secret.ResourceID,
			"rotation_begin_time":     secret.RotationBeginTime,
		}

		secretList = append(secretList, mapping)
		secretNames = append(secretNames, *secret.SecretName)
	}

	d.SetId(helper.DataResourceIdsHash(secretNames))
	if e := d.Set("secret_list", secretList); e != nil {
		log.Printf("[CRITAL]%s provider set SSM secret list fail, reason:%+v", logId, e)
		return e
	}

	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {
		return tccommon.WriteToFile(output.(string), secretList)
	}

	return nil
}
