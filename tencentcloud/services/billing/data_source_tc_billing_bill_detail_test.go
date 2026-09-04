package billing_test

import (
	"os"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	billingv20180709 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/billing/v20180709"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/billing"
)

// go test ./tencentcloud/services/billing/ -run "TestBillingBillDetailDS" -v -count=1 -gcflags="all=-l"

type mockMetaBillingBillDetailDS struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaBillingBillDetailDS) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaBillingBillDetailDS{}

func newMockMetaBillingBillDetailDS() *mockMetaBillingBillDetailDS {
	return &mockMetaBillingBillDetailDS{client: &connectivity.TencentCloudClient{}}
}

func billingBillDetailPtrStr(s string) *string {
	return &s
}

func billingBillDetailPtrUint64(n uint64) *uint64 {
	return &n
}

func TestBillingBillDetailDS_ReadBasic(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	billingClient := &billingv20180709.Client{}
	patches.ApplyMethodReturn(newMockMetaBillingBillDetailDS().client, "UseBillingV20180709Client", billingClient)

	patches.ApplyMethodFunc(billingClient, "DescribeBillDetail", func(request *billingv20180709.DescribeBillDetailRequest) (*billingv20180709.DescribeBillDetailResponse, error) {
		resp := billingv20180709.NewDescribeBillDetailResponse()
		resp.Response = &billingv20180709.DescribeBillDetailResponseParams{
			Total: billingBillDetailPtrUint64(1),
			DetailSet: []*billingv20180709.BillDetail{
				{
					BusinessCodeName: billingBillDetailPtrStr("CVM"),
					ProductCodeName:  billingBillDetailPtrStr("CVM-Standard S1"),
					PayModeName:      billingBillDetailPtrStr("postPay"),
					ProjectName:      billingBillDetailPtrStr("default"),
					RegionName:       billingBillDetailPtrStr("Guangzhou"),
					ZoneName:         billingBillDetailPtrStr("Guangzhou-3"),
					ResourceId:       billingBillDetailPtrStr("ins-xxxxxxxx"),
					ResourceName:     billingBillDetailPtrStr("test-instance"),
					ActionTypeName:   billingBillDetailPtrStr("postpay deduction"),
					OrderId:          billingBillDetailPtrStr(""),
					BillId:           billingBillDetailPtrStr("bill-001"),
					PayTime:          billingBillDetailPtrStr("2024-01-01 00:00:00"),
					FeeBeginTime:     billingBillDetailPtrStr("2024-01-01 00:00:00"),
					FeeEndTime:       billingBillDetailPtrStr("2024-01-02 00:00:00"),
					PayerUin:         billingBillDetailPtrStr("100000000001"),
					OwnerUin:         billingBillDetailPtrStr("100000000001"),
					OperateUin:       billingBillDetailPtrStr("100000000001"),
					BusinessCode:     billingBillDetailPtrStr("p_cvm"),
					ProductCode:      billingBillDetailPtrStr("p_cvm_s1"),
					ActionType:       billingBillDetailPtrStr("postPay"),
					RegionId:         billingBillDetailPtrStr("1"),
					Formula:          billingBillDetailPtrStr("cost=price*usage"),
					FormulaUrl:       billingBillDetailPtrStr("https://cloud.tencent.com/document/product/555/35762"),
					BillDay:          billingBillDetailPtrStr("2024-01-01"),
					BillMonth:        billingBillDetailPtrStr("2024-01"),
					Id:               billingBillDetailPtrStr("detail-001"),
					RegionType:       billingBillDetailPtrStr("domestic"),
					RegionTypeName:   billingBillDetailPtrStr("domestic"),
					ComponentSet: []*billingv20180709.BillDetailComponent{
						{
							ComponentCodeName: billingBillDetailPtrStr("CPU"),
							ItemCodeName:      billingBillDetailPtrStr("CPU core"),
							SinglePrice:       billingBillDetailPtrStr("0.50"),
							PriceUnit:         billingBillDetailPtrStr("yuan/hour"),
							UsedAmount:        billingBillDetailPtrStr("1"),
							UsedAmountUnit:    billingBillDetailPtrStr("core"),
							TimeSpan:          billingBillDetailPtrStr("1"),
							TimeUnitName:      billingBillDetailPtrStr("hour"),
							Cost:              billingBillDetailPtrStr("0.50"),
							Discount:          billingBillDetailPtrStr("1.00"),
							RealCost:          billingBillDetailPtrStr("0.50"),
							ComponentConfig: []*billingv20180709.BillDetailComponentConfig{
								{
									Name:  billingBillDetailPtrStr("cpu"),
									Value: billingBillDetailPtrStr("1 core"),
								},
							},
						},
					},
					Tags: []*billingv20180709.BillTagInfo{
						{
							TagKey:   billingBillDetailPtrStr("env"),
							TagValue: billingBillDetailPtrStr("prod"),
						},
					},
					PriceInfo: []*string{billingBillDetailPtrStr("attr1")},
					AssociatedOrder: &billingv20180709.BillDetailAssociatedOrder{
						PrepayPurchase: billingBillDetailPtrStr("order-001"),
					},
				},
			},
			Context: billingBillDetailPtrStr("ctx-token"),
		}
		return resp, nil
	})

	meta := newMockMetaBillingBillDetailDS()
	res := billing.DataSourceTencentCloudBillingBillDetail()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"month":           "2024-01",
		"need_record_num": 1,
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	detailSet := d.Get("detail_set").([]interface{})
	assert.Len(t, detailSet, 1)

	detail0 := detailSet[0].(map[string]interface{})
	assert.Equal(t, "CVM", detail0["business_code_name"].(string))
	assert.Equal(t, "ins-xxxxxxxx", detail0["resource_id"].(string))
	assert.Equal(t, "2024-01", detail0["bill_month"].(string))

	componentSet := detail0["component_set"].([]interface{})
	assert.Len(t, componentSet, 1)
	component0 := componentSet[0].(map[string]interface{})
	assert.Equal(t, "CPU", component0["component_code_name"].(string))

	componentConfig := component0["component_config"].([]interface{})
	assert.Len(t, componentConfig, 1)
	config0 := componentConfig[0].(map[string]interface{})
	assert.Equal(t, "cpu", config0["name"].(string))

	tags := detail0["tags"].([]interface{})
	assert.Len(t, tags, 1)
	tag0 := tags[0].(map[string]interface{})
	assert.Equal(t, "env", tag0["tag_key"].(string))

	priceInfo := detail0["price_info"].([]interface{})
	assert.Len(t, priceInfo, 1)
	assert.Equal(t, "attr1", priceInfo[0].(string))

	associatedOrder := detail0["associated_order"].([]interface{})
	assert.Len(t, associatedOrder, 1)
	order0 := associatedOrder[0].(map[string]interface{})
	assert.Equal(t, "order-001", order0["prepay_purchase"].(string))

	assert.Equal(t, 1, d.Get("total").(int))
	assert.Equal(t, "ctx-token", d.Get("context").(string))
}

