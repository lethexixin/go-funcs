#### 关于 confluent-kafka-go 静态编译
当你需要使用 confluent-kafka-go 组件时, 由于它依赖于 librdkafka, 如若要实现静态编译, 需要设置 CGO_ENABLED=1, 并且编译命令需要依赖 musl-gcc, 类似如下命令:
```shell
CC=/path/to/musl-gcc go build --ldflags '-linkmode external -extldflags "-static"' -tags musl
```

##### WSL2(Ubuntu) 安装 musl-gcc
```shell
sudo apt install musl-tools
```

##### [官网介绍](https://github.com/confluentinc/confluent-kafka-go): Static builds on Linux
----------------------

Since we are using `cgo`, Go builds a dynamically linked library even when using
the prebuilt, statically-compiled librdkafka as described in the **librdkafka**
chapter.

For `glibc` based systems, if the system where the client is being compiled is
different from the target system, especially when the target system is older,
there is a `glibc` version error when trying to run the compiled client.

Unfortunately, if we try building a statically linked binary, it doesn't solve the problem,
since there is no way to have truly static builds using `glibc`. This is
because there are some functions in `glibc`, like `getaddrinfo` which need the shared
version of the library even when the code is compiled statically.

One way around this is to either use a container/VM to build the binary, or install
an older version of `glibc` on the system where the client is being compiled.

The other way is using `musl` to create truly static builds for Linux. To do this,
[install it for your system](https://wiki.musl-libc.org/getting-started.html).

Static compilation command, meant to be used alongside the prebuilt librdkafka bundle:
```bash
CC=/path/to/musl-gcc go build --ldflags '-linkmode external -extldflags "-static"' -tags musl
```

##### 具体命令可参考
```shell
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 GO111MODULE=on CC=musl-gcc go build --ldflags '-linkmode external -extldflags "-static"' -tags musl
```
