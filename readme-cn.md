# draw-graph

[English](readme.md) | 简体中文

## 项目说明

`draw-gragh` 是一个 Go 画图工具，只支持有向无环图(DAG, Directed Acyclic Graph)

## 运行

* cmd
    * main
      读取 `.json` 文件，每个json文件生成一个图片

        * `-d` 读取指定目录下所有 `json` 文件
        * `-i` 读取指定的 `json` 文件，多个文件用逗号分割
        * `-o` 指定输出的目录

## 数据格式

数据可以通过 `json` 文件提供，参考 [示例](cmd/data/example.json) 里的格式。大致如下：

## 致谢

gridder - github.com/shomali11/gridder