func TestBillingBillDetailDS_ReadWithEmptyResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	billingClient := &billingv20180709.Client{}
	patches.ApplyMethodReturn(newMockMetaBillingBillDetailDS().client, "UseBillingV20180709Client", billingClient)

	patches.ApplyMethodFunc(billingClient, "DescribeBillDetail", func(request *billingv20180709.DescribeBillDetailRequest) (*billingv20180709.DescribeBillDetailResponse, error) {
		resp := billingv20180709.NewDescribeBillDetailResponse()
		resp.Response = &billingv20180709.DescribeBillDetailResponseParams{
			Total:     billingBillDetailPtrUint64(0),
			DetailSet: []*billingv20180709.BillDetail{},
		}
		return resp, nil
	})

	meta := newMockMetaBillingBillDetailDS()
	res := billing.DataSourceTencentCloudBillingBillDetail()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"month": "2024-01",
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)

	detailSet := d.Get("detail_set").([]interface{})
	assert.Len(t, detailSet, 0)
	assert.Equal(t, 0, d.Get("total").(int))
}

func TestBillingBillDetailDS_ReadPagination(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	billingClient := &billingv20180709.Client{}
	patches.ApplyMethodReturn(newMockMetaBillingBillDetailDS().client, "UseBillingV20180709Client", billingClient)

	callCount := 0
	patches.ApplyMethodFunc(billingClient, "DescribeBillDetail", func(request *billingv20180709.DescribeBillDetailRequest) (*billingv20180709.DescribeBillDetailResponse, error) {
		callCount++
		resp := billingv20180709.NewDescribeBillDetailResponse()
		var detailSet []*billingv20180709.BillDetail
		if callCount == 1 {
			// first page returns 300 items, total is 301
			for i := 0; i < 300; i++ {
				detailSet = append(detailSet, &billingv20180709.BillDetail{
					BillId: billingBillDetailPtrStr("bill-page1"),
				})
			}
			resp.Response = &billingv20180709.DescribeBillDetailResponseParams{
				Total:     billingBillDetailPtrUint64(301),
				DetailSet: detailSet,
			}
		} else {
			// second page returns 1 item
			detailSet = append(detailSet, &billingv20180709.BillDetail{
				BillId: billingBillDetailPtrStr("bill-page2"),
			})
			resp.Response = &billingv20180709.DescribeBillDetailResponseParams{
				Total:     billingBillDetailPtrUint64(301),
				DetailSet: detailSet,
			}
		}
		return resp, nil
	})

	meta := newMockMetaBillingBillDetailDS()
	res := billing.DataSourceTencentCloudBillingBillDetail()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"month": "2024-01",
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, 2, callCount)

	detailSet := d.Get("detail_set").([]interface{})
	assert.Len(t, detailSet, 301)
	assert.Equal(t, 301, d.Get("total").(int))
}

