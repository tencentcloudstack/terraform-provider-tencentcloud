## 1. 修复 Read 函数

- [x] 1.1 将 Read 函数中 `if respData == nil {` 修改为 `if respData == nil || len(respData.Rules) == 0 {`
