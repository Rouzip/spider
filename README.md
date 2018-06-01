# spider
分布式爬虫，分为request，parse and store部分

## 分工

request：Rouzip

store：Equim

parse：Ein

## request

向store请求url

负责并发进行请求

将网页发送到parse

## store

暂存html作为队列交给parse进行解析

最终结果的存储

## parse

向store请求html

解析html

横向扩展

向store发送url  

### 通信标准
成功的操作——200 或 204
暂时没有资源，需要等待——404 (这个只有获取 HTML 的一方(Ein)会收到)
没有资源，今后也不会有，可以停止了——410
你的请求有错误——400
我这处理有错误——500

### URL
base URL 为 https://api.ekyu.moe/d-spider/v1  
GET /URL  
POST /URL  
GET /HTML  
POST /HTML  
