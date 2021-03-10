:rabbit2: Hare
==============

[RabbitMQ](https://www.rabbitmq.com/) CLI toolkit.

### Basics

Download and compile Hare:
```bash
git clone https://github.com/mono83/hare.git
cd hare
mkdir -p target
go build -o target/hare main.go
target/hare --help
```

Or build Docker image:
```bash
docker build -t hare .
docker run --rm hare --help
```

RabbitMQ address and credentials are configured using `--uri` command line 
argument and have default value `amqp://guest:guest@localhost:5672/` 


#### Commands

#### ping

Tests current RabbitMQ connection settings.

```
$ hare ping --help 

Usage:
  hare ping [flags]
Aliases:
  ping, test
```


#### visit

Reads messages from queue and then requeues them without `redelivered` flag

```
$ hare visit --help

Usage:
  hare visit queue [flags]
Aliases:
  visit, view, look
Flags:
  -c, --count int   Count to get (default 1)
```


#### copy

Copies all messages from source queue to target queue

```
$ hare copy --help

Usage:
  hare copy source target [flags]
```


#### move

Moves messages from one queue to another

```
$ hare move --help

Usage:
  hare move source target [flags]
```

#### download

Downloads messages from queue to local file (line-separated JSON). By default will not delete messages from queue,
provide `-d` flag to change behaviour

```
$ hare download --help

Usage:
  hare download queue filename [flags]
Aliases:
  download, save, dump, down, flush
Flags:
  -a, --append   If true will append to file instead of replace it
  -d, --delete   If true, will delete messages from queue
```

#### upload

Uploads messages from local file to rabbitMQ. By default will upload into `queue`, passed as CLI argument, 
but if `exchange` flag is provided message will be sent to specified exchange with routing key equal to `queue`.

```
$ hare upload --help

Usage:
  hare upload queue filename [flags]
Aliases:
  upload, restore, up
Flags:
  -e, --exchange string   Exchange to use for uploaded messages
  -r, --replicate int     Amount of copies per single line (default 1)
```

#### circuit-breaker

Starts consumer with [high priority](https://www.rabbitmq.com/consumer-priority.html) that will intercept all messages in given queue.

```
Usage:
  hare circuit-breaker queue filename [flags]

Aliases:
  circuit-breaker, cb, break, circuit-break

Flags:
  -a, --append         If true will append to file instead of replace it
  -p, --priority int   Consumer priority (default 10)
```