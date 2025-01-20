# 获取Cookie
```javascript
var cookie=document.cookie;var ask=confirm('Cookie:'+cookie+'\n\nDo you want to copy the cookie to the clipboard?');if(ask==true){copy(cookie);msg=cookie}else{msg='Cancel'}
```


# 使用
两个子程序
1. 签到

编译
```shell
make sign
```
需要获取`https://api-takumi.miyoushe.com`和`https://api-takumi.mihoyo.com`下的Cookie，后续再优化逻辑
获取后填写到conf/config.yml中
```yaml
aliyun_auth:
    access_key: ""
    secret_key: ""
miyoushe_cookie: stuid=xxxx;stoken=xxx
mihoyo_cookie: account_id_v2=xxxx; cookie_token_v2=xxxx

```
然后运行./sign即可


2. 驱动盘分析

目前只做了截图提取驱动盘信息，代理人分析放在miyabi_test.go这种里面临时测试，后续再完善。

并且ocr使用的是阿里云的API，需要填写AK/SK到conf/config.yml
```yaml
aliyun_auth:
    access_key: "xxx"
    secret_key: "xxx"
```

```shell
# 编译
make parser

# 提取驱动盘截图里的信息
./parser ocr -f test/test2/Snipaste_2025-01-19_00-34-28.png
./parser ocr -f test/test2

# 如果存在某些识别结果解析失败，可手动调整然后解析
./parser ocr2 -c "河豚电音[5] 5 等级15/15 主属性 穿透率 24% 副属性 防御力 +2 45 穿透值 +2 27 攻击力 +1 38 异常精通 9"
```