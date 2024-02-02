package main

import (
    "context"
    "fmt"
    "github.com/go-pay/gopay"
    "github.com/go-pay/gopay/alipay"
    "github.com/go-pay/gopay/pkg/xlog"
)

// Alipay 支付创建
func Alipay(out_trade_no string, price string) string {
    client, err := alipay.NewClient("2021004130632232", "MIIEvwIBADANBgkqhkiG9w0BAQEFAASCBKkwggSlAgEAAoIBAQCxo+vNCXQi9B032U2OsrHAeY5DO+mCMUiNpvBa4Id8IIj3KP/fQb9+SwU+U0EvpDN/d2wf9ArvK4MUDWhcYZn50/eJ2hCqqoXJ+8KP6gAYNGzPHfb42pPRqW23Z/0QpIVhX30EQBTm326HWwCb9jsq6aIRHzAF9apQavTyzyqSVaAbFk3nHARBG0aTHo0P9Zp6vz1nbAgyPwQykuGLdYmtQDEHu53IA73aWHoV3Fqznftq/n6mo/PaJBe21LWm1jqe+EW8vgmLhdq+J1uymdmIxmCzqygxbP9q6ET608l41XtPUxtJX4EkfNucebJPT80t3nP/I3TYu9z/bSoRsqjNAgMBAAECggEBAIqG/SSIwcNcQMjDseKc4VbqtBEkUIWRrzWfwIHt5Fnc+VJc66SLnix7jlw7CnN/hhVZ6LzGUByQ/wgNKJwaFLwpGzmqDyM4FVsc9G3MGkTF5TBi+qy3r1xp1sCW3Fc1JTp4/4HoAyTDimsNgV0eWKevSA44FvgeyrEpp1kOSvGPfNhYVdT5Ak+zICApCwP1/CrTxQjSKPP/HHwQRbrhVVvu7jh1xzoTjIfecATi9D20txBWYm+fYEGBjR6T8FNHP5pjTr52ty84wz07uhTXqOX8gTbSDb/BJLa3/+aR7X4oNvsSa9iriqcs1xMq4KkB9bdZ8EzQ53HM3RHnJ20eP6ECgYEA5i/pIEdT4/vQr3KKyuAJ7cYWuD1fWkegA2FqyeLk1MmBCBNjnTCJWMicobKE+ihK2wNUDeuOHPG3sVuw0VzPjQ+lBUYovW9V1UFX4E2o/akNXyBKQJG3sbJWRee6s3aCqZVqbhnLO7t1DxDAId37H7xYDN/qObnOoXAG3reynhkCgYEAxY+GK0IgntEl8dSg3dE+XoO9X62fXze+8Z0QLhVbiKhwAqxIamJ1OjWWzGy5NtXOItvwkumw3ydJ4QXvZZpb8DzgOir9E9u18oh7WkVzL9kkynqfSk3XDXR2iI0s86VvrlSaValSyRsciRF1QEHznNJbjT4QRS45EjDxH/dlztUCgYEA0CMEPh6g2WXp5arBmv4HnEtgYcmEvcJECqp8f/48gbeOh7nYedrYZkJHduJP4U6rmOuihk+3Ga7rNWC+OiEcvuUlhuZQkjHov8Ks7fHq2yqQH7K30Tixi+jAn8cQB5QiQ6sKKHIEVYeEKlIwGK96kdChIUsapIXBNDJy09HwnYkCgYBzKyke/KzBiNFi+f5RcVK3jHsQVMnMm1XPyi0NgFvc/bxWgpKwmfcW2Piw8UzDv74sqiTDsEHwxRmXeXtGssaX9RUOM9NXCUU3PwMR69yrbx24f+VuTpRofpU/I3WqD65cZWuXNl9RZ2GqMig1Ln1S1XqTizO28KxKg4d9iB6shQKBgQDFfupwsResXN/v3bTZEMGJJqegRsnn0Tiw6XDbbSn7VgIIf4AXhyONpCgZ92F3z9f+wjbBPqdNP2pLVz1ZZ0Qnhjv2UkeWZbXl8VRYa+LGKwo7yTmoV8WcD6SOtIYwqmZyuwVTgXlaItYhSAeczx4+wwhqNwnQnQWbVKIVkEjR/A==", true)
    if err != nil {
        xlog.Error(err)
        return ""
    }

    client.DebugSwitch = gopay.DebugOn

    client.SetLocation(alipay.LocationShanghai).
        SetCharset(alipay.UTF8).
        SetSignType(alipay.RSA2).
        SetReturnUrl("https://www.fmm.ink").
        SetNotifyUrl("https://www.fmm.ink")

    bm := make(gopay.BodyMap)
    bm.Set("subject", "手机网站支付").
        Set("out_trade_no", out_trade_no).
        Set("total_amount", price).
        Set("timeout_express", "2m")

    aliRsp, err := client.TradeWapPay(context.Background(), bm)
    if err != nil {
        if bizErr, ok := alipay.IsBizError(err); ok {
            xlog.Errorf("%+v", bizErr)
            // do something
            return "11111"
        }
        xlog.Errorf("client.TradePay(%+v),err:%+v", bm, err)
        return "22222"
    }
    return aliRsp
}

