# 💅🏽 Makeup 💄

> A local development tool to replace Docker Compose, based on Make.

Makeup uses simple Makefiles to create a faster developer workflow compared to Docker and docker-compose. It uses locally installed tools and a version-checking mechanism to get reasonably-sane builds, and runs things (like multiple webservices) in parallel, a la Compose.

There are three main components of `makeup`, the `main.mk` file (which includes tool version checks), individial component `.mk` files (with `build`, `run`, `test`, `env`, and `clean` targets), and the optional generated `Makefile`, which allows anyone who isn't using `makeup` to build and run a project just as easily. 

The `makeup` tool is required to get started (and to get the most top-notch experience), but not required for anyone else to use your project (which is part of the beauty).

✨ Check out the full [usage instructions](./USAGE.md) ✨

### Commands

Implemented:
- `makeup`: builds each component sequentially, and then runs your entire project
- `makeup test` : tests each component sequentially
- `makeup clean` : cleans each of the components in the project
In progress:
- `makeup generate` : generates the main `Makefile` for anyone to use.

That's it for now.

Copyright Connor Hicks and external contributors, 2022. Apache-2.0 licensed, see LICENSE file.