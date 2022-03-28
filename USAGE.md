## Using Makeup

To use Makeup, you need two things, a `main.mk` in the root of your project, and a `.mk` in the subdirectory for each 'component' of your application. A component is anything you deem needing a build, test, run, clean loop in your workflow.

To start, create `main.mk`:
```makefile
# check go version
# equal 1.18

include ./testapp/testapp.mk
include ./testapp2/testapp2.mk
```

The first two lines are 'checks', they ensure that the correct versions of tools are being used. This can ensure that you and your team are all using the same version of a dependency, for example. You can have as many checks as you want, but they must be sequential lines that start with `# check` and `# equal`. 

> The equality is actually a 'contains' check, so even though `go version` outputs something like `go version go1.18 darwin/arm64`, since it contains `1.18`, it passes the check.

The `include` statements are how you add components to the project. Each `.mk` file included in `main.mk` must have the following targets (even if they're empty): `build`, `run`, `test`, `clean`, `env`. For example:
```makefile
build:
	go build -o ${BIN_DEST}

run:
	${BIN_DEST}

test:
	go test -v ./...

env:
	echo "SOMETHING_IMPORTANT=important value"

clean:
	rm ${BIN_DEST}
```

As you can see, Makeup uses these standardized targets to control the lifecycle of your environment.

You can run `makeup` to build each component and start them all, together. 

Other commands include `makeup test` and `makeup clean` which run the `test` and `clean` targets on each of your components, sequentially.

## Generate Makefile
Coming soon is the ability to run `makeup generate`. This will generate a `Makefile` that will simulate the workflow of makeup so that anyone can take advantage of these abilities, even if they don't have makeup installed. They'll just be able to run `make up`😉