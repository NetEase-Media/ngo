```seq
Title: ngo 接入oidc流程
web端->oidc: 唤起登录
oidc-->web端: 登录成功，重定向到web端，附带code
web端->ngo: 调用获取token接口
ngo->oidc: 根据code获取token
oidc-->ngo: 返回access_token，id_token
ngo->ngo: 验证id_token
ngo->oidc: 获取用户信息
oidc-->ngo: 返回用户信息
ngo-->web端: 生成ngo token，返回
web端->ngo: 正常业务接口
ngo-->web端: 正常业务返回
web端->ngo: 当token快超时，调用刷新token接口
ngo-->web端: 生成新的ngo token，返回
web端->ngo: 当长时间不操作，refresh token过期，调用业务接口
ngo-->web端: 返回token已过期，进入最开始重新登录流程
```