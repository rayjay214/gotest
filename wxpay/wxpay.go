package main

import (
    "context"
    "encoding/json"
    "fmt"
    "github.com/wechatpay-apiv3/wechatpay-go/core"
    "github.com/wechatpay-apiv3/wechatpay-go/core/option"
    "github.com/wechatpay-apiv3/wechatpay-go/services/payments/h5"
    "github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
    "github.com/wechatpay-apiv3/wechatpay-go/services/payments/native"
    "github.com/wechatpay-apiv3/wechatpay-go/utils"
    "io/ioutil"
    "log"
    "net/http"
    "time"
)

func nativepay() {
    var (
        mchID                      string = "1665176707"                               // 商户号
        mchCertificateSerialNumber string = "1CDD40BCFF230F0724B9F22AF05FCABF9AFA6D8F" // 商户证书序列号
        mchAPIv3Key                string = "9336ebf25087d91c818ee6e9ec29f8c1"         // 商户APIv3密钥
    )
    // 使用 utils 提供的函数从本地文件中加载商户私钥，商户私钥会用来生成请求的签名
    mchPrivateKey, err := utils.LoadPrivateKeyWithPath("apiclient_key.pem")
    if err != nil {
        log.Fatal("load merchant private key error")
    }
    ctx := context.Background()
    // 使用商户私钥等初始化 client，并使它具有自动定时获取微信支付平台证书的能力
    opts := []core.ClientOption{
        option.WithWechatPayAutoAuthCipher(mchID, mchCertificateSerialNumber, mchPrivateKey, mchAPIv3Key),
    }
    client, err := core.NewClient(ctx, opts...)
    if err != nil {
        log.Fatalf("new wechat pay client err:%s", err)
    }
    // 以 Native 支付为例
    svc := native.NativeApiService{Client: client}
    // 发送请求
    resp, result, err := svc.Prepay(ctx,
        native.PrepayRequest{
            Appid:       core.String("wx1dc66fa3e240e8cb"),
            Mchid:       core.String("1665176707"),
            Description: core.String("Image形象店-深圳腾大-QQ公仔"),
            OutTradeNo:  core.String("1217752501201407033233368018"),
            Attach:      core.String("自定义数据说明"),
            NotifyUrl:   core.String("https://www.weixin.qq.com/wxpay/pay.php"),
            Amount: &native.Amount{
                Total: core.Int64(100),
            },
        },
    )
    // 使用微信扫描 resp.code_url 对应的二维码，即可体验Native支付
    log.Printf("status=%d resp=%s", result.Response.StatusCode, resp)
}

func ExampleJsapiApiService_Prepay() {
    var (
        mchID                      string = "1665176707"                               // 商户号
        mchCertificateSerialNumber string = "1CDD40BCFF230F0724B9F22AF05FCABF9AFA6D8F" // 商户证书序列号
        mchAPIv3Key                string = "9336ebf25087d91c818ee6e9ec29f8c1"         // 商户APIv3密钥
    )

    // 使用 utils 提供的函数从本地文件中加载商户私钥，商户私钥会用来生成请求的签名
    mchPrivateKey, err := utils.LoadPrivateKeyWithPath("apiclient_key.pem")
    if err != nil {
        log.Print("load merchant private key error")
    }

    ctx := context.Background()
    // 使用商户私钥等初始化 client，并使它具有自动定时获取微信支付平台证书的能力
    opts := []core.ClientOption{
        option.WithWechatPayAutoAuthCipher(mchID, mchCertificateSerialNumber, mchPrivateKey, mchAPIv3Key),
    }
    client, err := core.NewClient(ctx, opts...)
    if err != nil {
        log.Printf("new wechat pay client err:%s", err)
    }

    svc := jsapi.JsapiApiService{Client: client}
    resp, result, err := svc.PrepayWithRequestPayment(ctx,
        jsapi.PrepayRequest{
            Appid:         core.String("wx1dc66fa3e240e8cb"),
            Mchid:         core.String("1665176707"),
            Description:   core.String("Image形象店-深圳腾大-QQ公仔"),
            OutTradeNo:    core.String("1217752501201407033233368018aaa"),
            TimeExpire:    core.Time(time.Now()),
            Attach:        core.String("自定义数据说明"),
            NotifyUrl:     core.String("https://www.weixin.qq.com/wxpay/pay.php"),
            GoodsTag:      core.String("WXG"),
            LimitPay:      []string{"LimitPay_example"},
            SupportFapiao: core.Bool(false),
            Amount: &jsapi.Amount{
                Currency: core.String("CNY"),
                Total:    core.Int64(100),
            },
            Payer: &jsapi.Payer{
                Openid: core.String("oYPaz66mYti0IfNLlgmvq5uIva0g"),
            },
            Detail: &jsapi.Detail{
                CostPrice: core.Int64(608800),
                GoodsDetail: []jsapi.GoodsDetail{jsapi.GoodsDetail{
                    GoodsName:        core.String("iPhoneX 256G"),
                    MerchantGoodsId:  core.String("ABC"),
                    Quantity:         core.Int64(1),
                    UnitPrice:        core.Int64(828800),
                    WechatpayGoodsId: core.String("1001"),
                }},
                InvoiceId: core.String("wx123"),
            },
            SceneInfo: &jsapi.SceneInfo{
                DeviceId:      core.String("013467007045764"),
                PayerClientIp: core.String("14.23.150.211"),
                StoreInfo: &jsapi.StoreInfo{
                    Address:  core.String("广东省深圳市南山区科技中一道10000号"),
                    AreaCode: core.String("440305"),
                    Id:       core.String("0001"),
                    Name:     core.String("腾讯大厦分店"),
                },
            },
            SettleInfo: &jsapi.SettleInfo{
                ProfitSharing: core.Bool(false),
            },
        },
    )

    if err != nil {
        // 处理错误
        log.Printf("call Prepay err:%s", err)
    } else {
        // 处理返回结果
        log.Printf("status=%d resp=%s", result.Response.StatusCode, resp)
    }
}

