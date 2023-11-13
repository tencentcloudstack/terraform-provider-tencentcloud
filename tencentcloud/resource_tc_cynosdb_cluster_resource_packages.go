/*
Provides a resource to create a cynosdb cluster_resource_packages

Example Usage

```hcl
resource "tencentcloud_cynosdb_cluster_resource_packages" "cluster_resource_packages" {
  package_ids =
  cluster_id = "cynosdb-qwerty"
}
```

Import

cynosdb cluster_resource_packages can be imported using the id, e.g.

```
terraform import tencentcloud_cynosdb_cluster_resource_packages.cluster_resource_packages cluster_resource_packages_id
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudCynosdbClusterResourcePackages() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCynosdbClusterResourcePackagesCreate,
		Read:   resourceTencentCloudCynosdbClusterResourcePackagesRead,
		Delete: resourceTencentCloudCynosdbClusterResourcePackagesDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"package_ids": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Resource Package Unique ID.",
			},

			"cluster_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Cluster ID.",
			},
		},
	}
}

func resourceTencentCloudCynosdbClusterResourcePackagesCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_cluster_resource_packages.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request   = cynosdb.NewBindClusterResourcePackagesRequest()
		response  = cynosdb.NewBindClusterResourcePackagesResponse()
		packageId string
	)
	if v, ok := d.GetOk("package_ids"); ok {
		packageIdsSet := v.(*schema.Set).List()
		for i := range packageIdsSet {
			packageIds := packageIdsSet[i].(string)
			request.PackageIds = append(request.PackageIds, &packageIds)
		}
	}

	if v, ok := d.GetOk("cluster_id"); ok {
		request.ClusterId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCynosdbClient().BindClusterResourcePackages(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cynosdb clusterResourcePackages failed, reason:%+v", logId, err)
		return err
	}

	packageId = *response.Response.PackageId
	d.SetId(packageId)

	return resourceTencentCloudCynosdbClusterResourcePackagesRead(d, meta)
}

func resourceTencentCloudCynosdbClusterResourcePackagesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_cluster_resource_packages.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	clusterResourcePackagesId := d.Id()

	clusterResourcePackages, err := service.DescribeCynosdbClusterResourcePackagesById(ctx, packageId)
	if err != nil {
		return err
	}

	if clusterResourcePackages == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CynosdbClusterResourcePackages` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if clusterResourcePackages.PackageIds != nil {
		_ = d.Set("package_ids", clusterResourcePackages.PackageIds)
	}

	if clusterResourcePackages.ClusterId != nil {
		_ = d.Set("cluster_id", clusterResourcePackages.ClusterId)
	}

	return nil
}

func resourceTencentCloudCynosdbClusterResourcePackagesDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_cluster_resource_packages.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}
	clusterResourcePackagesId := d.Id()

	if err := service.DeleteCynosdbClusterResourcePackagesById(ctx, packageId); err != nil {
		return err
	}

	return nil
}
