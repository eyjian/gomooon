# 写在前面

微信支付官方的 SDK 没有提供商家转账到零钱等的下载转账电子回单和下载转账账单接口。

实现分为三大部分，但均为遵照[微信支付官方的API文档](https://pay.weixin.qq.com/wiki/doc/apiv3/wxpay/pages/index.shtml)实现：

* 第一部分：运营工具下的商家转账到零钱
* 第二部分：支付产品下的资金/交易账单
* 第三部分：扩展工具下的分账

# 运营工具下的商家转账到零钱

转账电子回单分为两种：

* 转账批次（batch）电子回单：含批次中的所有转账记录
* 转账明细（detail）电子回单：仅批次中单笔转账记录

一个转账单可批量给多人转，对应的转账账单存在以上两种。

## apply_change_bill_receipt

申请商家转账到零钱的转账电子回单

## query_change_bill_receipt

查询商家转账到零钱的转账电子回单

## down_change_bill_receipt

下载商家转账到零钱的转账电子回单

# 支付产品下的资金/交易账单

## apply_bill

申请转账账单，支持交易账单和资金账单

## download_bill

下载账单，支持交易账单和资金账单

# 扩展工具下的分账

## apply_sharing_bill

申请分账账单

## download_sharing_bill

下载分账账单
