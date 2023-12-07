package tencentcloud

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cos "github.com/tencentyun/cos-go-sdk-v5"
)

func resourceTencentCloudCosBucketInventory() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCosBucketInventoryCreate,
		Read:   resourceTencentCloudCosBucketInventoryRead,
		Update: resourceTencentCloudCosBucketInventoryUpdate,
		Delete: resourceTencentCloudCosBucketInventoryDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Bucket name.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Inventory Name.",
			},
			"is_enabled": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Whether to enable the inventory. true or false.",
			},
			"included_object_versions": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Whether to include object versions in the inventory. All or No.",
			},
			"filter": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Filters objects prefixed with the specified value to analyze.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"prefix": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Prefix of the objects to analyze.",
						},
						"period": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Creation time range of the objects to analyze.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"start_time": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Creation start time of the objects to analyze. The parameter is a timestamp in seconds, for example, 1568688761.",
									},
									"end_time": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Creation end time of the objects to analyze. The parameter is a timestamp in seconds, for example, 1568688762.",
									},
								},
							},
						},
					},
				},
			},
			"optional_fields": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Description: "Analysis items to include in the inventory result	.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"fields": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "Optional analysis items to include in the inventory result. The optional fields include Size, LastModifiedDate, StorageClass, ETag, IsMultipartUploaded, ReplicationStatus, Tag, Crc64, and x-cos-meta-*.",
						},
					},
				},
			},
			"schedule": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Required:    true,
				Description: "Inventory job cycle.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"frequency": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Frequency of the inventory job. Enumerated values: Daily, Weekly.",
						},
					},
				},
			},
			"destination": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Required:    true,
				Description: "Information about the inventory result destination.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bucket": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Bucket name.",
						},
						"account_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "ID of the bucket owner.",
						},
						"prefix": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Prefix of the inventory result.",
						},
						"format": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Format of the inventory result. Valid value: CSV.",
						},
						"encryption": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Server-side encryption for the inventory result.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"sse_cos": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Encryption with COS-managed key. This field can be left empty.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudCosBucketInventoryCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cos_bucket_inventory.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	bucket := d.Get("bucket").(string)
	name := d.Get("name").(string)
	isEnabled := d.Get("is_enabled").(string)
	includedObjectVersions := d.Get("included_object_versions").(string)

	var filter cos.BucketInventoryFilter
	if v, ok := d.GetOk("filter"); ok && len(v.([]interface{})) != 0 {
		var period cos.BucketInventoryFilterPeriod
		filterMap := v.([]interface{})[0].(map[string]interface{})
		if v, ok := filterMap["period"]; ok && len(v.([]interface{})) > 0 {
			periodMap := v.([]interface{})[0].(map[string]interface{})
			if v, ok := periodMap["start_time"]; ok && v.(string) != "" {
				vStr, err := strconv.ParseInt(v.(string), 10, 64)
				if err != nil {
					return err
				}
				period.StartTime = vStr
			}
			if v, ok := periodMap["end_time"]; ok && v.(string) != "" {
				vStr, err := strconv.ParseInt(v.(string), 10, 64)
				if err != nil {
					return err
				}
				period.EndTime = vStr
			}
			filter.Period = &period
		}
		if v, ok := filterMap["prefix"]; ok {
			filter.Prefix = v.(string)
		}
	}
	var optionalFields cos.BucketInventoryOptionalFields
	if v, ok := d.GetOk("optional_fields"); ok && len(v.([]interface{})) != 0 {
		optionalFieldsMap := v.([]interface{})[0].(map[string]interface{})
		if v, ok := optionalFieldsMap["fields"]; ok {
			optionalFields.BucketInventoryFields = make([]string, 0)
			for _, field := range v.(*schema.Set).List() {
				optionalFields.BucketInventoryFields = append(optionalFields.BucketInventoryFields, field.(string))
			}
		}
	}

	var schedule cos.BucketInventorySchedule
	if v, ok := d.GetOk("schedule"); ok && len(v.([]interface{})) != 0 {
		scheduleMap := v.([]interface{})[0].(map[string]interface{})
		if v, ok := scheduleMap["frequency"]; ok {
			schedule.Frequency = v.(string)
		}
	}

	var destination cos.BucketInventoryDestination
	if v, ok := d.GetOk("destination"); ok && len(v.([]interface{})) != 0 {
		destinationMap := v.([]interface{})[0].(map[string]interface{})
		if v, ok := destinationMap["bucket"]; ok {
			destination.Bucket = v.(string)
		}
		if v, ok := destinationMap["account_id"]; ok {
			destination.AccountId = v.(string)
		}
		if v, ok := destinationMap["prefix"]; ok {
			destination.Prefix = v.(string)
		}
		if v, ok := destinationMap["format"]; ok {
			destination.Format = v.(string)
		}
		if v, ok := destinationMap["encryption"]; ok && len(v.([]interface{})) > 0 {
			encryptionMap := v.([]interface{})[0].(map[string]interface{})
			if v, ok := encryptionMap["sse_cos"]; ok {
				destination.Encryption = &cos.BucketInventoryEncryption{
					SSECOS: v.(string),
				}

			}
		}
	}

	opt := &cos.BucketPutInventoryOptions{
		ID:                     name,
		IsEnabled:              isEnabled,
		IncludedObjectVersions: includedObjectVersions,
		Filter:                 &filter,
		OptionalFields:         &optionalFields,
		Schedule:               &schedule,
		Destination:            &destination,
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		req, _ := json.Marshal(opt)
		resp, e := meta.(*TencentCloudClient).apiV3Conn.UseTencentCosClient(bucket).Bucket.PutInventory(ctx, name, opt)
		responseBody, _ := json.Marshal(resp.Body)
		if e != nil {
			log.Printf("[DEBUG]%s api[PutInventory] success, request body [%s], response body [%s], err: [%s]\n", logId, req, responseBody, e.Error())
			return retryError(e)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cos bucketInventory failed, reason:%+v", logId, err)
		return err
	}
	d.SetId(bucket + FILED_SP + name)

	return resourceTencentCloudCosBucketInventoryRead(d, meta)
}

func resourceTencentCloudCosBucketInventoryRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cos_bucket_inventory.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	bucket := idSplit[0]
	name := idSplit[1]
	result, _, err := meta.(*TencentCloudClient).apiV3Conn.UseTencentCosClient(bucket).Bucket.GetInventory(ctx, name)
	if err != nil {
		log.Printf("[CRITAL]%s get cos bucketInventory failed, reason:%+v", logId, err)
		return err
	}
	_ = d.Set("bucket", bucket)
	_ = d.Set("name", name)
	_ = d.Set("is_enabled", result.IsEnabled)
	_ = d.Set("included_object_versions", result.IncludedObjectVersions)
	filterMap := make(map[string]interface{})
	if result.Filter != nil {
		filterMap["prefix"] = result.Filter.Prefix
		periodMap := make(map[string]interface{})
		if result.Filter.Period != nil {
			if result.Filter.Period.StartTime != 0 {
				periodMap["start_time"] = strconv.FormatInt(result.Filter.Period.StartTime, 10)
			}
			if result.Filter.Period.EndTime != 0 {
				periodMap["end_time"] = strconv.FormatInt(result.Filter.Period.EndTime, 10)
			}
			filterMap["period"] = []interface{}{periodMap}
		}
	}
	_ = d.Set("filter", []interface{}{filterMap})
	optionalFieldsMap := make(map[string]interface{})
	if result.OptionalFields != nil {
		fields := make([]string, 0)
		if result.OptionalFields.BucketInventoryFields != nil {
			fields = append(fields, result.OptionalFields.BucketInventoryFields...)
			optionalFieldsMap["fields"] = fields
		}
	}
	_ = d.Set("optional_fields", []interface{}{optionalFieldsMap})

	scheduleMap := make(map[string]interface{})
	if result.Schedule != nil {
		scheduleMap["frequency"] = result.Schedule.Frequency
	}
	_ = d.Set("schedule", []interface{}{scheduleMap})

	destinationMap := make(map[string]interface{})
	if result.Destination != nil {
		destinationMap["bucket"] = result.Destination.Bucket
		destinationMap["account_id"] = result.Destination.AccountId
		destinationMap["prefix"] = result.Destination.Prefix
		destinationMap["format"] = result.Destination.Format
		if result.Destination.Encryption != nil && result.Destination.Encryption.SSECOS != "" {
			encryptionMap := make(map[string]interface{})

			encryptionMap["sse_cos"] = result.Destination.Encryption.SSECOS
			destinationMap["encryption"] = []interface{}{encryptionMap}

		}
	}
	_ = d.Set("destination", []interface{}{destinationMap})

	return nil
}

func resourceTencentCloudCosBucketInventoryUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cos_bucket_inventory.update")()
	defer inconsistentCheck(d, meta)()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	bucket := idSplit[0]
	name := idSplit[1]
	if !d.HasChange("is_enabled") && !d.HasChange("included_object_versions") && !d.HasChange("filter") && !d.HasChange("optional_fields") && !d.HasChange("schedule") && !d.HasChange("destination") {
		return resourceTencentCloudCosBucketInventoryRead(d, meta)
	}
	isEnabled := d.Get("is_enabled").(string)
	includedObjectVersions := d.Get("included_object_versions").(string)

	var filter cos.BucketInventoryFilter
	if v, ok := d.GetOk("filter"); ok && len(v.([]interface{})) != 0 {
		var period cos.BucketInventoryFilterPeriod
		filterMap := v.([]interface{})[0].(map[string]interface{})
		if v, ok := filterMap["period"]; ok && len(v.([]interface{})) > 0 {
			periodMap := v.([]interface{})[0].(map[string]interface{})
			if v, ok := periodMap["start_time"]; ok && v.(string) != "" {
				vStr, err := strconv.ParseInt(v.(string), 10, 64)
				if err != nil {
					return err
				}
				period.StartTime = vStr
			}
			if v, ok := periodMap["end_time"]; ok && v.(string) != "" {
				vStr, err := strconv.ParseInt(v.(string), 10, 64)
				if err != nil {
					return err
				}
				period.EndTime = vStr
			}
			filter.Period = &period
		}
		if v, ok := filterMap["prefix"]; ok {
			filter.Prefix = v.(string)
		}
	}
	var optionalFields cos.BucketInventoryOptionalFields
	if v, ok := d.GetOk("optional_fields"); ok && len(v.([]interface{})) != 0 {
		optionalFieldsMap := v.([]interface{})[0].(map[string]interface{})
		if v, ok := optionalFieldsMap["fields"]; ok {
			optionalFields.BucketInventoryFields = make([]string, 0)
			for _, field := range v.(*schema.Set).List() {
				optionalFields.BucketInventoryFields = append(optionalFields.BucketInventoryFields, field.(string))
			}
		}
	}

	var schedule cos.BucketInventorySchedule
	if v, ok := d.GetOk("schedule"); ok && len(v.([]interface{})) != 0 {
		scheduleMap := v.([]interface{})[0].(map[string]interface{})
		if v, ok := scheduleMap["frequency"]; ok {
			schedule.Frequency = v.(string)
		}
	}

	var destination cos.BucketInventoryDestination
	if v, ok := d.GetOk("destination"); ok && len(v.([]interface{})) != 0 {
		destinationMap := v.([]interface{})[0].(map[string]interface{})
		if v, ok := destinationMap["bucket"]; ok {
			destination.Bucket = v.(string)
		}
		if v, ok := destinationMap["account_id"]; ok {
			destination.AccountId = v.(string)
		}
		if v, ok := destinationMap["prefix"]; ok {
			destination.Prefix = v.(string)
		}
		if v, ok := destinationMap["format"]; ok {
			destination.Format = v.(string)
		}
		if v, ok := destinationMap["encryption"]; ok && len(v.([]interface{})) > 0 {
			encryptionMap := v.([]interface{})[0].(map[string]interface{})
			if v, ok := encryptionMap["sse_cos"]; ok {
				destination.Encryption = &cos.BucketInventoryEncryption{
					SSECOS: v.(string),
				}

			}
		}
	}

	opt := &cos.BucketPutInventoryOptions{
		ID:                     name,
		IsEnabled:              isEnabled,
		IncludedObjectVersions: includedObjectVersions,
		Filter:                 &filter,
		OptionalFields:         &optionalFields,
		Schedule:               &schedule,
		Destination:            &destination,
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		req, _ := json.Marshal(opt)
		resp, e := meta.(*TencentCloudClient).apiV3Conn.UseTencentCosClient(bucket).Bucket.PutInventory(ctx, name, opt)
		responseBody, _ := json.Marshal(resp.Body)
		if e != nil {
			log.Printf("[DEBUG]%s api[PutInventory] success, request body [%s], response body [%s], err: [%s]\n", logId, req, responseBody, e.Error())
			return retryError(e)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cos bucketInventory failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCosBucketInventoryRead(d, meta)
}

func resourceTencentCloudCosBucketInventoryDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cos_bucket_inventory.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	bucket := idSplit[0]
	name := idSplit[1]

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		resp, e := meta.(*TencentCloudClient).apiV3Conn.UseTencentCosClient(bucket).Bucket.DeleteInventory(ctx, name)
		if e != nil {
			log.Printf("[CRITAL][retry]%s api[%s] fail, resp body [%s], reason[%s]\n",
				logId, "DeleteInventory ", resp.Body, e.Error())
			return retryError(e)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete cos bucketInventory failed, reason:%+v", logId, err)
		return err
	}

	return nil
}
