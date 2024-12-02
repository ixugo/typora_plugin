
typora 在导出 HTML 文件时，其图片位置是根据写 markdown 时间定义的，比如在 markdown 中是网络图片，则导出 HTML 也是网络图片。

图片在编辑者本地，或图片在网络但 HTML 使用在离线环境，此时图片是无法正确加载的。

通过设置 typora 导出后执行命令，将图片替换成 base64 来解决以上问题。

[查看使用详情](https://blog.golang.space/p/typora-%E5%AF%BC%E5%87%BA%E5%B8%A6%E5%9B%BE%E7%89%87%E7%9A%84-html/)