package main

/*
	静态构建是指构建后的应用运行所需的所有符号、指令和数据都包含在自身的
	二进制文件当中, 没有任何对外部动态共享库的依赖; 静态构建出的二进制文件
	由于包含所有符号、指令和数据, 因而通常要比非静态构建的应用大许多;
	默认情况下, Go没有采用静态构建; ./server.go
	<go语言精进之路2>, 60.6 节, 此节内容基于go1.14, 后续版本相关内容有所变化
*/
