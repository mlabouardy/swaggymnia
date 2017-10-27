[![CircleCI](https://circleci.com/gh/mlabouardy/swaggymnia/tree/master.svg?style=svg&circle-token=bcfce92d1e46aaf0d50b4b3fa8baf8406d4bc115)](https://circleci.com/gh/mlabouardy/swaggymnia/tree/master) [![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)

Convert Insomnia REST API to Swagger Docs

# To Do

- Support JSON Body

# Docker image

You can use swaggymnia through official Docker image (automated build on DockerHub from this repository)

```
docker run --rm -v $PWD:/data mlabouardy/swaggymnia [options and arguments...]
```

You need to mount some local directory at /data so input and output files can be accessed by the container

## Maintainers

- Mohamed Labouardy
