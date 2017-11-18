# mirrorfs

基于[https://github.com/bazil/fuse](https://github.com/bazil/fuse)实现的一个 用户态文件系统 - mirrorfs(镜像文件系统)。 <br/>
执行 **./progname -mount path1 -mirror path2** 之后：<br>
就将我们的文件系统mirrorfs挂载到了path1上，之后对path1目录的操作实际上都是在操作path2目录，也就是将对path1的操作
"镜像"到了path2。<br>
示例讲解以及bazil/fuse的使用参看[这里](https://github.com/k2huang/blogpost/blob/master/golang/%E5%BA%94%E7%94%A8%E7%A8%8B%E5%BA%8F/%E7%94%A8Go%E5%92%8CFUSE%E8%87%AA%E5%B7%B1%E7%9A%84%E6%96%87%E4%BB%B6%E7%B3%BB%E7%BB%9F/README.md)。<br/>

注意：为了不依赖golang.org/x/net/context，我将自己的代码和bazil/fuse库中编译过程中会报错的地方都统一换成了标准库的context。

