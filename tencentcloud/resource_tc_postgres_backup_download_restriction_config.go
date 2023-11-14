/*
Provides a resource to create a postgres backup_download_restriction_config

Example Usage

```hcl
resource "tencentcloud_postgres_backup_download_restriction_config" "backup_download_restriction_config" {
  restriction_type = ""
  vpc_restriction_effect = ""
  vpc_id_set =
  ip_restriction_effect = ""
  ip_set =
}
```

Import

postgres backup_download_restriction_config can be imported using the id, e.g.

```
terraform import tencentcloud_postgres_backup_download_restriction_config.backup_download_restriction_config backup_download_restriction_config_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	postgres "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudPostgresBackupDownloadRestrictionConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPostgresBackupDownloadRestrictionConfigCreate,
		Read:   resourceTencentCloudPostgresBackupDownloadRestrictionConfigRead,
		Update: resourceTencentCloudPostgresBackupDownloadRestrictionConfigUpdate,
		Delete: resourceTencentCloudPostgresBackupDownloadRestrictionConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"restriction_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Backup file download restriction type:NONE:Unlimited, both internal and external networks can be downloaded.INTRANET:Only intranet downloads are allowed.CUSTOMIZE:Customize the vpc or ip that limits downloads.",
			},

			"vpc_restriction_effect": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Vpc limit Strategy:ALLOW,DENY.",
			},

			"vpc_id_set": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The list of vpcIds that allow or deny downloading of backup files.",
			},

			"ip_restriction_effect": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Ip limit Strategy:ALLOW,DENY.",
			},

			"ip_set": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The list of ip&amp;amp;#39;s that are allowed or denied to download backup files.",
			},
		},
	}
}

func resourceTencentCloudPostgresBackupDownloadRestrictionConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgres_backup_download_restriction_config.create")()
	defer inconsistentCheck(d, meta)()

	var restrictionType string
	if v, ok := d.GetOk("restriction_type"); ok {
		restrictionType = v.(string)
	}

	d.SetId(restrictionType)

	return resourceTencentCloudPostgresBackupDownloadRestrictionConfigUpdate(d, meta)
}

func resourceTencentCloudPostgresBackupDownloadRestrictionConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgres_backup_download_restriction_config.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := PostgresService{client: meta.(*TencentCloudClient).apiV3Conn}

	backupDownloadRestrictionConfigId := d.Id()

	BackupDownloadRestrictionConfig, err := service.DescribePostgresBackupDownloadRestrictionConfigById(ctx, restrictionType)
	if err != nil {
		return err
	}

	if BackupDownloadRestrictionConfig == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `PostgresBackupDownloadRestrictionConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if BackupDownloadRestrictionConfig.RestrictionType != nil {
		_ = d.Set("restriction_type", BackupDownloadRestrictionConfig.RestrictionType)
	}

	if BackupDownloadRestrictionConfig.VpcRestrictionEffect != nil {
		_ = d.Set("vpc_restriction_effect", BackupDownloadRestrictionConfig.VpcRestrictionEffect)
	}

	if BackupDownloadRestrictionConfig.VpcIdSet != nil {
		_ = d.Set("vpc_id_set", BackupDownloadRestrictionConfig.VpcIdSet)
	}

	if BackupDownloadRestrictionConfig.IpRestrictionEffect != nil {
		_ = d.Set("ip_restriction_effect", BackupDownloadRestrictionConfig.IpRestrictionEffect)
	}

	if BackupDownloadRestrictionConfig.IpSet != nil {
		_ = d.Set("ip_set", BackupDownloadRestrictionConfig.IpSet)
	}

	return nil
}

func resourceTencentCloudPostgresBackupDownloadRestrictionConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgres_backup_download_restriction_config.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := postgres.NewModifyBackupDownloadRestrictionRequest()

	backupDownloadRestrictionConfigId := d.Id()

	request.RestrictionType = &restrictionType

	immutableArgs := []string{"restriction_type", "vpc_restriction_effect", "vpc_id_set", "ip_restriction_effect", "ip_set"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("vpc_restriction_effect") {
		if v, ok := d.GetOk("vpc_restriction_effect"); ok {
			request.VpcRestrictionEffect = helper.String(v.(string))
		}
	}

	if d.HasChange("vpc_id_set") {
		if v, ok := d.GetOk("vpc_id_set"); ok {
			vpcIdSetSet := v.(*schema.Set).List()
			for i := range vpcIdSetSet {
				vpcIdSet := vpcIdSetSet[i].(string)
				request.VpcIdSet = append(request.VpcIdSet, &vpcIdSet)
			}
		}
	}

	if d.HasChange("ip_restriction_effect") {
		if v, ok := d.GetOk("ip_restriction_effect"); ok {
			request.IpRestrictionEffect = helper.String(v.(string))
		}
	}

	if d.HasChange("ip_set") {
		if v, ok := d.GetOk("ip_set"); ok {
			ipSetSet := v.(*schema.Set).List()
			for i := range ipSetSet {
				ipSet := ipSetSet[i].(string)
				request.IpSet = append(request.IpSet, &ipSet)
			}
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UsePostgresClient().ModifyBackupDownloadRestriction(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update postgres BackupDownloadRestrictionConfig failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudPostgresBackupDownloadRestrictionConfigRead(d, meta)
}

func resourceTencentCloudPostgresBackupDownloadRestrictionConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgres_backup_download_restriction_config.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
