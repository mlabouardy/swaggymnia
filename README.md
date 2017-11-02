[![CircleCI](https://circleci.com/gh/mlabouardy/swaggymnia/tree/master.svg?style=svg&circle-token=bcfce92d1e46aaf0d50b4b3fa8baf8406d4bc115)](https://circleci.com/gh/mlabouardy/swaggymnia/tree/master) [![Go Report Card](https://goreportcard.com/badge/github.com/mlabouardy/swaggymnia)](https://goreportcard.com/report/github.com/mlabouardy/swaggymnia) [![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)

# Swaggymnia

<<<<<<< HEAD
Generate Swagger Documentation from Insomnia REST Client
=======
# Build

```
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .
```

# To Do
>>>>>>> 312d4e9649f6c252873b4af70ce166edb7a753d4

## Download

Below are the available downloads for the latest version of Swaggymnia  (1.0.0-beta). Please download the proper package for your operating system and architecture.

### Linux:

```
wget https://s3.eu-west-2.amazonaws.com/swaggymnia/1.0.0-beta/linux/swaggymnia
```

### Windows:

```
wget https://s3.eu-west-2.amazonaws.com/swaggymnia/1.0.0-beta/windows/swaggymnia
```

## How to use it

See usage with:

```
$ swaggymnia --help
```

Generate Swagger documentation:

```
$ swaggymnia generate -insomnia INSOMNIA_EXPORTED_FILE -config CONFIG_FILE -output FORMAT
```

| Option | Description |
| ------ | ----------- |
| -insomnia | Insomnia exported file |
| -config | API Global Configuration file (see [API Configuration Structure File](#License))|
| -output | Insomnia output format (json or yaml, default json)  |


## Example

Let's convert the following Insomnia API documentation to Swagger:

<div align="center">
  <img src="insomnia.png"/>
</div>

```
$ swaggymnia generate -i examples/watchnow.json -c examples/config.json -o json
```

<div align="center">
  <img src="swagger.png"/>
</div>

## API Configuration file

```
{
  "title" : "API Name",
  "version" : "API version",
  "host" : "API URL",
  "bastPath" : "Base URL",
  "schemes" : "HTTP protocol",
  "description" : "API description"
}
```

## Maintainers

- Mohamed Labouardy - mohamed@labouardy.com

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details
