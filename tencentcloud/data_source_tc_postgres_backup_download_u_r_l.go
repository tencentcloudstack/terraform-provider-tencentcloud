/*
Use this data source to query detailed information of postgres backup_download_u_r_l

Example Usage

```hcl
data "tencentcloud_postgres_backup_download_u_r_l" "backup_download_u_r_l" {
  d_b_instance_id = ""
  backup_type = ""
  backup_id = ""
  u_r_l_expire_time =
  backup_download_restriction {
		restriction_type = ""
		vpc_restriction_effect = ""
		vpc_id_set =
		ip_restriction_effect = ""
		ip_set =

  }
    tags = {
    "createdBy" = "terraform"
  }
}
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudPostgresBackupDownloadURL() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudPostgresBackupDownloadURLRead,
		Schema: map[string]*schema.Schema{
			"d_b_instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},

			"backup_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Backup type. Valid values: `LogBackup`, `BaseBackup`.",
			},

			"backup_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Unique backup ID.",
			},

			"u_r_l_expire_time": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Validity period of a URL, which is 12 hours by default.",
			},

			"backup_download_restriction": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Backup download restriction.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"restriction_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Type of the network restrictions for downloading backup files. Valid values: `NONE` (backups can be downloaded over both private and public networks), `INTRANET` (backups can only be downloaded over the private network), `CUSTOMIZE` (backups can be downloaded over specified VPCs or at specified IPs).",
						},
						"vpc_restriction_effect": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Whether VPC is allowed. Valid values: `ALLOW` (allow), `DENY` (deny).",
						},
						"vpc_id_set": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "Whether it is allowed to download the VPC ID list of the backup files.",
						},
						"ip_restriction_effect": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Whether IP is allowed. Valid values: `ALLOW` (allow), `DENY` (deny).",
						},
						"ip_set": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "Whether it is allowed to download IP list of the backup files.",
						},
					},
				},
			},

			"backup_download_u_r_l": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Backup download URL.",
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag description list.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudPostgresBackupDownloadURLRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_postgres_backup_download_u_r_l.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("d_b_instance_id"); ok {
		paramMap["DBInstanceId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("backup_type"); ok {
		paramMap["BackupType"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("backup_id"); ok {
		paramMap["BackupId"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("u_r_l_expire_time"); v != nil {
		paramMap["URLExpireTime"] = helper.IntUint64(v.(int))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "backup_download_restriction"); ok {
		backupDownloadRestriction := postgres.BackupDownloadRestriction{}
		if v, ok := dMap["restriction_type"]; ok {
			backupDownloadRestriction.RestrictionType = helper.String(v.(string))
		}
		if v, ok := dMap["vpc_restriction_effect"]; ok {
			backupDownloadRestriction.VpcRestrictionEffect = helper.String(v.(string))
		}
		if v, ok := dMap["vpc_id_set"]; ok {
			vpcIdSetSet := v.(*schema.Set).List()
			backupDownloadRestriction.VpcIdSet = helper.InterfacesStringsPoint(vpcIdSetSet)
		}
		if v, ok := dMap["ip_restriction_effect"]; ok {
			backupDownloadRestriction.IpRestrictionEffect = helper.String(v.(string))
		}
		if v, ok := dMap["ip_set"]; ok {
			ipSetSet := v.(*schema.Set).List()
			backupDownloadRestriction.IpSet = helper.InterfacesStringsPoint(ipSetSet)
		}
		paramMap["backup_download_restriction"] = &backupDownloadRestriction
	}

	service := PostgresService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribePostgresBackupDownloadURLByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		backupDownloadURL = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(backupDownloadURL))
	if backupDownloadURL != nil {
		_ = d.Set("backup_download_u_r_l", backupDownloadURL)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string)); e != nil {
			return e
		}
	}
	return nil
}
