# Zen Lang

The Zen language interpreter, built as an experiment while learning how to write my own custom programming language, and mostly just for fun. The goal is to eventually build and host a website running on a webserver that it built and hosted with Zen.

## Running Zen

You have two main options when it comes to running Zen, Docker and the binary.

### Binary

You can download the latest version of Zen from the [GitHub releases](https://github.com/Senither/zen-lang/releases), alternatively you can build your own.

### Docker

If you have Docker installed, getting started with Zen is easy. You can run the latest version of Zen, or specify a specific version you want to use, and then run the following command:

```
docker run --rm -it ghcr.io/senither/zen-lang <command>
```

However, if you want to run a file, you’ll need to mount a volume with the file in the container. Alternatively, you can create your own custom image with the Zen image as a base.

```Dockerfile
FROM ghcr.io/senither/zen-lang:latest AS builder

WORKDIR /src
COPY . .

RUN ["/zen-lang", "build", "main.zen", "-o", "program"]

FROM ghcr.io/senither/zen-lang:latest AS final

COPY --from=builder /src/program /program

ENTRYPOINT ["/zen-lang", "program"]
```

The Docker image produced will contain a single file with all the bytecode that's been compiled from the source, which will be run when the container starts.

## Building from source

> Note: As of right now the build process with the `make` command has only been verified to work on Windows.

The project uses `make` to simplify the build and testing process, so to get the project up and running you'll first need to install the dependencies with:

```
make install
```

Afterward you can build the project into a binary using:

```shell
# Builds a binary for your local OS
make build
# Builds a Docker image
make docker
```

## Testing

Zen comes with two primary ways of testing the language, there are the unit tests which uses the [go testing package](https://pkg.go.dev/testing) to ensure the internals of the language produces our expected output, and there is the language test which will run all the [Zen test files](tests), the language tests helps ensure code given to any of the runtimes will produce not only the expected result, but also the same result between them, allowing us to ensure that the evaluator and compiler are feature matching.

All the tests for Zen can be run using:

```
make test
```

## License

Zen Lang is open-sourced software licensed under the [MIT License](https://opensource.org/licenses/MIT).