type TradeQueryResponse struct {
    Response     *TradeQuery `json:"alipay_trade_query_response"`
    AlipayCertSn string      `json:"alipay_cert_sn,omitempty"`
    SignData     string      `json:"-"`
    Sign         string      `json:"sign"`
}
type ErrorResponse struct {
    Code    string `json:"code"`
    Msg     string `json:"msg"`
    SubCode string `json:"sub_code,omitempty"`
    SubMsg  string `json:"sub_msg,omitempty"`
}

type TradeFundBill struct {
    FundChannel string `json:"fund_channel,omitempty"` // 同步通知里是 fund_channel
    Amount      string `json:"amount,omitempty"`
    RealAmount  string `json:"real_amount,omitempty"`
    FundType    string `json:"fund_type,omitempty"`
}

type HbFqPayInfo struct {
    UserInstallNum string `json:"user_install_num,omitempty"`
}

type TradeSettleDetail struct {
    OperationType     string `json:"operation_type,omitempty"`
    OperationSerialNo string `json:"operation_serial_no,omitempty"`
    OperationDt       string `json:"operation_dt,omitempty"`
    TransOut          string `json:"trans_out,omitempty"`
    TransIn           string `json:"trans_in,omitempty"`
    Amount            string `json:"amount,omitempty"`
    OriTransOut       string `json:"ori_trans_out,omitempty"`
    OriTransIn        string `json:"ori_trans_in,omitempty"`
}

type TradeSettleInfo struct {
    TradeSettleDetailList []*TradeSettleDetail `json:"trade_settle_detail_list,omitempty"`
}

type BkAgentRespInfo struct {
    BindtrxId        string `json:"bindtrx_id,omitempty"`
    BindclrissrId    string `json:"bindclrissr_id,omitempty"`
    BindpyeracctbkId string `json:"bindpyeracctbk_id,omitempty"`
    BkpyeruserCode   string `json:"bkpyeruser_code,omitempty"`
    EstterLocation   string `json:"estter_location,omitempty"`
}

type ChargeInfo struct {
    ChargeFee               string    `json:"charge_fee,omitempty"`
    OriginalChargeFee       string    `json:"original_charge_fee,omitempty"`
    SwitchFeeRate           string    `json:"switch_fee_rate,omitempty"`
    IsRatingOnTradeReceiver string    `json:"is_rating_on_trade_receiver,omitempty"`
    IsRatingOnSwitch        string    `json:"is_rating_on_switch,omitempty"`
    ChargeType              string    `json:"charge_type,omitempty"`
    SubFeeDetailList        []*SubFee `json:"sub_fee_detail_list,omitempty"`
}

type SubFee struct {
    ChargeFee         string `json:"charge_fee,omitempty"`
    OriginalChargeFee string `json:"original_charge_fee,omitempty"`
    SwitchFeeRate     string `json:"switch_fee_rate,omitempty"`
}

