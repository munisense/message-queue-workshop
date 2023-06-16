# Syntax Workshop June 2023

Welcome!

Today we will be showing what a Message Queue is and why it's awesome. We will be doing that in conjunction with the
programming language [Golang](https://go.dev/), because it's builtin concurrency makes it a breeze to work with a
Message Queue. First of all lets make sure everybody is set up.

## Prerequisites

- Install Golang ([go.dev/doc/install](https://go.dev/doc/install))
- Verify the installation by running `go version` in your terminal.
    - If the output is similar to `go version go1.20.2 linux/amd64` you are good to Go!
- Install an IDE, we recommend GoLand from JetBrains ([jetbrains.com/go](https://www.jetbrains.com/go/))
    - Should be free to use for students.
- Download this repository to your laptop.
    - `git clone https://github.com/munisense/syntax-workshop-2023.git`

# Step 1: Hello World

We will be using `go run` here, this compiles and runs your code in a single command.

```shell
go run cmd/01_hello_world/main.go
```

# Step 2: Getting a message from a queue

Copy the file `.env.default` to a new file called `.env`, and fill in the credentials.

```shell
go run cmd/02_get_from_a_queue/main.go
```

**Question**: Does everybody get all data?

# Step 3: Consuming from a queue

Now we will all be reading (consuming) messages of a queue.

```shell
go run cmd/03_consume_from_a_queue/main.go
```
