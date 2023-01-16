package tencentcloud

var TRADE_RETRYABLE_ERROR = []string{
	"InternalError.TradeError",  //mysql
	"FailedOperation.PayFailed", //redis
}

// deal status: https://cloud.tencent.com/document/product/555/19179

var DEAL_STATUS_CODE = []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}

var DEAL_TERMINATE_STATUS_CODE = []int64{4, 5, 6, 7, 8, 9, 10, 11}