func TestBillingBillDetailDS_ReadPaginationNoTotal(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	billingClient := &billingv20180709.Client{}
	patches.ApplyMethodReturn(newMockMetaBillingBillDetailDS().client, "UseBillingV20180709Client", billingClient)

	callCount := 0
	patches.ApplyMethodFunc(billingClient, "DescribeBillDetail", func(request *billingv20180709.DescribeBillDetailRequest) (*billingv20180709.DescribeBillDetailResponse, error) {
		callCount++
		resp := billingv20180709.NewDescribeBillDetailResponse()
		var detailSet []*billingv20180709.BillDetail
		if callCount == 1 {
			// first page returns 300 items, Total is NOT returned (need_record_num not set)
			for i := 0; i < 300; i++ {
				detailSet = append(detailSet, &billingv20180709.BillDetail{
					BillId: billingBillDetailPtrStr("bill-page1"),
				})
			}
		} else {
			// second page returns 1 item
			detailSet = append(detailSet, &billingv20180709.BillDetail{
				BillId: billingBillDetailPtrStr("bill-page2"),
			})
		}
		resp.Response = &billingv20180709.DescribeBillDetailResponseParams{
			DetailSet: detailSet,
		}
		return resp, nil
	})

	meta := newMockMetaBillingBillDetailDS()
	res := billing.DataSourceTencentCloudBillingBillDetail()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"month": "2024-01",
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, 2, callCount)

	detailSet := d.Get("detail_set").([]interface{})
	assert.Len(t, detailSet, 301)
}

func TestBillingBillDetailDS_ReadParamPassing(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	billingClient := &billingv20180709.Client{}
	patches.ApplyMethodReturn(newMockMetaBillingBillDetailDS().client, "UseBillingV20180709Client", billingClient)

	patches.ApplyMethodFunc(billingClient, "DescribeBillDetail", func(request *billingv20180709.DescribeBillDetailRequest) (*billingv20180709.DescribeBillDetailResponse, error) {
		// 断言所有 Optional 入参都被正确透传到 request（「有」分支）
		assert.Equal(t, "2024-01", *request.Month)
		assert.Equal(t, "2024-01-01 00:00:00", *request.BeginTime)
		assert.Equal(t, "2024-01-31 23:59:59", *request.EndTime)
		assert.Equal(t, int64(1), *request.NeedRecordNum)
		assert.Equal(t, "postPay", *request.PayMode)
		assert.Equal(t, "ins-xxxxxxxx", *request.ResourceId)
		assert.Equal(t, "按量计费日结", *request.ActionType)
		assert.Equal(t, int64(1002227), *request.ProjectId)
		assert.Equal(t, "p_cvm", *request.BusinessCode)
		assert.Equal(t, "ctx-input", *request.Context)
		assert.Equal(t, "909619400", *request.PayerUin)

		resp := billingv20180709.NewDescribeBillDetailResponse()
		resp.Response = &billingv20180709.DescribeBillDetailResponseParams{
			Total:     billingBillDetailPtrUint64(1),
			DetailSet: []*billingv20180709.BillDetail{{BillId: billingBillDetailPtrStr("bill-001")}},
			Context:   billingBillDetailPtrStr("ctx-output"),
		}
		return resp, nil
	})

	meta := newMockMetaBillingBillDetailDS()
	res := billing.DataSourceTencentCloudBillingBillDetail()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"month":           "2024-01",
		"begin_time":      "2024-01-01 00:00:00",
		"end_time":        "2024-01-31 23:59:59",
		"need_record_num": 1,
		"pay_mode":        "postPay",
		"resource_id":     "ins-xxxxxxxx",
		"action_type":     "按量计费日结",
		"project_id":      1002227,
		"business_code":   "p_cvm",
		"context":         "ctx-input",
		"payer_uin":       "909619400",
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, 1, d.Get("total").(int))
}

