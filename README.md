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

Copies all messages from source queue to target queue.

With `-e exchange` parameter it is possible to send data into exchange instead of queue, queue name
in that case will be used as routing key. If queue name start with `#` hare will replace routing key
with corresponding header value. For example passing `#class` as queue causes reading routing key 
from `class` header of every message.

```
$ hare copy --help

Usage:
  hare copy source target [flags]

Flags:
  -e, --exchange string   Exchange to use for copied messages
```


#### move

Moves messages from one queue to another

With `-e exchange` parameter it is possible to send data into exchange instead of queue, queue name
in that case will be used as routing key. If queue name start with `#` hare will replace routing key
with corresponding header value. For example passing `#class` as queue causes reading routing key
from `class` header of every message.

```
$ hare move --help

Usage:
  hare move source target [flags]

Flags:
  -e, --exchange string   Exchange to use for moved messages
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

Uploads messages from local file to RabbitMQ. 

With `-e exchange` parameter it is possible to send data into exchange instead of queue, queue name
in that case will be used as routing key. If queue name start with `#` hare will replace routing key
with corresponding header value. For example passing `#class` as queue causes reading routing key
from `class` header of every message.


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