# Message Queue Workshop

Welcome!

Today we will be showing what a Message Queue is and why it's awesome. We will be doing that in conjunction with the
programming language [Golang](https://go.dev/), because it's builtin concurrency makes it a breeze to work with a
Message Queue. First of all lets make sure everybody is set up.

## Prerequisites

- Install Golang ([go.dev/doc/install](https://go.dev/doc/install))
- Verify the installation by running `go version` in your terminal.
    - If the output is similar to `go version go1.21.0 linux/amd64` you are good to Go!
- Install an IDE, we recommend GoLand from JetBrains ([jetbrains.com/go](https://www.jetbrains.com/go/))
    - Should be free to use for students.
- Download this repository to your laptop.
    - `git clone https://github.com/munisense/message-queue-workshop.git`

# Step 1: Hello World

We will be using `go run` here, this compiles and runs your code in a single command.

```shell
go run cmd/01_hello_world/main.go
```

# Step 2: Getting a message from a queue

Copy the file `.env.default` to a new file called `.env`, and fill in the credentials (see the presentation).

```shell
go run cmd/02_get_from_a_queue/main.go
```

# Step 3: Consuming from a queue

Now we will all be reading (consuming) messages of a queue. This application will keep running until you stop it with `ctrl-c`.

```shell
go run cmd/03_consume_from_a_queue/main.go
```

# Question: Does everybody get all the data?

Yes? No? Why?

# Step 4: Getting all the data

As seen in the presentation: "wie het eerst komt die het eerst maalt". We will now create our own (exclusive) queue in order to guarantee we receive all the data.

```shell
go run cmd/04_consume_all_messages_from_a_queue/main.go
```

# Step 5: Golang also has queues

Oh sorry, they are called "channels" actually and don't have the "routingkey" that RabbitMQ uses. However, you can still use it to shovel data around.
Perhaps we want to parse our JSON string into an object. Possible that is an expensive operation, or maybe you just want to split off the logic.
Now, consume from the amqp queue, log it, put the message on a channel and have another Go routine handle the parsing!

```shell
go run cmd/05_golang_also_has_queues/main.go
```

# Step 6: Publishing!

Now you will be publishing your own message to the queue. Can you also read back this message by changing the routingkey in the code from step 4?

```shell
go run cmd/06_publish_to_an_exchange/main.go
```
