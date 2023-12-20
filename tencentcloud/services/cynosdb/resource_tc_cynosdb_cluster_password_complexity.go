package cynosdb

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCynosdbClusterPasswordComplexity() *schema.Resource {
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
				ForceNew:    true,
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
	defer tccommon.LogElapsed("resource.tencentcloud_cynosdb_cluster_password_complexity.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var (
		request   = cynosdb.NewOpenClusterPasswordComplexityRequest()
		response  = cynosdb.NewOpenClusterPasswordComplexityResponse()
		clusterId string
	)
	if v, ok := d.GetOk("cluster_id"); ok {
		clusterId = v.(string)
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

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCynosdbClient().OpenClusterPasswordComplexity(request)
		if e != nil {
			return tccommon.RetryError(e)
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

	d.SetId(clusterId)

	flowId := *response.Response.FlowId
	service := CynosdbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	err = resource.Retry(6*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ok, err := service.DescribeFlow(ctx, flowId)
		if err != nil {
			if _, ok := err.(*sdkErrors.TencentCloudSDKError); !ok {
				return resource.RetryableError(err)
			} else {
				return resource.NonRetryableError(err)
			}
		}
		if ok {
			return nil
		} else {
			return resource.RetryableError(fmt.Errorf("create cynosdb clusterPasswordComplexity is processing"))
		}
	})

	if err != nil {
		log.Printf("[CRITAL]%s create cynosdb clusterPasswordComplexity fail, reason:%s\n", logId, err.Error())
		return err
	}

	err = service.CopyClusterPasswordComplexity(ctx, clusterId)
	if err != nil {
		log.Printf("[CRITAL]%s create cynosdb copyClusterPasswordComplexity fail, reason:%s\n", logId, err.Error())
		return err
	}

	return resourceTencentCloudCynosdbClusterPasswordComplexityRead(d, meta)
}

func resourceTencentCloudCynosdbClusterPasswordComplexityRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cynosdb_cluster_password_complexity.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := CynosdbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	clusterId := d.Id()

	clusterPasswordComplexity, err := service.DescribeCynosdbClusterPasswordComplexityById(ctx, clusterId)
	if err != nil {
		return err
	}

	if clusterPasswordComplexity == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CynosdbClusterPasswordComplexity` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("cluster_id", clusterId)

	if clusterPasswordComplexity.ValidatePasswordLength != nil {
		currentValue, err := strconv.ParseInt(*clusterPasswordComplexity.ValidatePasswordLength.CurrentValue, 10, 64)
		if err != nil {
			return err
		}
		_ = d.Set("validate_password_length", currentValue)
	}

	if clusterPasswordComplexity.ValidatePasswordMixedCaseCount != nil {
		currentValue, err := strconv.ParseInt(*clusterPasswordComplexity.ValidatePasswordMixedCaseCount.CurrentValue, 10, 64)
		if err != nil {
			return err
		}
		_ = d.Set("validate_password_mixed_case_count", currentValue)
	}

	if clusterPasswordComplexity.ValidatePasswordSpecialCharCount != nil {
		currentValue, err := strconv.ParseInt(*clusterPasswordComplexity.ValidatePasswordSpecialCharCount.CurrentValue, 10, 64)
		if err != nil {
			return err
		}
		_ = d.Set("validate_password_special_char_count", currentValue)
	}

	if clusterPasswordComplexity.ValidatePasswordNumberCount != nil {
		currentValue, err := strconv.ParseInt(*clusterPasswordComplexity.ValidatePasswordNumberCount.CurrentValue, 10, 64)
		if err != nil {
			return err
		}
		_ = d.Set("validate_password_number_count", currentValue)
	}

	if clusterPasswordComplexity.ValidatePasswordPolicy != nil {
		_ = d.Set("validate_password_policy", clusterPasswordComplexity.ValidatePasswordPolicy.CurrentValue)
	}

	if clusterPasswordComplexity.ValidatePasswordDictionary != nil {
		if clusterPasswordComplexity.ValidatePasswordDictionary.CurrentValue != nil {
			dictionary := strings.Split(*clusterPasswordComplexity.ValidatePasswordDictionary.CurrentValue, ",")
			_ = d.Set("validate_password_dictionary", dictionary)
		}
	}

	return nil
}

func resourceTencentCloudCynosdbClusterPasswordComplexityUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cynosdb_cluster_password_complexity.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var (
		request  = cynosdb.NewModifyClusterPasswordComplexityRequest()
		response = cynosdb.NewModifyClusterPasswordComplexityResponse()
	)

	clusterId := d.Id()
	request.ClusterId = &clusterId

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

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCynosdbClient().ModifyClusterPasswordComplexity(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update cynosdb clusterPasswordComplexity failed, reason:%+v", logId, err)
		return err
	}

	flowId := *response.Response.FlowId
	service := CynosdbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	err = resource.Retry(6*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ok, err := service.DescribeFlow(ctx, flowId)
		if err != nil {
			if _, ok := err.(*sdkErrors.TencentCloudSDKError); !ok {
				return resource.RetryableError(err)
			} else {
				return resource.NonRetryableError(err)
			}
		}
		if ok {
			return nil
		} else {
			return resource.RetryableError(fmt.Errorf("update cynosdb clusterPasswordComplexity is processing"))
		}
	})

	if err != nil {
		log.Printf("[CRITAL]%s update cynosdb clusterPasswordComplexity fail, reason:%s\n", logId, err.Error())
		return err
	}

	err = service.CopyClusterPasswordComplexity(ctx, clusterId)
	if err != nil {
		log.Printf("[CRITAL]%s update cynosdb copyClusterPasswordComplexity fail, reason:%s\n", logId, err.Error())
		return err
	}

	return resourceTencentCloudCynosdbClusterPasswordComplexityRead(d, meta)
}

func resourceTencentCloudCynosdbClusterPasswordComplexityDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cynosdb_cluster_password_complexity.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := CynosdbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	clusterId := d.Id()

	flowId, err := service.DeleteCynosdbClusterPasswordComplexityById(ctx, clusterId)
	if err != nil {
		return err
	}

	err = resource.Retry(6*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ok, err := service.DescribeFlow(ctx, flowId)
		if err != nil {
			if _, ok := err.(*sdkErrors.TencentCloudSDKError); !ok {
				return resource.RetryableError(err)
			} else {
				return resource.NonRetryableError(err)
			}
		}
		if ok {
			return nil
		} else {
			return resource.RetryableError(fmt.Errorf("delete cynosdb clusterPasswordComplexity is processing"))
		}
	})

	if err != nil {
		log.Printf("[CRITAL]%s delete cynosdb clusterPasswordComplexity fail, reason:%s\n", logId, err.Error())
		return err
	}

	return nil
}
