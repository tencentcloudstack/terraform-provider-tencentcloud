package tencentcloud

import (
	"context"
	"encoding/json"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudCosBucketInventorys() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCosBucketInventorysRead,

		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Bucket.",
			},
			"inventorys": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Multiple batch processing task information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Whether to enable the inventory. true or false.",
						},
						"is_enabled": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Whether to enable the inventory. true or false.",
						},
						"included_object_versions": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Whether to include object versions in the inventory. All or No.",
						},
						"filter": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Filters objects prefixed with the specified value to analyze.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"prefix": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Prefix of the objects to analyze.",
									},
									"period": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Creation time range of the objects to analyze.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"start_time": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Creation start time of the objects to analyze. The parameter is a timestamp in seconds, for example, 1568688761.",
												},
												"end_time": {
													Type:        schema.TypeString,
													Computed:    true,
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
							Optional: true,
							Description: "Analysis items to include in the inventory result	.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"fields": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "Optional analysis items to include in the inventory result. The optional fields include Size, LastModifiedDate, StorageClass, ETag, IsMultipartUploaded, ReplicationStatus, Tag, Crc64, and x-cos-meta-*.",
									},
								},
							},
						},
						"schedule": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Inventory job cycle.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"frequency": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Frequency of the inventory job. Enumerated values: Daily, Weekly.",
									},
								},
							},
						},
						"destination": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Information about the inventory result destination.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"bucket": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Bucket name.",
									},
									"account_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ID of the bucket owner.",
									},
									"prefix": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Prefix of the inventory result.",
									},
									"format": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Format of the inventory result. Valid value: CSV.",
									},
									"encryption": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Server-side encryption for the inventory result.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"sse_cos": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Encryption with COS-managed key. This field can be left empty.",
												},
											},
										},
									},
								},
							},
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

func dataSourceTencentCloudCosBucketInventorysRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_cos_bucket_inventorys.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	bucket := d.Get("bucket").(string)
	inventoryConfigurations := make([]map[string]interface{}, 0)
	token := ""
	ids := make([]string, 0)
	for {
		result, response, err := meta.(*TencentCloudClient).apiV3Conn.UseTencentCosClient(bucket).Bucket.ListInventoryConfigurations(ctx, token)
		responseBody, _ := json.Marshal(response.Body)
		log.Printf("[DEBUG]%s api[ListInventoryConfigurations] success, response body [%s]\n", logId, responseBody)
		if err != nil {
			return err
		}

		for _, item := range result.InventoryConfigurations {
			itemMap := make(map[string]interface{})
			itemMap["id"] = item.ID
			itemMap["is_enabled"] = item.IsEnabled
			itemMap["included_object_versions"] = item.IncludedObjectVersions

			filterMap := make(map[string]interface{})
			if item.Filter != nil {
				filterMap["prefix"] = item.Filter.Prefix
				periodMap := make(map[string]interface{})
				if item.Filter.Period != nil {
					if item.Filter.Period.StartTime != 0 {
						periodMap["start_time"] = strconv.FormatInt(item.Filter.Period.StartTime, 10)
					}
					if item.Filter.Period.EndTime != 0 {
						periodMap["end_time"] = strconv.FormatInt(item.Filter.Period.EndTime, 10)
					}
					filterMap["period"] = []interface{}{periodMap}
				}
				itemMap["filter"] = []interface{}{filterMap}
			}
			if item.OptionalFields != nil {
				optionalFieldsMap := make(map[string]interface{})
				fields := make([]string, 0)
				if item.OptionalFields.BucketInventoryFields != nil {
					fields = append(fields, item.OptionalFields.BucketInventoryFields...)
					optionalFieldsMap["fields"] = fields
				}
				itemMap["optional_fields"] = []interface{}{optionalFieldsMap}
			}

			if item.Schedule != nil {
				scheduleMap := make(map[string]interface{})
				scheduleMap["frequency"] = item.Schedule.Frequency
				itemMap["schedule"] = []interface{}{scheduleMap}
			}

			if item.Destination != nil {
				destinationMap := make(map[string]interface{})
				destinationMap["bucket"] = item.Destination.Bucket
				destinationMap["account_id"] = item.Destination.AccountId
				destinationMap["prefix"] = item.Destination.Prefix
				destinationMap["format"] = item.Destination.Format
				if item.Destination.Encryption != nil && item.Destination.Encryption.SSECOS != "" {
					encryptionMap := make(map[string]interface{})

					encryptionMap["sse_cos"] = item.Destination.Encryption.SSECOS
					destinationMap["encryption"] = []interface{}{encryptionMap}

				}
				itemMap["destination"] = []interface{}{destinationMap}
			}
			ids = append(ids, item.ID)
			inventoryConfigurations = append(inventoryConfigurations, itemMap)
		}
		if result.NextContinuationToken != "" {
			token = result.NextContinuationToken
		} else {
			break
		}
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	_ = d.Set("inventorys", inventoryConfigurations)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if err := writeToFile(output.(string), inventoryConfigurations); err != nil {
			return err
		}
	}

	return nil
}
