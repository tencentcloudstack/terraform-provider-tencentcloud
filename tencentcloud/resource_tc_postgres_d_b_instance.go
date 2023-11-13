/*
Provides a resource to create a postgres d_b_instance

Example Usage

```hcl
resource "tencentcloud_postgres_d_b_instance" "d_b_instance" {
  spec_code = "cdb.pg.sh1.2g"
  storage = 10
  instance_count = 1
  period = 1
  zone = "ap-guangzhou-7"
  charset = "UTF8"
  admin_name = "user"
  admin_password = "password!@#123ABCabc"
  project_id = 0
  d_b_version = ""
  instance_charge_type = "POSTPAID_BY_HOUR"
  auto_voucher =
  voucher_ids =
  vpc_id = "vpc-xxxx"
  subnet_id = "subnet-xxxx"
  auto_renew_flag =
  activity_id =
  name = ""
  need_support_ipv6 =
  tag_list {
		tag_key = ""
		tag_value = ""

  }
  security_group_ids =
  d_b_major_version = ""
  d_b_kernel_version = ""
  d_b_node_set {
		role = ""
		zone = ""

  }
  need_support_t_d_e =
  k_m_s_key_id = ""
  k_m_s_region = ""
  d_b_engine = ""
  d_b_engine_config = ""
  tags = {
    "createdBy" = "terraform"
  }
}
```

Import

postgres d_b_instance can be imported using the id, e.g.

```
terraform import tencentcloud_postgres_d_b_instance.d_b_instance d_b_instance_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	postgres "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"time"
)

func resourceTencentCloudPostgresDBInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPostgresDBInstanceCreate,
		Read:   resourceTencentCloudPostgresDBInstanceRead,
		Update: resourceTencentCloudPostgresDBInstanceUpdate,
		Delete: resourceTencentCloudPostgresDBInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"spec_code": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Purchasable specification ID, which can be obtained through the SpecCode field in the returned value of the DescribeProductConfig API.",
			},

			"storage": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Instance storage capacity in GB.",
			},

			"instance_count": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "The number of instances purchased at a time. Value range:1-10.",
			},

			"period": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Valid period in months of purchased instances. Valid values:1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 24, 36. This parameter is set to 1 when the pay-as-you-go billing mode is used.",
			},

			"zone": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Availability zone ID, which can be obtained through the Zone field in the returned value of the DescribeZones API.",
			},

			"charset": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance character set. Valid values:UTF8, LATIN1.",
			},

			"admin_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance root account name.",
			},

			"admin_password": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance root account password.",
			},

			"project_id": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Project ID.",
			},

			"d_b_version": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "PostgreSQL version. If it is specified, an instance running the latest kernel of PostgreSQL DBVersion will be created. You must pass in at least one of the following parameters:DBVersion, DBMajorVersion, DBKernelVersion.",
			},

			"instance_charge_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Instance billing mode. Valid values:PREPAID (monthly subscription), POSTPAID_BY_HOUR (pay-as-you-go).",
			},

			"auto_voucher": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Whether to automatically use vouchers.Valid values:1(yes),0(no).Default value:0.",
			},

			"voucher_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Voucher ID list. Currently, you can specify only one voucher.",
			},

			"vpc_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "VPC ID.",
			},

			"subnet_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "ID of a subnet in the VPC specified by VpcId.",
			},

			"auto_renew_flag": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Renewal flag. Valid values:0 (manual renewal), 1 (auto-renewal). Default value:0.",
			},

			"activity_id": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Campaign ID.",
			},

			"name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Instance name.",
			},

			"need_support_ipv6": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Whether to support IPv6 address access. Valid values:1 (yes), 0 (no). Default value:0.",
			},

			"tag_list": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "The information of tags to be associated with instances. This parameter is left empty by default.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Tag key.",
						},
						"tag_value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Tag value.",
						},
					},
				},
			},

			"security_group_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Security group IDs.",
			},

			"d_b_major_version": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "PostgreSQL major version. Valid values:10, 11, 12, 13. If it is specified, an instance running the latest kernel of PostgreSQL DBMajorVersion will be created. You must pass in at least one of the following parameters:DBMajorVersion, DBVersion, DBKernelVersion.",
			},

			"d_b_kernel_version": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "PostgreSQL kernel version. If it is specified, an instance running the latest kernel of PostgreSQL DBKernelVersion will be created. You must pass in one of the following parameters:DBKernelVersion, DBVersion, DBMajorVersion.",
			},

			"d_b_node_set": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Instance node information, which is required if you purchase a multi-AZ deployed instance.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"role": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Node type. Valid values:Primary; Standby.",
						},
						"zone": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "AZ where the node resides, such as ap-guangzhou-1.",
						},
					},
				},
			},

			"need_support_t_d_e": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Whether to support transparent data encryption. Valid values:1 (yes), 0 (no). Default value:0.",
			},

			"k_m_s_key_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "KeyId of custom key, which is required if you select custom key encryption. It is also the unique CMK identifier.",
			},

			"k_m_s_region": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The region where the KMS service is enabled. When KMSRegion is left empty, the KMS of the current region will be enabled by default. If the current region is not supported, you need to select another region supported by KMS.",
			},

			"d_b_engine": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Database engine. Valid values:1. postgresql (TencentDB for PostgreSQL)2. mssql_compatible（MSSQL compatible-TencentDB for PostgreSQL)Default value: postgresql.",
			},

			"d_b_engine_config": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Configuration information of database engine in the following format:{$key1:$value1, $key2:$value2}Valid values:1. mssql_compatible engine：migrationMode: Database mode. Valid values: single-db (single-database mode), multi-db (multi-database mode). Default value: single-db.defaultLocale: Default locale, which can’t be modified after the initialization. Default value: en_US. Valid values:af_ZA, sq_AL, ar_DZ, ar_BH, ar_EG, ar_IQ, ar_JO, ar_KW, ar_LB, ar_LY, ar_MA, ar_OM, ar_QA, ar_SA, ar_SY, ar_TN, ar_AE, ar_YE, hy_AM, az_Cyrl_AZ, az_Latn_AZ, eu_ES, be_BY, bg_BG, ca_ES, zh_HK, zh_MO, zh_CN, zh_SG, zh_TW, hr_HR, cs_CZ, da_DK, nl_BE, nl_NL, en_AU, en_BZ, en_CA, en_IE, en_JM, en_NZ, en_PH, en_ZA, en_TT, en_GB, en_US, en_ZW, et_EE, fo_FO, fa_IR, fi_FI, fr_BE, fr_CA, fr_FR, fr_LU, fr_MC, fr_CH, mk_MK, ka_GE, de_AT, de_DE, de_LI, de_LU, de_CH, el_GR, gu_IN, he_IL, hi_IN, hu_HU, is_IS, id_ID, it_IT, it_CH, ja_JP, kn_IN, kok_IN, ko_KR, ky_KG, lv_LV, lt_LT, ms_BN, ms_MY, mr_IN, mn_MN, nb_NO, nn_NO, pl_PL, pt_BR, pt_PT, pa_IN, ro_RO, ru_RU, sa_IN, sr_Cyrl_RS, sr_Latn_RS, sk_SK, sl_SI, es_AR, es_BO, es_CL, es_CO, es_CR, es_DO, es_EC, es_SV, es_GT, es_HN, es_MX, es_NI, es_PA, es_PY,es_PE, es_PR, es_ES, es_TRADITIONAL, es_UY, es_VE, sw_KE, sv_FI, sv_SE, tt_RU, te_IN, th_TH, tr_TR, uk_UA, ur_IN, ur_PK, uz_Cyrl_UZ, uz_Latn_UZ, vi_VN.serverCollationName: Name of collation rule, which can’t be modified after the initialization. Default value: sql_latin1_general_cp1_ci_as. Valid values:bbf_unicode_general_ci_as, bbf_unicode_cp1_ci_as, bbf_unicode_CP1250_ci_as, bbf_unicode_CP1251_ci_as, bbf_unicode_cp1253_ci_as, bbf_unicode_cp1254_ci_as, bbf_unicode_cp1255_ci_as, bbf_unicode_cp1256_ci_as, bbf_unicode_cp1257_ci_as, bbf_unicode_cp1258_ci_as, bbf_unicode_cp874_ci_as, sql_latin1_general_cp1250_ci_as, sql_latin1_general_cp1251_ci_as, sql_latin1_general_cp1_ci_as, sql_latin1_general_cp1253_ci_as, sql_latin1_general_cp1254_ci_as, sql_latin1_general_cp1255_ci_as,sql_latin1_general_cp1256_ci_as, sql_latin1_general_cp1257_ci_as, sql_latin1_general_cp1258_ci_as, chinese_prc_ci_as, cyrillic_general_ci_as, finnish_swedish_ci_as, french_ci_as, japanese_ci_as, korean_wansung_ci_as, latin1_general_ci_as, modern_spanish_ci_as, polish_ci_as, thai_ci_as, traditional_spanish_ci_as, turkish_ci_as, ukrainian_ci_as, vietnamese_ci_as.",
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag description list.",
			},
		},
	}
}

func resourceTencentCloudPostgresDBInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgres_d_b_instance.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request      = postgres.NewCreateInstancesRequest()
		response     = postgres.NewCreateInstancesResponse()
		dBInstanceId string
	)
	if v, ok := d.GetOk("spec_code"); ok {
		request.SpecCode = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("storage"); ok {
		request.Storage = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("instance_count"); ok {
		request.InstanceCount = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("period"); ok {
		request.Period = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("zone"); ok {
		request.Zone = helper.String(v.(string))
	}

	if v, ok := d.GetOk("charset"); ok {
		request.Charset = helper.String(v.(string))
	}

	if v, ok := d.GetOk("admin_name"); ok {
		request.AdminName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("admin_password"); ok {
		request.AdminPassword = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("project_id"); ok {
		request.ProjectId = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("d_b_version"); ok {
		request.DBVersion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_charge_type"); ok {
		request.InstanceChargeType = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("auto_voucher"); ok {
		request.AutoVoucher = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("voucher_ids"); ok {
		voucherIdsSet := v.(*schema.Set).List()
		for i := range voucherIdsSet {
			voucherIds := voucherIdsSet[i].(string)
			request.VoucherIds = append(request.VoucherIds, &voucherIds)
		}
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		request.VpcId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("subnet_id"); ok {
		request.SubnetId = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("auto_renew_flag"); ok {
		request.AutoRenewFlag = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("activity_id"); ok {
		request.ActivityId = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("need_support_ipv6"); ok {
		request.NeedSupportIpv6 = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("tag_list"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			tag := postgres.Tag{}
			if v, ok := dMap["tag_key"]; ok {
				tag.TagKey = helper.String(v.(string))
			}
			if v, ok := dMap["tag_value"]; ok {
				tag.TagValue = helper.String(v.(string))
			}
			request.TagList = append(request.TagList, &tag)
		}
	}

	if v, ok := d.GetOk("security_group_ids"); ok {
		securityGroupIdsSet := v.(*schema.Set).List()
		for i := range securityGroupIdsSet {
			securityGroupIds := securityGroupIdsSet[i].(string)
			request.SecurityGroupIds = append(request.SecurityGroupIds, &securityGroupIds)
		}
	}

	if v, ok := d.GetOk("d_b_major_version"); ok {
		request.DBMajorVersion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("d_b_kernel_version"); ok {
		request.DBKernelVersion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("d_b_node_set"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			dBNode := postgres.DBNode{}
			if v, ok := dMap["role"]; ok {
				dBNode.Role = helper.String(v.(string))
			}
			if v, ok := dMap["zone"]; ok {
				dBNode.Zone = helper.String(v.(string))
			}
			request.DBNodeSet = append(request.DBNodeSet, &dBNode)
		}
	}

	if v, ok := d.GetOkExists("need_support_t_d_e"); ok {
		request.NeedSupportTDE = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("k_m_s_key_id"); ok {
		request.KMSKeyId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("k_m_s_region"); ok {
		request.KMSRegion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("d_b_engine"); ok {
		request.DBEngine = helper.String(v.(string))
	}

	if v, ok := d.GetOk("d_b_engine_config"); ok {
		request.DBEngineConfig = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UsePostgresClient().CreateInstances(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create postgres DBInstance failed, reason:%+v", logId, err)
		return err
	}

	dBInstanceId = *response.Response.DBInstanceId
	d.SetId(dBInstanceId)

	service := PostgresService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"running"}, 120*readRetryTimeout, time.Second, service.PostgresDBInstanceStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := TagService{client: meta.(*TencentCloudClient).apiV3Conn}
		region := meta.(*TencentCloudClient).apiV3Conn.Region
		resourceName := fmt.Sprintf("qcs::postgres:%s:uin/:dbInstanceId/%s", region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudPostgresDBInstanceRead(d, meta)
}

func resourceTencentCloudPostgresDBInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgres_d_b_instance.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := PostgresService{client: meta.(*TencentCloudClient).apiV3Conn}

	dBInstanceId := d.Id()

	DBInstance, err := service.DescribePostgresDBInstanceById(ctx, dBInstanceId)
	if err != nil {
		return err
	}

	if DBInstance == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `PostgresDBInstance` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if DBInstance.SpecCode != nil {
		_ = d.Set("spec_code", DBInstance.SpecCode)
	}

	if DBInstance.Storage != nil {
		_ = d.Set("storage", DBInstance.Storage)
	}

	if DBInstance.InstanceCount != nil {
		_ = d.Set("instance_count", DBInstance.InstanceCount)
	}

	if DBInstance.Period != nil {
		_ = d.Set("period", DBInstance.Period)
	}

	if DBInstance.Zone != nil {
		_ = d.Set("zone", DBInstance.Zone)
	}

	if DBInstance.Charset != nil {
		_ = d.Set("charset", DBInstance.Charset)
	}

	if DBInstance.AdminName != nil {
		_ = d.Set("admin_name", DBInstance.AdminName)
	}

	if DBInstance.AdminPassword != nil {
		_ = d.Set("admin_password", DBInstance.AdminPassword)
	}

	if DBInstance.ProjectId != nil {
		_ = d.Set("project_id", DBInstance.ProjectId)
	}

	if DBInstance.DBVersion != nil {
		_ = d.Set("d_b_version", DBInstance.DBVersion)
	}

	if DBInstance.InstanceChargeType != nil {
		_ = d.Set("instance_charge_type", DBInstance.InstanceChargeType)
	}

	if DBInstance.AutoVoucher != nil {
		_ = d.Set("auto_voucher", DBInstance.AutoVoucher)
	}

	if DBInstance.VoucherIds != nil {
		_ = d.Set("voucher_ids", DBInstance.VoucherIds)
	}

	if DBInstance.VpcId != nil {
		_ = d.Set("vpc_id", DBInstance.VpcId)
	}

	if DBInstance.SubnetId != nil {
		_ = d.Set("subnet_id", DBInstance.SubnetId)
	}

	if DBInstance.AutoRenewFlag != nil {
		_ = d.Set("auto_renew_flag", DBInstance.AutoRenewFlag)
	}

	if DBInstance.ActivityId != nil {
		_ = d.Set("activity_id", DBInstance.ActivityId)
	}

	if DBInstance.Name != nil {
		_ = d.Set("name", DBInstance.Name)
	}

	if DBInstance.NeedSupportIpv6 != nil {
		_ = d.Set("need_support_ipv6", DBInstance.NeedSupportIpv6)
	}

	if DBInstance.TagList != nil {
		tagListList := []interface{}{}
		for _, tagList := range DBInstance.TagList {
			tagListMap := map[string]interface{}{}

			if DBInstance.TagList.TagKey != nil {
				tagListMap["tag_key"] = DBInstance.TagList.TagKey
			}

			if DBInstance.TagList.TagValue != nil {
				tagListMap["tag_value"] = DBInstance.TagList.TagValue
			}

			tagListList = append(tagListList, tagListMap)
		}

		_ = d.Set("tag_list", tagListList)

	}

	if DBInstance.SecurityGroupIds != nil {
		_ = d.Set("security_group_ids", DBInstance.SecurityGroupIds)
	}

	if DBInstance.DBMajorVersion != nil {
		_ = d.Set("d_b_major_version", DBInstance.DBMajorVersion)
	}

	if DBInstance.DBKernelVersion != nil {
		_ = d.Set("d_b_kernel_version", DBInstance.DBKernelVersion)
	}

	if DBInstance.DBNodeSet != nil {
		dBNodeSetList := []interface{}{}
		for _, dBNodeSet := range DBInstance.DBNodeSet {
			dBNodeSetMap := map[string]interface{}{}

			if DBInstance.DBNodeSet.Role != nil {
				dBNodeSetMap["role"] = DBInstance.DBNodeSet.Role
			}

			if DBInstance.DBNodeSet.Zone != nil {
				dBNodeSetMap["zone"] = DBInstance.DBNodeSet.Zone
			}

			dBNodeSetList = append(dBNodeSetList, dBNodeSetMap)
		}

		_ = d.Set("d_b_node_set", dBNodeSetList)

	}

	if DBInstance.NeedSupportTDE != nil {
		_ = d.Set("need_support_t_d_e", DBInstance.NeedSupportTDE)
	}

	if DBInstance.KMSKeyId != nil {
		_ = d.Set("k_m_s_key_id", DBInstance.KMSKeyId)
	}

	if DBInstance.KMSRegion != nil {
		_ = d.Set("k_m_s_region", DBInstance.KMSRegion)
	}

	if DBInstance.DBEngine != nil {
		_ = d.Set("d_b_engine", DBInstance.DBEngine)
	}

	if DBInstance.DBEngineConfig != nil {
		_ = d.Set("d_b_engine_config", DBInstance.DBEngineConfig)
	}

	tcClient := meta.(*TencentCloudClient).apiV3Conn
	tagService := &TagService{client: tcClient}
	tags, err := tagService.DescribeResourceTags(ctx, "postgres", "dbInstanceId", tcClient.Region, d.Id())
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudPostgresDBInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgres_d_b_instance.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := postgres.NewModifyDBInstanceSpecRequest()

	dBInstanceId := d.Id()

	request.DBInstanceId = &dBInstanceId

	immutableArgs := []string{"spec_code", "storage", "instance_count", "period", "zone", "charset", "admin_name", "admin_password", "project_id", "d_b_version", "instance_charge_type", "auto_voucher", "voucher_ids", "vpc_id", "subnet_id", "auto_renew_flag", "activity_id", "name", "need_support_ipv6", "tag_list", "security_group_ids", "d_b_major_version", "d_b_kernel_version", "d_b_node_set", "need_support_t_d_e", "k_m_s_key_id", "k_m_s_region", "d_b_engine", "d_b_engine_config"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("storage") {
		if v, ok := d.GetOkExists("storage"); ok {
			request.Storage = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("auto_voucher") {
		if v, ok := d.GetOkExists("auto_voucher"); ok {
			request.AutoVoucher = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("voucher_ids") {
		if v, ok := d.GetOk("voucher_ids"); ok {
			voucherIdsSet := v.(*schema.Set).List()
			for i := range voucherIdsSet {
				voucherIds := voucherIdsSet[i].(string)
				request.VoucherIds = append(request.VoucherIds, &voucherIds)
			}
		}
	}

	if d.HasChange("activity_id") {
		if v, ok := d.GetOkExists("activity_id"); ok {
			request.ActivityId = helper.IntInt64(v.(int))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UsePostgresClient().ModifyDBInstanceSpec(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update postgres DBInstance failed, reason:%+v", logId, err)
		return err
	}

	service := PostgresService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"running"}, 600*readRetryTimeout, time.Second, service.PostgresDBInstanceStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	if d.HasChange("tags") {
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		tcClient := meta.(*TencentCloudClient).apiV3Conn
		tagService := &TagService{client: tcClient}
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		resourceName := BuildTagResourceName("postgres", "dbInstanceId", tcClient.Region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
	}

	return resourceTencentCloudPostgresDBInstanceRead(d, meta)
}

func resourceTencentCloudPostgresDBInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgres_d_b_instance.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := PostgresService{client: meta.(*TencentCloudClient).apiV3Conn}
	dBInstanceId := d.Id()

	if err := service.DeletePostgresDBInstanceById(ctx, dBInstanceId); err != nil {
		return err
	}

	service := PostgresService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"running"}, 120*readRetryTimeout, time.Second, service.PostgresDBInstanceStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return nil
}
