/*
Provides a resource to create a cynosdb resource_package

Example Usage

```hcl
resource "tencentcloud_cynosdb_resource_package" "resource_package" {
  instance_type = "cdb"
  package_region = "china"
  package_type = "CCU"
  package_version = "base"
  package_spec =
  expire_day = 180
  package_count = 1
  package_name = "PackageName"
}
```

Import

cynosdb resource_package can be imported using the id, e.g.

```
terraform import tencentcloud_cynosdb_resource_package.resource_package resource_package_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudCynosdbResourcePackage() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCynosdbResourcePackageCreate,
		Read:   resourceTencentCloudCynosdbResourcePackageRead,
		Update: resourceTencentCloudCynosdbResourcePackageUpdate,
		Delete: resourceTencentCloudCynosdbResourcePackageDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance Type.",
			},

			"package_region": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Resource package usage region China - common in mainland China, overseas - common in Hong Kong, Macao, Taiwan, and overseas.",
			},

			"package_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Resource package type: CCU computing resource package, DISK storage resource package.",
			},

			"package_version": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Resource package version base basic version, common general version, enterprise enterprise version.",
			},

			"package_spec": {
				Required:    true,
				Type:        schema.TypeFloat,
				Description: "Resource package size, calculated in 10000 units; Storage resources: GB.",
			},

			"expire_day": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Validity period of resource package, in days.",
			},

			"package_count": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Number of purchased resource packs.",
			},

			"package_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Resource Package Name.",
			},
		},
	}
}

func resourceTencentCloudCynosdbResourcePackageCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_resource_package.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request = cynosdb.NewCreateResourcePackageRequest()
		// response  = cynosdb.NewCreateResourcePackageResponse()
		// packageId string
	)
	if v, ok := d.GetOk("instance_type"); ok {
		request.InstanceType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("package_region"); ok {
		request.PackageRegion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("package_type"); ok {
		request.PackageType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("package_version"); ok {
		request.PackageVersion = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("package_spec"); ok {
		request.PackageSpec = helper.Float64(v.(float64))
	}

	if v, ok := d.GetOkExists("expire_day"); ok {
		request.ExpireDay = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("package_count"); ok {
		request.PackageCount = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("package_name"); ok {
		request.PackageName = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCynosdbClient().CreateResourcePackage(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		// response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cynosdb resourcePackage failed, reason:%+v", logId, err)
		return err
	}

	// packageId = *response.Response.PackageId
	// d.SetId(helper.String(packageId))

	return resourceTencentCloudCynosdbResourcePackageRead(d, meta)
}

func resourceTencentCloudCynosdbResourcePackageRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_resource_package.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	packageId := d.Id()
	resourcePackage, err := service.DescribeCynosdbResourcePackageById(ctx, packageId)
	if err != nil {
		return err
	}

	if resourcePackage == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CynosdbResourcePackage` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	// if resourcePackage.InstanceType != nil {
	// 	_ = d.Set("instance_type", resourcePackage.InstanceType)
	// }

	// if resourcePackage.PackageRegion != nil {
	// 	_ = d.Set("package_region", resourcePackage.PackageRegion)
	// }

	// if resourcePackage.PackageType != nil {
	// 	_ = d.Set("package_type", resourcePackage.PackageType)
	// }

	// if resourcePackage.PackageVersion != nil {
	// 	_ = d.Set("package_version", resourcePackage.PackageVersion)
	// }

	// if resourcePackage.PackageSpec != nil {
	// 	_ = d.Set("package_spec", resourcePackage.PackageSpec)
	// }

	// if resourcePackage.ExpireDay != nil {
	// 	_ = d.Set("expire_day", resourcePackage.ExpireDay)
	// }

	// if resourcePackage.PackageCount != nil {
	// 	_ = d.Set("package_count", resourcePackage.PackageCount)
	// }

	// if resourcePackage.PackageName != nil {
	// 	_ = d.Set("package_name", resourcePackage.PackageName)
	// }

	return nil
}

func resourceTencentCloudCynosdbResourcePackageUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_resource_package.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := cynosdb.NewModifyResourcePackageNameRequest()

	packageId := d.Id()

	request.PackageId = &packageId

	immutableArgs := []string{"instance_type", "package_region", "package_type", "package_version", "package_spec", "expire_day", "package_count", "package_name"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("package_name") {
		if v, ok := d.GetOk("package_name"); ok {
			request.PackageName = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCynosdbClient().ModifyResourcePackageName(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update cynosdb resourcePackage failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCynosdbResourcePackageRead(d, meta)
}

func resourceTencentCloudCynosdbResourcePackageDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_resource_package.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}
	packageId := d.Id()

	if err := service.DeleteCynosdbResourcePackageById(ctx, packageId); err != nil {
		return err
	}

	return nil
}
