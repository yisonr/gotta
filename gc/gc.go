package main

/*
	垃圾回收
	TODO: 有余力的情况下研读<垃圾回收算法手册>
	在计算机科学中,垃圾回收(Garbage Collection, GC) 是自动内存管理的一种形式,
	通常由垃圾收集器收集并适时回收或重用不再被对象占用的内存;
	垃圾回收作为内存管理的一部分, 包含3个重要的功能:
	- 分配和管理新对象
	- 识别正在使用的对象
	- 清除不再使用的对象

	垃圾回收会让开发变得更加简单, 屏蔽了复杂且容易出错的操作


	垃圾回收的好处:
	- 减少错误和复杂性
		没有垃圾回收就必须手动分配内存, 释放内存; 不管是内存泄露还是野指针都
		会增加开发程序的难度, 虽然垃圾回收不保证完全不产生内存泄露, 但其提供了
		重要的保障, 即不再被引用的对象最终被收集, 这种设定同样避免了悬空
		指针(TODO), 多次释放等手动管理内存时会出现的问题, 具有垃圾回收功能的
		语言屏蔽了内存管理的复杂性, 开发者可以更好的关注核心的业务逻辑;
	- 解耦
		当两个模块同时维护一个内存时, 释放内存必须特别小心, 手动分配的问题在
		于难以在本地模块内做出全局的决定, 而具有垃圾回收功能的语言将垃圾收集
		的工作托管给了具有全局视野的运行时代码, 开发者编写的业务模块将真正实
		现解耦, 从而利于开发和调试;

	因为垃圾回收带来额外的成本, 需要保存内存的状态信息(例如是否使用, 是否包含
	指针)并扫描内存, 很多时候还需要中断整个程序来处理垃圾回收; 因此垃圾回收
	对于要求极致的速度和内存要求极小的场景并不适用, 却是开发大规模, 分布式, 微
	服务集群的极佳选择;

	内存管理和垃圾回收都属于 go 语言最复杂的模块, 没有最完美的垃圾回收算法,
	因为每个应用程序的硬件条件, 工作负载和性能要求都是不同的, 理论上, 可以为
	单独的应用程序设计最佳的内存分配方案;
	通用的具有垃圾回收的编程语言会提供通用的垃圾会收拾算法, 并且每一种语言侧重
	的垃圾回收目标不尽相同; 垃圾回收的常见指标包括程序暂停时间, 空间开销, 回收
	的及时性等, 更具设计目标的侧重点不同有不同的垃圾回收算法;


	以下介绍垃圾回收的5种经典算法:
	TODO: 了解熟悉经典的垃圾回收算法, 精通golang的三色标记算法

	标记-清扫(Mark-Sweep):
		主要分两个阶段, 第一阶段扫描并标记当前活着的对象, 第2阶段是清扫没有被标记
		的垃圾对象; 算是一种间接的垃圾回收算法;

		扫描一般从栈上的根对象开始, 只要对象引用了其他对象, 就会一直向下扫描,
		因此可以采取深度优先搜索或广度优先搜索的方式进行扫描;

		在扫描阶段, 为了管理扫描对象的状态, 引入了经典的三色标记算法(TODO);

		标记-清扫算法的缺点在于可能产生内存碎片或空洞, 这会导致新对象分配失败;

	标记-压缩(Mark-Compact):
		通过将分散的, 活着的对象移动到更紧密的空间来解决内存碎片问题; 分为
		标记和压缩阶段:
		- 标记过程与标记-清扫中的标记过程类似;
		- 压缩阶段需要扫描活着的对象并将其压缩到空闲的区域, 这可以保证压缩后的
			空间更紧凑, 从而解决内存碎片问题; 同时, 压缩后的空间能以更快的
			速度查找到空闲的内存区域(在已经使用内存的后方);

		标记-压缩算法的缺点在于内存对象在内存的位置是随机的, 这通常会破坏缓存
		的局部性, 并且时常需要一些额外的空间来标记当前对象已经移动到了其他地方;
		在压缩阶段, 如果B对象发生了转移, 那么必须更新所有引用了B对象的对象的
		指针, 这无疑增加了实现的复杂性;

	半空间复制(Semispace Copy):
		一种空间换时间的算法; 经典的半空间复制算法只能使用一半的内存空间, 保留
		另一半的内存空间用于快速压缩内存;
		半空间复制的压缩性消除了内存碎片问题, 同时, 其压缩时间比标记-压缩算法
		更短, 半空间复制不分阶段, 在扫描对象时就可以直接压缩, 每个扫描到的对象
		都会从 fromspace 空间复制到 tospace 的空间(TODO); 因此, 一旦扫描完成就
		得到了一个压缩后的副本;

	引用计数(Reference Counting):
		每个对象都包含一个引用计数, 每当其他对象引用了此对象时, 引用计数就会
		增加; 反之取消引用后, 引用计数就会减少; 一旦引用计数为0, 就表明该对象
		为垃圾对象, 需要被回收; 引用计数算法简单高效, 在垃圾回收阶段不需要额外
		占用大量内存, 即便垃圾回收系统的一部分出现异常, 也有一部分对象被正常
		回收(TODO);
		致命缺点: 一些没有破坏性的操作, 如只读操作, 循环迭代操作也需要引用计数,
		栈上的内存操作或寄存器操作更新引用计数是不合理的(TODO:?); 同时, 引用
		计数必须原子更新, 并发操作同一个对象会导致引用计数难以处理自引用(TODO)
		的对象;

	分代GC:
		分代GC指将对象按照存活时间进行划分, 前提是: 死去的对象一般都是新创建
		不久的, 因此没有必要反复的扫描旧对象, 这样大概率会加快垃圾回收的速度,
		提高处理能力和吞吐量, 减少程序暂停的时间; 但是分代GC没有办法及时回收
		老一代的对象, 并且需要额外开销引用和区分新老对象, 特别是在有多代对象
		的时候(TODO);


	Golang 采用了并发三色标记算法进行垃圾回收, 三色标记是最简单的垃圾回收算法,
	其实现也相对简单; 引用计数由于其固有的缺陷在并发时很少使用, 不适合go这样的
	高并发语言, 同时应该想到go最重要的就是高并发能力, 所以语言上的很多设计都
	要更好的适配GMP;

	go 不选择压缩gc的原因: (TODO: 研究TCmalloc内存分配算法)
		压缩算法额主要优势是减少碎片且快速分配, go 使用了现代内存分配算法
		TCmalloc, 已经很好的解决了内存碎片问题, 并且由于需要加锁, 压缩算法
		并不适合在并发程序中使用, 且其实现比三色标记算法更加复杂;

	go 不选择分代gc的原因:(TODO:逃逸分析, 写屏障)
		分代gc的主要假设是大部分变成垃圾的对象都是新创建的, 但是由于编译器的
		优化, go语言通过内存逃逸的机制将会继续使用的对象转移到了堆中, 大部分
		生命周期很短的对象会在栈中分配, 这和其他使用分代GC的编程语言有显著的
		不用, 减弱了分代GC的优势; 同时分代GC需要额外的写屏障来保护并发垃圾
		回收时对像的隔代性, 会减慢GC的速度; 分代GC是被golang尝试过并抛弃的方案;


	由于垃圾回收具有减少复杂性及解耦的性质, 现代高级语言基本都具有垃圾回收
	这样的自动内存管理功能; go 语言的垃圾回收算法经历了复杂的演技过程, 从最初
	的单携程垃圾回收到最后实现了并发的三色抽象标记算法, 在演进的过程中大大减少
	了STW(TODO)的时间, 从最初的几百毫秒降到了现在的微妙级别, 对大部分场景来说
	几乎是无感知的;

	垃圾回收是go语言中最复杂的模块, 其贯穿于程序的整个生命周期, 涉及编译器,
	调度器, 内存分配, 栈扫描, 位图等技术; 掌握了垃圾回收就相当于掌握了go语言的
	运行时;
*/