func ExampleH5ApiService_Prepay() {
    var (
        mchID                      string = "1665176707"                               // 商户号
        mchCertificateSerialNumber string = "1CDD40BCFF230F0724B9F22AF05FCABF9AFA6D8F" // 商户证书序列号
        mchAPIv3Key                string = "9336ebf25087d91c818ee6e9ec29f8c1"         // 商户APIv3密钥
    )

    // 使用 utils 提供的函数从本地文件中加载商户私钥，商户私钥会用来生成请求的签名
    mchPrivateKey, err := utils.LoadPrivateKeyWithPath("apiclient_key.pem")
    if err != nil {
        log.Print("load merchant private key error")
    }

    ctx := context.Background()
    // 使用商户私钥等初始化 client，并使它具有自动定时获取微信支付平台证书的能力
    opts := []core.ClientOption{
        option.WithWechatPayAutoAuthCipher(mchID, mchCertificateSerialNumber, mchPrivateKey, mchAPIv3Key),
    }
    client, err := core.NewClient(ctx, opts...)
    if err != nil {
        log.Printf("new wechat pay client err:%s", err)
    }

    svc := h5.H5ApiService{Client: client}
    resp, result, err := svc.Prepay(ctx,
        h5.PrepayRequest{
            Appid:         core.String("wx1dc66fa3e240e8cb"),
            Mchid:         core.String("1665176707"),
            Description:   core.String("Image形象店-深圳腾大-QQ公仔"),
            OutTradeNo:    core.String("1217752501201407033233368018"),
            TimeExpire:    core.Time(time.Now()),
            Attach:        core.String("自定义数据说明"),
            NotifyUrl:     core.String("https://www.weixin.qq.com/wxpay/pay.php"),
            GoodsTag:      core.String("WXG"),
            LimitPay:      []string{"LimitPay_example"},
            SupportFapiao: core.Bool(false),
            Amount: &h5.Amount{
                Currency: core.String("CNY"),
                Total:    core.Int64(100),
            },
            Detail: &h5.Detail{
                CostPrice: core.Int64(608800),
                GoodsDetail: []h5.GoodsDetail{h5.GoodsDetail{
                    GoodsName:        core.String("iPhoneX 256G"),
                    MerchantGoodsId:  core.String("ABC"),
                    Quantity:         core.Int64(1),
                    UnitPrice:        core.Int64(828800),
                    WechatpayGoodsId: core.String("1001"),
                }},
                InvoiceId: core.String("wx123"),
            },
            SceneInfo: &h5.SceneInfo{
                DeviceId: core.String("013467007045764"),
                H5Info: &h5.H5Info{
                    AppName:     core.String("王者荣耀"),
                    AppUrl:      core.String("https://pay.qq.com"),
                    BundleId:    core.String("com.tencent.wzryiOS"),
                    PackageName: core.String("com.tencent.tmgp.sgame"),
                    Type:        core.String("iOS"),
                },
                PayerClientIp: core.String("14.23.150.211"),
                StoreInfo: &h5.StoreInfo{
                    Address:  core.String("广东省深圳市南山区科技中一道10000号"),
                    AreaCode: core.String("440305"),
                    Id:       core.String("0001"),
                    Name:     core.String("腾讯大厦分店"),
                },
            },
            SettleInfo: &h5.SettleInfo{
                ProfitSharing: core.Bool(false),
            },
        },
    )

    if err != nil {
        // 处理错误
        log.Printf("call Prepay err:%s", err)
    } else {
        // 处理返回结果
        log.Printf("status=%d resp=%s", result.Response.StatusCode, resp)
    }
}

type OpenIdResponse struct {
    AccessToken string `json:"access_token"`
    ExpiresIn   int    `json:"expires_in"`
    OpenID      string `json:"openid"`
    Scope       string `json:"scope"`
}

func getOpenID() (string, error) {
    code := "0a3cjQ000dwonR1Nio200KbDmK2cjQ0c"
    appID := "wx1dc66fa3e240e8cb"
    appSecret := "d5bde1c75d156b28bae56a4aeed1c55d"

    url := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code", appID, appSecret, code)

    resp, err := http.Get(url)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    // 读取响应体
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }
    fmt.Println(string(body))

    // 解析 JSON 响应
    var openIdResponse OpenIdResponse
    err = json.Unmarshal(body, &openIdResponse)
    if err != nil {
        return "", err
    }

    fmt.Println(openIdResponse.OpenID)
    return openIdResponse.OpenID, nil
}

type AccessTokenResponse struct {
    AccessToken string `json:"access_token"`
    ExpiresIn   int    `json:"expires_in"`
}

func getAccessToken() (string, error) {
    appID := "wx1dc66fa3e240e8cb"
    appSecret := "d5bde1c75d156b28bae56a4aeed1c55d"

    url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s", appID, appSecret)

    resp, err := http.Get(url)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }

    // 解析 JSON 响应
    var accessTokenResponse AccessTokenResponse
    err = json.Unmarshal(body, &accessTokenResponse)
    if err != nil {
        return "", err
    }

    return accessTokenResponse.AccessToken, nil
}

func main() {
    ExampleJsapiApiService_Prepay()
    //ExampleH5ApiService_Prepay()

    //accessToken, _ := getAccessToken()
    //getOpenID()
}
