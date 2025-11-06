[![CircleCI](https://circleci.com/gh/mlabouardy/nexus-cli.svg?style=svg)](https://circleci.com/gh/mlabouardy/nexus-cli) [![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)

<div align="center">
<img src="logo.png" width="60%"/>
</div>

Nexus CLI for Docker Registry

## Usage

<div align="center">
<img src="example.png"/>
</div>

## Download

Below are the available downloads for the latest version of Nexus CLI (1.0.1-beta). Please download the proper package for your operating system and architecture.


## Available Commands

```
$ nexus-cli configure
```

```
$ nexus-cli image ls
```

```
$ nexus-cli image tags -name mlabouardy/nginx
```

```
$ nexus-cli image info -name mlabouardy/nginx -tag 1.2.0
```

```
$ nexus-cli image delete -name mlabouardy/nginx -tag 1.2.0
```

```
$ nexus-cli image delete -name mlabouardy/nginx -keep 4
```

```
$ nexus-cli image size -name mlabouardy/nginx
```

## New Available Commands

```
# 按时间查看标签（从旧到新）
nexus-cli image tags --name myimage --sort time

# 按时间查看标签（从新到旧）
nexus-cli image tags --name myimage --sort time-desc

# 按时间排序删除，保留最新的5个
nexus-cli image delete --name myimage --keep 5 --sort-by time

# 删除最新的几个（保留旧的）
nexus-cli image delete --name myimage --keep 5 --sort-by time-desc
```

## Tutorials

* [Cleanup old Docker images from Nexus Repository](http://www.blog.labouardy.com/cleanup-old-docker-images-from-nexus-repository/)
