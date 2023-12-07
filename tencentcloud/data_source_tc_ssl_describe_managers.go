package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudSslDescribeManagers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudSslDescribeManagersRead,
		Schema: map[string]*schema.Schema{
			"company_id": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Company ID.",
			},

			"manager_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Manager&amp;#39;s name (will be abandoned), please use Searchkey.",
			},

			"manager_mail": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Vague query manager email (will be abandoned), please use Searchkey.",
			},

			"status": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Filter according to the status of the manager, and the value is available&amp;#39;None&amp;#39; Unable to submit review&amp;#39;Audit&amp;#39;, Asian Credit Review&amp;#39;Caaudit&amp;#39; CA review&amp;#39;OK&amp;#39; has been reviewed&amp;#39;Invalid&amp;#39; review failed&amp;#39;Expiring&amp;#39; is about to expire&amp;#39;Expired&amp;#39; expired.",
			},

			"search_key": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Manager&amp;#39;s surname/Manager name/mailbox/department precise matching.",
			},

			"managers": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Company Manager List.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Status: Audit: OK during the review: review passed inValid: expired expiRing: is about to expire Expired: expired.",
						},
						"manager_first_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Manager name.",
						},
						"manager_last_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Manager name.",
						},
						"manager_position": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Manager position.",
						},
						"manager_phone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Manager phone call.",
						},
						"manager_mail": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Manager mailbox.",
						},
						"manager_department": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Administrator department.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation timeNote: This field may return NULL, indicating that the valid value cannot be obtained.",
						},
						"domain_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of administrators.",
						},
						"cert_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of administrative certificates.",
						},
						"manager_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Manager ID.",
						},
						"expire_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Examine the validity expiration timeNote: This field may return NULL, indicating that the valid value cannot be obtained.",
						},
						"submit_audit_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The last time the review timeNote: This field may return NULL, indicating that the valid value cannot be obtained.",
						},
						"verify_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Examination timeNote: This field may return NULL, indicating that the valid value cannot be obtained.",
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

func dataSourceTencentCloudSslDescribeManagersRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_ssl_describe_managers.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	var companyId string
	paramMap := make(map[string]interface{})
	if v, _ := d.GetOk("company_id"); v != nil {
		companyId = helper.IntToStr(v.(int))
		paramMap["CompanyId"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("manager_name"); ok {
		paramMap["ManagerName"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("manager_mail"); ok {
		paramMap["ManagerMail"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("status"); ok {
		paramMap["Status"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("search_key"); ok {
		paramMap["SearchKey"] = helper.String(v.(string))
	}

	service := SslService{client: meta.(*TencentCloudClient).apiV3Conn}

	var managers []*ssl.ManagerInfo

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeSslDescribeManagersByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		managers = result
		return nil
	})
	if err != nil {
		return err
	}

	tmpList := make([]map[string]interface{}, 0, len(managers))

	if managers != nil {
		for _, managerInfo := range managers {
			managerInfoMap := map[string]interface{}{}

			if managerInfo.Status != nil {
				managerInfoMap["status"] = managerInfo.Status
			}

			if managerInfo.ManagerFirstName != nil {
				managerInfoMap["manager_first_name"] = managerInfo.ManagerFirstName
			}

			if managerInfo.ManagerLastName != nil {
				managerInfoMap["manager_last_name"] = managerInfo.ManagerLastName
			}

			if managerInfo.ManagerPosition != nil {
				managerInfoMap["manager_position"] = managerInfo.ManagerPosition
			}

			if managerInfo.ManagerPhone != nil {
				managerInfoMap["manager_phone"] = managerInfo.ManagerPhone
			}

			if managerInfo.ManagerMail != nil {
				managerInfoMap["manager_mail"] = managerInfo.ManagerMail
			}

			if managerInfo.ManagerDepartment != nil {
				managerInfoMap["manager_department"] = managerInfo.ManagerDepartment
			}

			if managerInfo.CreateTime != nil {
				managerInfoMap["create_time"] = managerInfo.CreateTime
			}

			if managerInfo.DomainCount != nil {
				managerInfoMap["domain_count"] = managerInfo.DomainCount
			}

			if managerInfo.CertCount != nil {
				managerInfoMap["cert_count"] = managerInfo.CertCount
			}

			if managerInfo.ManagerId != nil {
				managerInfoMap["manager_id"] = managerInfo.ManagerId
			}

			if managerInfo.ExpireTime != nil {
				managerInfoMap["expire_time"] = managerInfo.ExpireTime
			}

			if managerInfo.SubmitAuditTime != nil {
				managerInfoMap["submit_audit_time"] = managerInfo.SubmitAuditTime
			}

			if managerInfo.VerifyTime != nil {
				managerInfoMap["verify_time"] = managerInfo.VerifyTime
			}

			tmpList = append(tmpList, managerInfoMap)
		}

		_ = d.Set("managers", tmpList)
	}

	d.SetId(companyId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
