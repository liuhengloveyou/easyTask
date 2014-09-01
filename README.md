easyTask
========

简单任务分发系统


接口:
-----

##### 1. 添加任务
    GET /putask?type=任务类型名&rid=记录ID&info=任务描述

##### 2. 获取任务
    GET /getask?type=任务类型名&name=工兵名字&num=取任务条数

##### 3. 更新任务
    GET /uptask?type=任务类型名&name=工兵名字&tid=任务ID&stat=任务状态(成功=1|失败=-1)&msg=错误信息

##### 4. 向服务器打招乎. 工兵启动要先向服务器打招乎才能取任务,相当于注册的动作
    GET /sayhi?type=任务类型名&name=工兵名字

##### 5. 添加新的任务类型. 会在数据库创建相应的表
    GET /newtype?name=任务类型名

##### 6. 心跳. 每个工兵每5秒向服务器发一个心跳包
    GET /beat?type=任务类型名&name=工兵名字
