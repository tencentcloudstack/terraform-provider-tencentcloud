/*
Provides a resource to create a cynosdb cluster_password_complexity

Example Usage

```hcl
resource "tencentcloud_cynosdb_cluster_password_complexity" "cluster_password_complexity" {
  cluster_id = "cynosdbpg-bzxxrmtq"
  validate_password_length = 8
  validate_password_mixed_case_count = 1
  validate_password_special_char_count = 1
  validate_password_number_count = 1
  validate_password_policy = "MEDIUM"
  validate_password_dictionary =
}
```

Import

cynosdb cluster_password_complexity can be imported using the id, e.g.

```
terraform import tencentcloud_cynosdb_cluster_password_complexity.cluster_password_complexity cluster_password_complexity_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"time"
)

func resourceTencentCloudCynosdbClusterPasswordComplexity() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCynosdbClusterPasswordComplexityCreate,
		Read:   resourceTencentCloudCynosdbClusterPasswordComplexityRead,
		Update: resourceTencentCloudCynosdbClusterPasswordComplexityUpdate,
		Delete: resourceTencentCloudCynosdbClusterPasswordComplexityDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Cluster ID.",
			},

			"validate_password_length": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Password length.",
			},

			"validate_password_mixed_case_count": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Number of uppercase and lowercase characters.",
			},

			"validate_password_special_char_count": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Number of special characters.",
			},

			"validate_password_number_count": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Number of digits.",
			},

			"validate_password_policy": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Password strength (MEDIUM, STRONG).",
			},

			"validate_password_dictionary": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Data dictionary.",
			},
		},
	}
}

func resourceTencentCloudCynosdbClusterPasswordComplexityCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_cluster_password_complexity.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request   = cynosdb.NewOpenClusterPasswordComplexityRequest()
		response  = cynosdb.NewOpenClusterPasswordComplexityResponse()
		paramName string
	)
	if v, ok := d.GetOk("cluster_id"); ok {
		request.ClusterId = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("validate_password_length"); ok {
		request.ValidatePasswordLength = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("validate_password_mixed_case_count"); ok {
		request.ValidatePasswordMixedCaseCount = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("validate_password_special_char_count"); ok {
		request.ValidatePasswordSpecialCharCount = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("validate_password_number_count"); ok {
		request.ValidatePasswordNumberCount = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("validate_password_policy"); ok {
		request.ValidatePasswordPolicy = helper.String(v.(string))
	}

	if v, ok := d.GetOk("validate_password_dictionary"); ok {
		validatePasswordDictionarySet := v.(*schema.Set).List()
		for i := range validatePasswordDictionarySet {
			validatePasswordDictionary := validatePasswordDictionarySet[i].(string)
			request.ValidatePasswordDictionary = append(request.ValidatePasswordDictionary, &validatePasswordDictionary)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCynosdbClient().OpenClusterPasswordComplexity(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cynosdb clusterPasswordComplexity failed, reason:%+v", logId, err)
		return err
	}

	paramName = *response.Response.ParamName
	d.SetId(paramName)

	service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"success"}, 30*readRetryTimeout, time.Second, service.CynosdbClusterPasswordComplexityStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudCynosdbClusterPasswordComplexityRead(d, meta)
}

func resourceTencentCloudCynosdbClusterPasswordComplexityRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_cluster_password_complexity.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	clusterPasswordComplexityId := d.Id()

	clusterPasswordComplexity, err := service.DescribeCynosdbClusterPasswordComplexityById(ctx, paramName)
	if err != nil {
		return err
	}

	if clusterPasswordComplexity == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CynosdbClusterPasswordComplexity` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if clusterPasswordComplexity.ClusterId != nil {
		_ = d.Set("cluster_id", clusterPasswordComplexity.ClusterId)
	}

	if clusterPasswordComplexity.ValidatePasswordLength != nil {
		_ = d.Set("validate_password_length", clusterPasswordComplexity.ValidatePasswordLength)
	}

	if clusterPasswordComplexity.ValidatePasswordMixedCaseCount != nil {
		_ = d.Set("validate_password_mixed_case_count", clusterPasswordComplexity.ValidatePasswordMixedCaseCount)
	}

	if clusterPasswordComplexity.ValidatePasswordSpecialCharCount != nil {
		_ = d.Set("validate_password_special_char_count", clusterPasswordComplexity.ValidatePasswordSpecialCharCount)
	}

	if clusterPasswordComplexity.ValidatePasswordNumberCount != nil {
		_ = d.Set("validate_password_number_count", clusterPasswordComplexity.ValidatePasswordNumberCount)
	}

	if clusterPasswordComplexity.ValidatePasswordPolicy != nil {
		_ = d.Set("validate_password_policy", clusterPasswordComplexity.ValidatePasswordPolicy)
	}

	if clusterPasswordComplexity.ValidatePasswordDictionary != nil {
		_ = d.Set("validate_password_dictionary", clusterPasswordComplexity.ValidatePasswordDictionary)
	}

	return nil
}

func resourceTencentCloudCynosdbClusterPasswordComplexityUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_cluster_password_complexity.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		modifyClusterPasswordComplexityRequest  = cynosdb.NewModifyClusterPasswordComplexityRequest()
		modifyClusterPasswordComplexityResponse = cynosdb.NewModifyClusterPasswordComplexityResponse()
	)

	clusterPasswordComplexityId := d.Id()

	request.ParamName = &paramName

	immutableArgs := []string{"cluster_id", "validate_password_length", "validate_password_mixed_case_count", "validate_password_special_char_count", "validate_password_number_count", "validate_password_policy", "validate_password_dictionary"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("cluster_id") {
		if v, ok := d.GetOk("cluster_id"); ok {
			request.ClusterId = helper.String(v.(string))
		}
	}

	if d.HasChange("validate_password_length") {
		if v, ok := d.GetOkExists("validate_password_length"); ok {
			request.ValidatePasswordLength = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("validate_password_mixed_case_count") {
		if v, ok := d.GetOkExists("validate_password_mixed_case_count"); ok {
			request.ValidatePasswordMixedCaseCount = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("validate_password_special_char_count") {
		if v, ok := d.GetOkExists("validate_password_special_char_count"); ok {
			request.ValidatePasswordSpecialCharCount = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("validate_password_number_count") {
		if v, ok := d.GetOkExists("validate_password_number_count"); ok {
			request.ValidatePasswordNumberCount = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("validate_password_policy") {
		if v, ok := d.GetOk("validate_password_policy"); ok {
			request.ValidatePasswordPolicy = helper.String(v.(string))
		}
	}

	if d.HasChange("validate_password_dictionary") {
		if v, ok := d.GetOk("validate_password_dictionary"); ok {
			validatePasswordDictionarySet := v.(*schema.Set).List()
			for i := range validatePasswordDictionarySet {
				validatePasswordDictionary := validatePasswordDictionarySet[i].(string)
				request.ValidatePasswordDictionary = append(request.ValidatePasswordDictionary, &validatePasswordDictionary)
			}
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCynosdbClient().ModifyClusterPasswordComplexity(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update cynosdb clusterPasswordComplexity failed, reason:%+v", logId, err)
		return err
	}

	service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"success"}, 30*readRetryTimeout, time.Second, service.CynosdbClusterPasswordComplexityStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"success"}, 30*readRetryTimeout, time.Second, service.CynosdbClusterPasswordComplexityStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudCynosdbClusterPasswordComplexityRead(d, meta)
}

func resourceTencentCloudCynosdbClusterPasswordComplexityDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_cluster_password_complexity.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}
	clusterPasswordComplexityId := d.Id()

	if err := service.DeleteCynosdbClusterPasswordComplexityById(ctx, paramName); err != nil {
		return err
	}

	service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"success"}, 30*readRetryTimeout, time.Second, service.CynosdbClusterPasswordComplexityStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return nil
}