func TestBillingBillDetailDS_ReadMinimalArgs(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	billingClient := &billingv20180709.Client{}
	patches.ApplyMethodReturn(newMockMetaBillingBillDetailDS().client, "UseBillingV20180709Client", billingClient)

	patches.ApplyMethodFunc(billingClient, "DescribeBillDetail", func(request *billingv20180709.DescribeBillDetailRequest) (*billingv20180709.DescribeBillDetailResponse, error) {
		// 必选参数 offset/limit 由内部自动兜底，用户无需传
		assert.Equal(t, uint64(0), *request.Offset)
		assert.Equal(t, uint64(300), *request.Limit)
		// 空入参 {}，其余可选入参均为 nil（未被设置）
		assert.Nil(t, request.Month)
		assert.Nil(t, request.BeginTime)
		assert.Nil(t, request.EndTime)
		assert.Nil(t, request.NeedRecordNum)
		assert.Nil(t, request.PayMode)
		assert.Nil(t, request.ResourceId)
		assert.Nil(t, request.ActionType)
		assert.Nil(t, request.ProjectId)
		assert.Nil(t, request.BusinessCode)
		assert.Nil(t, request.Context)
		assert.Nil(t, request.PayerUin)

		resp := billingv20180709.NewDescribeBillDetailResponse()
		resp.Response = &billingv20180709.DescribeBillDetailResponseParams{
			Total:     billingBillDetailPtrUint64(1),
			DetailSet: []*billingv20180709.BillDetail{{BillId: billingBillDetailPtrStr("bill-001")}},
		}
		return resp, nil
	})

	meta := newMockMetaBillingBillDetailDS()
	res := billing.DataSourceTencentCloudBillingBillDetail()
	// 最小测试单元：空入参 {}
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())
	assert.Equal(t, 1, d.Get("total").(int))
}

func TestBillingBillDetailDS_ReadPaginationWithContext(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	billingClient := &billingv20180709.Client{}
	patches.ApplyMethodReturn(newMockMetaBillingBillDetailDS().client, "UseBillingV20180709Client", billingClient)

	callCount := 0
	patches.ApplyMethodFunc(billingClient, "DescribeBillDetail", func(request *billingv20180709.DescribeBillDetailRequest) (*billingv20180709.DescribeBillDetailResponse, error) {
		callCount++
		resp := billingv20180709.NewDescribeBillDetailResponse()

		var detailSet []*billingv20180709.BillDetail
		if callCount == 1 {
			// 第一页：offset=0 + 用户传入的 context
			assert.Equal(t, uint64(0), *request.Offset)
			assert.Equal(t, "ctx-input", *request.Context)
			for i := 0; i < 300; i++ {
				detailSet = append(detailSet, &billingv20180709.BillDetail{BillId: billingBillDetailPtrStr("bill-page1")})
			}
			resp.Response = &billingv20180709.DescribeBillDetailResponseParams{
				DetailSet: detailSet,
				Context:   billingBillDetailPtrStr("ctx-page2"),
			}
		} else {
			// 第二页：offset=300 + 上一页返回的 context
			assert.Equal(t, uint64(300), *request.Offset)
			assert.Equal(t, "ctx-page2", *request.Context)
			detailSet = append(detailSet, &billingv20180709.BillDetail{BillId: billingBillDetailPtrStr("bill-page2")})
			resp.Response = &billingv20180709.DescribeBillDetailResponseParams{
				DetailSet: detailSet,
			}
		}
		return resp, nil
	})

	meta := newMockMetaBillingBillDetailDS()
	res := billing.DataSourceTencentCloudBillingBillDetail()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"month":   "2024-01",
		"context": "ctx-input",
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, 2, callCount)

	detailSet := d.Get("detail_set").([]interface{})
	assert.Len(t, detailSet, 301)
}

