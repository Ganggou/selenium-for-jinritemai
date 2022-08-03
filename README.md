# selenium-for-jinritemai

### 登录
* 执行命令 go run main.go -a login
* 登录成功后cookie将保存至cookies文件夹
* 退出

### 获取数据
* 执行命令 go run main.go
* 成功登录后设置筛选条件，点查询
* 在程序界面敲入任意命令，程序将逐页点开隐藏数据并采集，有时会因频率过高限制隐藏数据采集
* 确认扫描完成后，输入export命令导出至excel
* 退出
