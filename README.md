# stress-test-go
压测工具（go语言版）

使用方法：
- 1.先在custom目录下实现自己的压测任务逻辑，可参考GooleTask，实现CustomTaskRunner接口
- 2.在main中配置压测参数并运行
- 3.输出报告在reports/stats目录下，如果要查看请求详情，可查看reports/records目录