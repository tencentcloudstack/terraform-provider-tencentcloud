/*
Use this data source to query detailed information of cynosdb resource_package_sale_specs

Example Usage

```hcl
data "tencentcloud_cynosdb_resource_package_sale_specs" "resource_package_sale_specs" {
  instance_type  = "cynosdb-serverless"
  package_region = "china"
  package_type   = "CCU"
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudCynosdbResourcePackageSaleSpecs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCynosdbResourcePackageSaleSpecsRead,
		Schema: map[string]*schema.Schema{
			"instance_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance Type. Value range: cynosdb-serverless, cynosdb, cdb.",
			},
			"package_region": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Resource package usage region China - common in mainland China, overseas - common in Hong Kong, Macao, Taiwan, and overseas.",
			},
			"package_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Resource package type CCU - Computing resource package DISK - Storage resource package.",
			},
			"detail": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Resource package details note: This field may return null, indicating that a valid value cannot be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"package_region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"package_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Resource package type CCU - Compute resource package DISK - Store resource package Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"package_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Resource package version base basic version, common general version, enterprise enterprise version Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"min_package_spec": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "The minimum number of resources in the current version of the resource package, calculated in units of resources; Storage resource: GB Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"max_package_spec": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "The maximum number of resources in the current version of the resource package, calculated in units of resources; Storage resource: GB Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"expire_day": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Resource package validity period, in days. Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
					},
				},
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudCynosdbResourcePackageSaleSpecsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_cynosdb_resource_package_sale_specs.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId         = getLogId(contextNil)
		ctx           = context.WithValue(context.TODO(), logIdKey, logId)
		service       = CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}
		detail        []*cynosdb.SalePackageSpec
		instanceType  string
		packageRegion string
		packageType   string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_type"); ok {
		paramMap["InstanceType"] = helper.String(v.(string))
		instanceType = v.(string)
	}

	if v, ok := d.GetOk("package_region"); ok {
		paramMap["PackageRegion"] = helper.String(v.(string))
		packageRegion = v.(string)
	}

	if v, ok := d.GetOk("package_type"); ok {
		paramMap["PackageType"] = helper.String(v.(string))
		packageType = v.(string)
	}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCynosdbResourcePackageSaleSpecsByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}

		detail = result
		return nil
	})

	if err != nil {
		return err
	}

	ids := make([]string, 0, len(detail))
	ids = append(ids, instanceType, packageRegion, packageType)
	tmpList := make([]map[string]interface{}, 0, len(detail))

	if detail != nil {
		for _, salePackageSpec := range detail {
			salePackageSpecMap := map[string]interface{}{}

			if salePackageSpec.PackageRegion != nil {
				salePackageSpecMap["package_region"] = salePackageSpec.PackageRegion
			}

			if salePackageSpec.PackageType != nil {
				salePackageSpecMap["package_type"] = salePackageSpec.PackageType
			}

			if salePackageSpec.PackageVersion != nil {
				salePackageSpecMap["package_version"] = salePackageSpec.PackageVersion
			}

			if salePackageSpec.MinPackageSpec != nil {
				salePackageSpecMap["min_package_spec"] = salePackageSpec.MinPackageSpec
			}

			if salePackageSpec.MaxPackageSpec != nil {
				salePackageSpecMap["max_package_spec"] = salePackageSpec.MaxPackageSpec
			}

			if salePackageSpec.ExpireDay != nil {
				salePackageSpecMap["expire_day"] = salePackageSpec.ExpireDay
			}

			tmpList = append(tmpList, salePackageSpecMap)
		}

		_ = d.Set("detail", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}

	return nil
}
