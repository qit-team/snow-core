## v0.1.17(2019-11-18)

### New Features
- 去掉测试用fmt


## v0.1.16(2019-11-18)

### New Features
- 队列queue新增aliyunmq【rocketmq】驱动

## v0.1.15(2019-11-10)

### New Features
- 队列queue出队增加返回参数：消费次数

## v0.1.14(2019-11-10)

### New Features
- 更改mns队列默认可见时间

## v0.1.13(2019-10-20)

### New Features
- 缓存驱动增加，内存缓存

## v0.1.12(2019-09-29)

### New Features
- 缓存组件增加decr和incr方法
- 增加单测文件

## v0.1.11(2019-08-25)

### New Features
- 日志增加traceId
- middleware组件和httputil组件增加内容串联traceId
- ctxkit包优化，方便扩展


## v0.1.10(2019-08-02)

### New Features
- 补充单测案例

### Bug Fix
- cache和queue包在获取对象时读锁枷锁未配对解锁

## v0.1.9(2019-08-01)

### New Features
- 补充单测案例

### Changes
- 优化utils包HttpBuildQuery的map嵌套转换实现

## v0.1.8(2019-07-26)

### New Features
- rediscache的单元测试案例

### Changes
- rediscache的Get返回优化。若key不存在之前是返回错误类型ErrNil,现在不返回错误，返回字符串为空

### Bug Fix
- 修复rediscache的SetMulti实现bug

## v0.1.7(2019-07-25)

### Changes
- 更新qit-team/work包的版本号v0.3.3->v.0.3.4

## v0.1.6(2019-07-24)

### Bug Fix
- 修复utils包HttpBuildQuery的对值非字符串的处理bug

## v0.1.5(2019-07-23)

### New Features
- Command执行脚本模式支持

## v0.1.4(2019-07-23)

### Changes
- utils工具包
    - HTTP请求工具包封装建议的Get Post PostJson Request方法

## v0.1.3(2019-07-22)

### New Features
- Redis组件服务
- Log组件服务
- DB组件服务
- Config通用配置结构
- Cache缓存及驱动
- Queue队列及驱动
- Http的通用中间件和通用上下文kit
- Kernel内核包
    - close服务注册
    - provider组件注册
    - container容器注入
    - server通用服务启动
- utils工具包
    - HTTP请求工具包
    - 其他常用函数工具包