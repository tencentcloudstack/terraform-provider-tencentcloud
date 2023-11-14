/*
Provides a resource to create a cwp license_operation

Example Usage

```hcl
resource "tencentcloud_cwp_license_operation" "license_operation" {
  resource_id = ""
  license_type =
  is_all =
  quuid_list =
  task_id =
  filters {
		name = ""
		values =

  }
}
```

Import

cwp license_operation can be imported using the id, e.g.

```
terraform import tencentcloud_cwp_license_operation.license_operation license_operation_id
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cwp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cwp/v20180228"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudCwpLicenseOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCwpLicenseOperationCreate,
		Read:   resourceTencentCloudCwpLicenseOperationRead,
		Delete: resourceTencentCloudCwpLicenseOperationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"resource_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Resource ID.",
			},

			"license_type": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "LicenseType, 0 CWP Pro - Pay as you go ,1 CWP Pro - Monthly subscription , 2 CWP Ultimate - Monthly subscriptionDefault is 0.",
			},

			"is_all": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "All machines or not (If the total number of machines is greater than the number of licenses available on the current order, excess machines will be skipped).",
			},

			"quuid_list": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of machine quuid to be bound. This parameter is mandatory when IsAll = false. Otherwise, ignore this parameter. Maximum length =2000.",
			},

			"task_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Task ID.",
			},

			"filters": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				Description: "Status Binding progress status 0 running 1 complete 2 fail.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Filter key name。.",
						},
						"values": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "One or more filter values。.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudCwpLicenseOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cwp_license_operation.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		modifyLicenseBindsRequest  = cwp.NewModifyLicenseBindsRequest()
		modifyLicenseBindsResponse = cwp.NewModifyLicenseBindsResponse()
	)
	if v, ok := d.GetOk("resource_id"); ok {
		resourceId = v.(string)
		request.ResourceId = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("license_type"); ok {
		request.LicenseType = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("is_all"); ok {
		request.IsAll = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("quuid_list"); ok {
		quuidListSet := v.(*schema.Set).List()
		for i := range quuidListSet {
			quuidList := quuidListSet[i].(string)
			request.QuuidList = append(request.QuuidList, &quuidList)
		}
	}

	if v, ok := d.GetOkExists("task_id"); ok {
		request.TaskId = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("filters"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			filter := cwp.Filter{}
			if v, ok := dMap["name"]; ok {
				filter.Name = helper.String(v.(string))
			}
			if v, ok := dMap["values"]; ok {
				valuesSet := v.(*schema.Set).List()
				for i := range valuesSet {
					values := valuesSet[i].(string)
					filter.Values = append(filter.Values, &values)
				}
			}
			request.Filters = append(request.Filters, &filter)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCwpClient().ModifyLicenseBinds(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cwp licenseOperation failed, reason:%+v", logId, err)
		return err
	}

	resourceId = *response.Response.ResourceId
	d.SetId(resourceId)

	return resourceTencentCloudCwpLicenseOperationRead(d, meta)
}

func resourceTencentCloudCwpLicenseOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cwp_license_operation.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CwpService{client: meta.(*TencentCloudClient).apiV3Conn}

	licenseOperationId := d.Id()

	licenseOperation, err := service.DescribeCwpLicenseOperationById(ctx, resourceId)
	if err != nil {
		return err
	}

	if licenseOperation == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CwpLicenseOperation` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if licenseOperation.ResourceId != nil {
		_ = d.Set("resource_id", licenseOperation.ResourceId)
	}

	if licenseOperation.LicenseType != nil {
		_ = d.Set("license_type", licenseOperation.LicenseType)
	}

	if licenseOperation.IsAll != nil {
		_ = d.Set("is_all", licenseOperation.IsAll)
	}

	if licenseOperation.QuuidList != nil {
		_ = d.Set("quuid_list", licenseOperation.QuuidList)
	}

	if licenseOperation.TaskId != nil {
		_ = d.Set("task_id", licenseOperation.TaskId)
	}

	if licenseOperation.Filters != nil {
		filtersList := []interface{}{}
		for _, filters := range licenseOperation.Filters {
			filtersMap := map[string]interface{}{}

			if licenseOperation.Filters.Name != nil {
				filtersMap["name"] = licenseOperation.Filters.Name
			}

			if licenseOperation.Filters.Values != nil {
				filtersMap["values"] = licenseOperation.Filters.Values
			}

			filtersList = append(filtersList, filtersMap)
		}

		_ = d.Set("filters", filtersList)

	}

	return nil
}

func resourceTencentCloudCwpLicenseOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cwp_license_operation.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CwpService{client: meta.(*TencentCloudClient).apiV3Conn}
	licenseOperationId := d.Id()

	if err := service.DeleteCwpLicenseOperationById(ctx, resourceId); err != nil {
		return err
	}

	return nil
}