func TestBillingBillDetailDS_ReadWithResultOutputFile(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	billingClient := &billingv20180709.Client{}
	patches.ApplyMethodReturn(newMockMetaBillingBillDetailDS().client, "UseBillingV20180709Client", billingClient)

	patches.ApplyMethodFunc(billingClient, "DescribeBillDetail", func(request *billingv20180709.DescribeBillDetailRequest) (*billingv20180709.DescribeBillDetailResponse, error) {
		resp := billingv20180709.NewDescribeBillDetailResponse()
		resp.Response = &billingv20180709.DescribeBillDetailResponseParams{
			Total: billingBillDetailPtrUint64(1),
			DetailSet: []*billingv20180709.BillDetail{
				{
					BillId:          billingBillDetailPtrStr("bill-001"),
					ProductCodeName: billingBillDetailPtrStr("CVM-Standard S1"),
					ResourceId:      billingBillDetailPtrStr("ins-xxxxxxxx"),
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaBillingBillDetailDS()
	res := billing.DataSourceTencentCloudBillingBillDetail()
	tmpFile := "/tmp/billing_bill_detail_test_result.json"
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"month":              "2024-01",
		"result_output_file": tmpFile,
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)

	content, readErr := os.ReadFile(tmpFile)
	assert.NoError(t, readErr)
	assert.Contains(t, string(content), "bill-001")
	assert.Contains(t, string(content), "CVM-Standard S1")
	assert.Contains(t, string(content), "ins-xxxxxxxx")

	_ = os.Remove(tmpFile)
}

func TestBillingBillDetailDS_Schema(t *testing.T) {
	res := billing.DataSourceTencentCloudBillingBillDetail()

	assert.NotNil(t, res)
	assert.Contains(t, res.Schema, "detail_set")
	assert.Contains(t, res.Schema, "total")
	assert.Contains(t, res.Schema, "result_output_file")
	assert.Contains(t, res.Schema, "month")
	assert.Contains(t, res.Schema, "begin_time")
	assert.Contains(t, res.Schema, "end_time")

	monthSchema := res.Schema["month"]
	assert.Equal(t, schema.TypeString, monthSchema.Type)
	assert.True(t, monthSchema.Optional)

	totalSchema := res.Schema["total"]
	assert.Equal(t, schema.TypeInt, totalSchema.Type)
	assert.True(t, totalSchema.Computed)

	detailSetSchema := res.Schema["detail_set"]
	assert.Equal(t, schema.TypeList, detailSetSchema.Type)
	assert.True(t, detailSetSchema.Computed)

	elemRes := detailSetSchema.Elem.(*schema.Resource)
	assert.Contains(t, elemRes.Schema, "component_set")
	assert.Contains(t, elemRes.Schema, "tags")
	assert.Contains(t, elemRes.Schema, "associated_order")
	assert.Contains(t, elemRes.Schema, "price_info")
	assert.Contains(t, elemRes.Schema, "bill_month")

	componentSetSchema := elemRes.Schema["component_set"]
	assert.Equal(t, schema.TypeList, componentSetSchema.Type)
	componentElemRes := componentSetSchema.Elem.(*schema.Resource)
	assert.Contains(t, componentElemRes.Schema, "component_config")
	assert.Contains(t, componentElemRes.Schema, "single_price")

	associatedOrderSchema := elemRes.Schema["associated_order"]
	assert.Equal(t, schema.TypeList, associatedOrderSchema.Type)
	orderElemRes := associatedOrderSchema.Elem.(*schema.Resource)
	assert.Contains(t, orderElemRes.Schema, "prepay_purchase")
	assert.Contains(t, orderElemRes.Schema, "prepay_renew")
	assert.Contains(t, orderElemRes.Schema, "prepay_modify_up")
	assert.Contains(t, orderElemRes.Schema, "reverse_order")
	assert.Contains(t, orderElemRes.Schema, "new_order")
	assert.Contains(t, orderElemRes.Schema, "original")
}
