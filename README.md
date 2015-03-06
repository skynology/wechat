# Wechat SDK for Skynology Cloud Code
[![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg)](LICENSE.md)

其于 [@产先生](https://github.com/chanxuehong) 的golang版微信SDK开发, 少量增删一些功能. 便于在上空云的**云代码**环境中使用. 若您在想在独立项目中使用微信SDK, 请直接使用[产先生的微信SDK](https://github.com/chanxuehong/wechat) .

## 安装
无需安装, 已经直接集成到上空云的云代码环境中.

## 使用说明
具体使用说明, 请参考上空云文档中心的[云代码之微信SDK使用](http://developer.skynology.com/weixin-sdk.html)

## 主要修改部分.
* 删除了各种Handler
* 删除了多媒体文件上传/下载功能.(提供下载URL)
* 删除了二级package, 都同意放到跟目录下(便于云代码中调用)
* 所以和微信交互功能全部是主Client的Method.(云代码会提供主Client, 无需自己实例化)



## 代码许可
原代码中的已经有的版权/许可请参考[产先生的授权协议](https://github.com/chanxuehong/wechat/blob/master/LICENSE).   
其他增删部分以MIT授权方式开放.   
This SDK is released under the MIT license. See the [LICENSE](https://github.com/skynology/objc-sdk/blob/master/LICENSE) file for more details.