type TradeQuery struct {
    ErrorResponse
    TradeNo               string           `json:"trade_no,omitempty"`
    OutTradeNo            string           `json:"out_trade_no,omitempty"`
    BuyerLogonId          string           `json:"buyer_logon_id,omitempty"`
    TradeStatus           string           `json:"trade_status,omitempty"`
    TotalAmount           string           `json:"total_amount,omitempty"`
    TransCurrency         string           `json:"trans_currency,omitempty"`
    SettleCurrency        string           `json:"settle_currency,omitempty"`
    SettleAmount          string           `json:"settle_amount,omitempty"`
    PayCurrency           string           `json:"pay_currency,omitempty"`
    PayAmount             string           `json:"pay_amount,omitempty"`
    SettleTransRate       string           `json:"settle_trans_rate,omitempty"`
    TransPayRate          string           `json:"trans_pay_rate,omitempty"`
    BuyerPayAmount        string           `json:"buyer_pay_amount,omitempty"`
    PointAmount           string           `json:"point_amount,omitempty"`
    InvoiceAmount         string           `json:"invoice_amount,omitempty"`
    SendPayDate           string           `json:"send_pay_date,omitempty"`
    ReceiptAmount         string           `json:"receipt_amount,omitempty"`
    StoreId               string           `json:"store_id,omitempty"`
    TerminalId            string           `json:"terminal_id,omitempty"`
    FundBillList          []*TradeFundBill `json:"fund_bill_list"`
    StoreName             string           `json:"store_name,omitempty"`
    BuyerUserId           string           `json:"buyer_user_id,omitempty"`
    BuyerOpenId           string           `json:"buyer_open_id,omitempty"`
    DiscountGoodsDetail   string           `json:"discount_goods_detail,omitempty"`
    IndustrySepcDetail    string           `json:"industry_sepc_detail,omitempty"`
    IndustrySepcDetailGov string           `json:"industry_sepc_detail_gov,omitempty"`
    IndustrySepcDetailAcc string           `json:"industry_sepc_detail_acc,omitempty"`
    ChargeAmount          string           `json:"charge_amount,omitempty"`
    ChargeFlags           string           `json:"charge_flags,omitempty"`
    SettlementId          string           `json:"settlement_id,omitempty"`
    TradeSettleInfo       *TradeSettleInfo `json:"trade_settle_info,omitempty"`
    AuthTradePayMode      string           `json:"auth_trade_pay_mode,omitempty"`
    BuyerUserType         string           `json:"buyer_user_type,omitempty"`
    MdiscountAmount       string           `json:"mdiscount_amount,omitempty"`
    DiscountAmount        string           `json:"discount_amount,omitempty"`
    Subject               string           `json:"subject,omitempty"`
    Body                  string           `json:"body,omitempty"`
    AlipaySubMerchantId   string           `json:"alipay_sub_merchant_id,omitempty"`
    ExtInfos              string           `json:"ext_infos,omitempty"`
    PassbackParams        string           `json:"passback_params,omitempty"`
    HbFqPayInfo           *HbFqPayInfo     `json:"hb_fq_pay_info,omitempty"`
    CreditPayMode         string           `json:"credit_pay_mode,omitempty"`
    CreditBizOrderId      string           `json:"credit_biz_order_id,omitempty"`
    HybAmount             string           `json:"hyb_amount,omitempty"`
    BkagentRespInfo       *BkAgentRespInfo `json:"bkagent_resp_info,omitempty"`
    ChargeInfoList        []*ChargeInfo    `json:"charge_info_list,omitempty"`
    BizSettleMode         string           `json:"biz_settle_mode,omitempty"`
}

