# 简单的 Golang Gradle 尝试

## 场景

设想一个场景，我们在编译 golang 项目之前，需要 codegen，我们可以怎么办。

在本仓库的例子中，我们假设有一个程序他依赖了 Antlr Codegen。

我们可以看到，在 Gradle 的构建系统中，我们引入了 Antlr Codegen 流程、Go Workspace 流程 和 Go Build 流程。
并且这三个流程的运行顺序和我们在行文中的顺序一致。

## 思路

1. 先使用 Antlr Codegen 将程序需要用到的额外需要的 `go` 代码生成到 `build` 文件夹中。
2. 再使用 Go Workspace 将生成后的代码和原先的项目信息关联在一起
3. 最后使用 Go Build 将项目编译

## 构建

使用指令
```bash
export GOROOT="您的 Golang 目录"
./gradlew build
```

编译后的文件可以在 build 文件夹中找到，其文件名为 `gradle-go`。

## 编译结果功能介绍

一个普通的四则运算计算器。输入算式输出结果。
