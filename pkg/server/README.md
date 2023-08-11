# Bedrock / Backend

Hello! There is bedrock's backend source code. I think everything is blazing fast and alright!

Welcome to Pull Request and Issues! Any contribute is welcome!

Builtin Go with [Fiber](https://gofiber.io), [Fx](https://uber-go.github.io/fx/) and cola 🍻

## Usage

Easy to run, the frontend resources in embedded when building.

```shell
$ docker run code.smartsheep.studio/atom/bedrock \
    -p 9443:9443 \ 
    -v $(pwd)/config.toml:/app/config.toml
```