// client.TradeQuery()
func TradeSQuery(out_trade_no string) (resp *alipay.TradeQueryResponse) {
    //aliPayPublicKey := "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA1wn1sU/8Q0rYLlZ6sq3enrPZw2ptp6FecHR2bBFLjJ+sKzepROd0bKddgj+Mr1ffr3Ej78mLdWV8IzLfpXUi945DkrQcOUWLY0MHhYVG2jSs/qzFfpzmtut2Cl2TozYpE84zom9ei06u2AXLMBkU6VpznZl+R4qIgnUfByt3Ix5b3h4Cl6gzXMAB1hJrrrCkq+WvWb3Fy0vmk/DUbJEz8i8mQPff2gsHBE1nMPvHVAMw1GMk9ImB4PxucVek4ZbUzVqxZXphaAgUXFK2FSFU+Q+q1SPvHbUsjtIyL+cLA6H/6ybFF9Ffp27Y14AHPw29+243/SpMisbGcj2KD+evBwIDAQAB"
    client, err := alipay.NewClient("2021003166661162", "MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDVMppatSws5rsHcMcGAapEZWgDq6bZOHqfXfribVr/KkdWcVY1egaivGiw6w0eyLXdvcWioin/3I0pHAyAWXm/3V11orTeE6OTPWNn1TrFyhU3+k+1q3KwqiR55qY5AKKqZxfM8cEOdeYzFXFbEXnuE3ngDem06DQAdLkaOeThnujDcKBUNQuCRrCm5SUy/KSv9mdIT6wdG5M4+VgW4pN6QHILmEfjVfsRRhiyZVHoqVrFI2xezXs+QRgb8HAFoeV4d4KnkSRwpWqHlsnFgNAm/Gm+O90bzSYtrP2Fn0AJdgqCMcE3X+GJfhOrW88TeZHCT5gDkMHwFZZuPpIhzbbBAgMBAAECggEAYUVGMW6Jqi7XVEy7MV1fHveZXltZs6/WGpIZXmeAZf0XMTRineF/143YwidsBAiVGYd0/X+5Y9hvuzrl5UVtjAFmm75RgSU6s3oFuaEKwKUzyyd0aLHBkSL/o3J9knJcxRxmOoZui7d1AQeegWtW0y2lpHkUkQqEd1TKA1I4wEsYPPwoaPYXIE8On54peo9f9EVAHc1eJxSW+yWE4UQI4DnDwjBz6Uux4re6cTBTvNYCxwsMReIFyuN3kHFRsulWVlTwchVWUbt1Boc0HoUICyfby3QkyaJGWTXvhnfiPxpkv511MtBmG5+yizTV2UEKs4N+oQ4bz7cbjFdzUe7sRQKBgQDr95RNUrGh4aQq+gp2FjtBJsl6vgo+lcCszmlMNpGYtctH0T0WSD9m7oyyZP1ZN3mRh1o4v376UwZ622cRr/95TCQgenJt7Ina84vV9+ATgWgiQfUdHV8D5CGnydHPwYil5OKA0ucB+FKDBcgQGOphUTG1/y6AvSHr4+NsMhv6MwKBgQDnTClDiuwpsdeSAbW/nxb3TaUgFseiAeY2s1T2cJsAd5p3yaq0UVVgV9OUX3tLq21HWmIVTDj09G4xpp4NVNRRPcE2o/CFDh9SS2myqRXy4gH6saQp01bWh7uzmXD8BfafqebjK+0j7C53ScbMXBgHidUVHdTCLBpuR7iikpW/OwKBgQCAUZ9lOR398U2sTVMZClfowyX3yJabmCYyEwFx/47Ho7zK7j8w+dL4r6r1bDPVq3RBroBischkanfwoZV4KeRc2woeW1gU7Pe+iIi3r9c75DhzwLiBv7Im1I10yCx/tTgRNtnxwj77dEWymJdGIbZ7e4Lz/LQWMEPdGo1XDhzmvwKBgQC65yJZAAOCVcFarKMPKyFFyaprWb0LvvkmrpczZR77q6pYrc+RUj/pUE8akGVzah0uEW08xJEp7/KzkG4bW7cNxxdAbg1Hl3fb6jCJPHUOBW+QAsgjPDHpvVkB8jYIkVEPCB4Y6EACTTHnFujb7ndEcC6Nl8N6/GSHRNGAHW+ATQKBgDd47gCdbQ14uz/dk4LWMNckzBTDbeGz9rofR9+hOrT3be37dUoFeWgnsyWW6tErv40ViEDKs8Hh1gPNJec/MyA/hAAPs1pT8aw5AOCq1Um7r7nu4lf+J2m+25wPhC2TQ4BAACH8t2K00V0oMv4XmftZTnAeMvFCQhUS/6/FyxnF", true)
    if err != nil {
        xlog.Error(err)
    }
    //配置公共参数
    client.SetCharset("utf-8").
        SetSignType(alipay.RSA2)
    //SetAppAuthToken("201908BB03f542de8ecc42b985900f5080407abc")

    //请求参数
    bm := make(gopay.BodyMap)
    bm.Set("out_trade_no", out_trade_no)

    //查询订单
    aliRsp, errs := client.TradeQuery(context.Background(), bm)
    if err != nil {
        xlog.Error("err:", errs)
    }
    return aliRsp
}

