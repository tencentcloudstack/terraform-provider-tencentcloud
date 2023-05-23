/*
Provides a resource to create a postgresql backup_download_restriction_config

Example Usage

Unlimit the restriction of the backup file download.
```hcl
resource "tencentcloud_postgresql_backup_download_restriction_config" "backup_download_restriction_config" {
  restriction_type = "NONE"
}
```

Set the download only to allow the intranet downloads.
```hcl
resource "tencentcloud_postgresql_backup_download_restriction_config" "backup_download_restriction_config" {
  restriction_type = "INTRANET"
}
```

Restrict the backup file download by customizing.
```hcl
resource "tencentcloud_vpc" "pg_vpc" {
	name       = var.instance_name
	cidr_block = var.vpc_cidr
}

resource "tencentcloud_postgresql_backup_download_restriction_config" "backup_download_restriction_config" {
  restriction_type = "CUSTOMIZE"
  vpc_restriction_effect = "DENY"
  vpc_id_set = [tencentcloud_vpc.pg_vpc2.id]
  ip_restriction_effect = "DENY"
  ip_set = ["192.168.0.0"]
}
```

Import

postgresql backup_download_restriction_config can be imported using the id, e.g.

```
terraform import tencentcloud_postgresql_backup_download_restriction_config.backup_download_restriction_config backup_download_restriction_config_id
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	postgresql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudPostgresqlBackupDownloadRestrictionConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPostgresqlBackupDownloadRestrictionConfigCreate,
		Read:   resourceTencentCloudPostgresqlBackupDownloadRestrictionConfigRead,
		Update: resourceTencentCloudPostgresqlBackupDownloadRestrictionConfigUpdate,
		Delete: resourceTencentCloudPostgresqlBackupDownloadRestrictionConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"restriction_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Backup file download restriction type: NONE:Unlimited, both internal and external networks can be downloaded. INTRANET:Only intranet downloads are allowed. CUSTOMIZE:Customize the vpc or ip that limits downloads.",
			},

			"vpc_restriction_effect": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "vpc limit Strategy: ALLOW, DENY.",
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
				Description: "ip limit Strategy: ALLOW, DENY.",
			},

			"ip_set": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The list of ips that are allowed or denied to download backup files.",
			},
		},
	}
}

func resourceTencentCloudPostgresqlBackupDownloadRestrictionConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgresql_backup_download_restriction_config.create")()
	defer inconsistentCheck(d, meta)()

	var resType string
	if v, ok := d.GetOk("restriction_type"); ok {
		resType = v.(string)
	}

	d.SetId(resType)

	return resourceTencentCloudPostgresqlBackupDownloadRestrictionConfigUpdate(d, meta)
}

func resourceTencentCloudPostgresqlBackupDownloadRestrictionConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgresql_backup_download_restriction_config.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := PostgresqlService{client: meta.(*TencentCloudClient).apiV3Conn}

	resType := d.Id()

	BackupDownloadRestrictionConfig, err := service.DescribePostgresqlBackupDownloadRestrictionConfigById(ctx, resType)
	if err != nil {
		return err
	}

	if BackupDownloadRestrictionConfig == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `PostgresqlBackupDownloadRestrictionConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
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

func resourceTencentCloudPostgresqlBackupDownloadRestrictionConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgresql_backup_download_restriction_config.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := postgresql.NewModifyBackupDownloadRestrictionRequest()

	resType := d.Id()

	if d.HasChange("restriction_type") {
		if v, ok := d.GetOk("restriction_type"); ok {
			resType = v.(string)
		}
	}
	request.RestrictionType = &resType

	if d.HasChange("vpc_restriction_effect") {
		if v, ok := d.GetOk("vpc_restriction_effect"); ok {
			request.VpcRestrictionEffect = helper.String(v.(string))
		}
	}

	if d.HasChange("vpc_id_set") {
		if v, ok := d.GetOk("vpc_id_set"); ok {
			vpcIdSetSet := v.(*schema.Set).List()
			for i := range vpcIdSetSet {
				if vpcIdSetSet[i] != nil {
					vpcIdSet := vpcIdSetSet[i].(string)
					request.VpcIdSet = append(request.VpcIdSet, &vpcIdSet)
				}
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
				if ipSetSet[i] != nil {
					ipSet := ipSetSet[i].(string)
					request.IpSet = append(request.IpSet, &ipSet)
				}
			}
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UsePostgresqlClient().ModifyBackupDownloadRestriction(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update postgresql BackupDownloadRestrictionConfig failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudPostgresqlBackupDownloadRestrictionConfigRead(d, meta)
}

func resourceTencentCloudPostgresqlBackupDownloadRestrictionConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgresql_backup_download_restriction_config.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
