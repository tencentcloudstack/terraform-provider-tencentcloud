/*
Provides a resource to create a cdb backup_download_restriction

Example Usage

```hcl
resource "tencentcloud_cdb_backup_download_restriction" "backup_download_restriction" {
  limit_type = "NoLimit"
  vpc_comparison_symbol = "In"
  ip_comparison_symbol = "In"
  limit_vpc {
		region = "ap-guangzhou"
		vpc_list =

  }
  limit_ip =
}
```

Import

cdb backup_download_restriction can be imported using the id, e.g.

```
terraform import tencentcloud_cdb_backup_download_restriction.backup_download_restriction backup_download_restriction_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	"log"
)

func resourceTencentCloudCdbBackupDownloadRestriction() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCdbBackupDownloadRestrictionCreate,
		Read:   resourceTencentCloudCdbBackupDownloadRestrictionRead,
		Update: resourceTencentCloudCdbBackupDownloadRestrictionUpdate,
		Delete: resourceTencentCloudCdbBackupDownloadRestrictionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"limit_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "NoLimit No limit, both internal and external networks can be downloaded; LimitOnlyIntranet Only intranet can be downloaded; Customize user-defined vpc:ip can be downloaded. LimitVpc and LimitIp can be set only when the value is Customize.",
			},

			"vpc_comparison_symbol": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "This parameter only supports In, which means that the vpc specified by LimitVpc can be downloaded. The default is In.",
			},

			"ip_comparison_symbol": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "In: The specified ip can be downloaded; NotIn: The specified ip cannot be downloaded. The default is In.",
			},

			"limit_vpc": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Vpc settings to limit downloads.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Restrict downloads from regions. Currently only the current region is supported.",
						},
						"vpc_list": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "List of vpcs to limit downloads.",
						},
					},
				},
			},

			"limit_ip": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Ip settings to limit downloads.",
			},
		},
	}
}

func resourceTencentCloudCdbBackupDownloadRestrictionCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdb_backup_download_restriction.create")()
	defer inconsistentCheck(d, meta)()

	var idsHash string
	if v, ok := d.GetOk("ids_hash"); ok {
		idsHash = v.(string)
	}

	d.SetId(idsHash)

	return resourceTencentCloudCdbBackupDownloadRestrictionUpdate(d, meta)
}

func resourceTencentCloudCdbBackupDownloadRestrictionRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdb_backup_download_restriction.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	backupDownloadRestrictionId := d.Id()

	backupDownloadRestriction, err := service.DescribeCdbBackupDownloadRestrictionById(ctx, idsHash)
	if err != nil {
		return err
	}

	if backupDownloadRestriction == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CdbBackupDownloadRestriction` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if backupDownloadRestriction.LimitType != nil {
		_ = d.Set("limit_type", backupDownloadRestriction.LimitType)
	}

	if backupDownloadRestriction.VpcComparisonSymbol != nil {
		_ = d.Set("vpc_comparison_symbol", backupDownloadRestriction.VpcComparisonSymbol)
	}

	if backupDownloadRestriction.IpComparisonSymbol != nil {
		_ = d.Set("ip_comparison_symbol", backupDownloadRestriction.IpComparisonSymbol)
	}

	if backupDownloadRestriction.LimitVpc != nil {
		limitVpcList := []interface{}{}
		for _, limitVpc := range backupDownloadRestriction.LimitVpc {
			limitVpcMap := map[string]interface{}{}

			if backupDownloadRestriction.LimitVpc.Region != nil {
				limitVpcMap["region"] = backupDownloadRestriction.LimitVpc.Region
			}

			if backupDownloadRestriction.LimitVpc.VpcList != nil {
				limitVpcMap["vpc_list"] = backupDownloadRestriction.LimitVpc.VpcList
			}

			limitVpcList = append(limitVpcList, limitVpcMap)
		}

		_ = d.Set("limit_vpc", limitVpcList)

	}

	if backupDownloadRestriction.LimitIp != nil {
		_ = d.Set("limit_ip", backupDownloadRestriction.LimitIp)
	}

	return nil
}

func resourceTencentCloudCdbBackupDownloadRestrictionUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdb_backup_download_restriction.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := cdb.NewModifyBackupDownloadRestrictionRequest()

	backupDownloadRestrictionId := d.Id()

	request.IdsHash = &idsHash

	immutableArgs := []string{"limit_type", "vpc_comparison_symbol", "ip_comparison_symbol", "limit_vpc", "limit_ip"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCdbClient().ModifyBackupDownloadRestriction(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update cdb backupDownloadRestriction failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCdbBackupDownloadRestrictionRead(d, meta)
}

func resourceTencentCloudCdbBackupDownloadRestrictionDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdb_backup_download_restriction.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