func test_alipay() {
    client, err := alipay.NewClient("2021004130632232", "MIIEvwIBADANBgkqhkiG9w0BAQEFAASCBKkwggSlAgEAAoIBAQCxo+vNCXQi9B032U2OsrHAeY5DO+mCMUiNpvBa4Id8IIj3KP/fQb9+SwU+U0EvpDN/d2wf9ArvK4MUDWhcYZn50/eJ2hCqqoXJ+8KP6gAYNGzPHfb42pPRqW23Z/0QpIVhX30EQBTm326HWwCb9jsq6aIRHzAF9apQavTyzyqSVaAbFk3nHARBG0aTHo0P9Zp6vz1nbAgyPwQykuGLdYmtQDEHu53IA73aWHoV3Fqznftq/n6mo/PaJBe21LWm1jqe+EW8vgmLhdq+J1uymdmIxmCzqygxbP9q6ET608l41XtPUxtJX4EkfNucebJPT80t3nP/I3TYu9z/bSoRsqjNAgMBAAECggEBAIqG/SSIwcNcQMjDseKc4VbqtBEkUIWRrzWfwIHt5Fnc+VJc66SLnix7jlw7CnN/hhVZ6LzGUByQ/wgNKJwaFLwpGzmqDyM4FVsc9G3MGkTF5TBi+qy3r1xp1sCW3Fc1JTp4/4HoAyTDimsNgV0eWKevSA44FvgeyrEpp1kOSvGPfNhYVdT5Ak+zICApCwP1/CrTxQjSKPP/HHwQRbrhVVvu7jh1xzoTjIfecATi9D20txBWYm+fYEGBjR6T8FNHP5pjTr52ty84wz07uhTXqOX8gTbSDb/BJLa3/+aR7X4oNvsSa9iriqcs1xMq4KkB9bdZ8EzQ53HM3RHnJ20eP6ECgYEA5i/pIEdT4/vQr3KKyuAJ7cYWuD1fWkegA2FqyeLk1MmBCBNjnTCJWMicobKE+ihK2wNUDeuOHPG3sVuw0VzPjQ+lBUYovW9V1UFX4E2o/akNXyBKQJG3sbJWRee6s3aCqZVqbhnLO7t1DxDAId37H7xYDN/qObnOoXAG3reynhkCgYEAxY+GK0IgntEl8dSg3dE+XoO9X62fXze+8Z0QLhVbiKhwAqxIamJ1OjWWzGy5NtXOItvwkumw3ydJ4QXvZZpb8DzgOir9E9u18oh7WkVzL9kkynqfSk3XDXR2iI0s86VvrlSaValSyRsciRF1QEHznNJbjT4QRS45EjDxH/dlztUCgYEA0CMEPh6g2WXp5arBmv4HnEtgYcmEvcJECqp8f/48gbeOh7nYedrYZkJHduJP4U6rmOuihk+3Ga7rNWC+OiEcvuUlhuZQkjHov8Ks7fHq2yqQH7K30Tixi+jAn8cQB5QiQ6sKKHIEVYeEKlIwGK96kdChIUsapIXBNDJy09HwnYkCgYBzKyke/KzBiNFi+f5RcVK3jHsQVMnMm1XPyi0NgFvc/bxWgpKwmfcW2Piw8UzDv74sqiTDsEHwxRmXeXtGssaX9RUOM9NXCUU3PwMR69yrbx24f+VuTpRofpU/I3WqD65cZWuXNl9RZ2GqMig1Ln1S1XqTizO28KxKg4d9iB6shQKBgQDFfupwsResXN/v3bTZEMGJJqegRsnn0Tiw6XDbbSn7VgIIf4AXhyONpCgZ92F3z9f+wjbBPqdNP2pLVz1ZZ0Qnhjv2UkeWZbXl8VRYa+LGKwo7yTmoV8WcD6SOtIYwqmZyuwVTgXlaItYhSAeczx4+wwhqNwnQnQWbVKIVkEjR/A==", true)
    if err != nil {
        xlog.Error(err)
        return
    }

    client.DebugSwitch = gopay.DebugOn

    client.SetLocation(alipay.LocationShanghai).
        SetCharset(alipay.UTF8).
        SetSignType(alipay.RSA2).
        SetReturnUrl("https://www.fmm.ink").
        SetNotifyUrl("https://www.fmm.ink")

    bm := make(gopay.BodyMap)
    bm.Set("subject", "电脑网站支付").
        Set("out_trade_no", "GZ201909081743431444").
        Set("total_amount", "0.01").
        Set("timeout_express", "2m")

    aliRsp, err := client.TradePagePay(context.Background(), bm)
    fmt.Println(aliRsp)
    if err != nil {
        if bizErr, ok := alipay.IsBizError(err); ok {
            xlog.Errorf("%+v", bizErr)
            // do something
            return
        }
        xlog.Errorf("client.TradePay(%+v),err:%+v", bm, err)
        return
    }
}

func main() {
    test_alipay()
}
