easyTask
========

简单任务分发系统


接口:
-----

##### 1. 添加任务
    GET /putask?type=任务类型名&rid=记录ID&info=任务描述
    成功返回任务ID;出错返回错误信息串
    
##### 2. 获取任务
    GET /getask?type=任务类型名&name=工兵名字&num=取任务条数
    成功返回任务信息:[{"Tid":"任务ID", "Rid":"记录ID", "Info":"任务内容"},...];出错返回错误信息串

##### 3. 更新任务
    GET /uptask?type=任务类型名&name=工兵名字&tid=任务ID&stat=任务状态(成功=1|失败=-1)&msg=错误信息
    成功返回"OK";出错返回错误信息串
    
##### 4. 向服务器打招乎. 工兵启动要先向服务器打招乎才能取任务,相当于注册的动作
    GET /sayhi?type=任务类型名&name=工兵名字
    成功返回"OK";出错返回错误信息串
    
##### 5. 添加新的任务类型. 会在数据库创建相应的表
    GET /newtype?name=任务类型名
    成功返回"OK";出错返回错误信息串
    
##### 6. 心跳. 每个工兵每5秒向服务器发一个心跳包
    GET /beat?type=任务类型名&name=工兵名字
    成功返回"OK";出错返回错误信息串
    
数据库:
------
	CREATE DATABASE IF NOT EXISTS `taskManager` DEFAULT CHARACTER SET utf8;
	
	CREATE TABLE `tasks_demo` (
		`id` int(11) NOT NULL AUTO_INCREMENT,
		`tid` varchar(33) NOT NULL,
		`rid` varchar(32) NOT NULL,
		`info` varchar(1024) NOT NULL,
		`stat` int(11) NOT NULL DEFAULT '0', -- 1 = 新任务; 2 = 正在处理; 3 = 处理成功; -1 = 处理出错
		`addTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
		`overTime` timestamp NULL DEFAULT '0000-00-00 00:00:00',
		`rapper` varchar(256) DEFAULT NULL,
		`client` varchar(256) DEFAULT NULL,
		`remark` text,
		PRIMARY KEY (`id`),
		UNIQUE KEY `inx_tid` (`tid`)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8;




