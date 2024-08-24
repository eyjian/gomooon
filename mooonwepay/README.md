微信支付官方的 SDK 没有提供商家转账到零钱的下载转账电子回单接口，转账电子回单分为两种：

* 转账批次（batch）电子回单：含批次中的所有转账记录
* 转账明细（detail）电子回单：仅批次中单笔转账记录

一个转账单可批量给多人转，对应的转账账单存在以上两种。

# get_change_bill_receipt

商家转账到零钱的转账电子回单申请受理

# query_change_bill_receipt

查询商家转账到零钱的转账电子回单

# down_change_bill_receipt

下载商家转账到零钱的电子回单
