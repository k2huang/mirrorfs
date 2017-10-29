# mirrorfs

基于[https://github.com/bazil/fuse](https://github.com/bazil/fuse)和[https://github.com/elgutierrez/mirrorfs](https://github.com/elgutierrez/mirrorfs)<br/>
来学习用FUSE写用户态文件系统的示例程序。<br/>

注意：为了不依赖golang.org/x/net/context，我将自己的代码和bazil/fuse库中编译过程中会报错的地方都统一换成了标准库的context。

