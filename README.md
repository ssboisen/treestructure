Hello Stranger
====

You know the drill: `docker-compose up`

Afterwards you can hit up `http://localhost:8080/`

More specifically you could do:

`curl -XPOST http://localhost:8080/nodes/f/o/o/b/a/r` to create a few nodes

`curl -XPOST http://localhost:8080/nodes/f/o/o/b/a/z` to create some additional nodes

`curl -XPOST http://localhost:8080/nodes/f/o/o/q/u/e/x` to create even more nodes

`curl -XPOST http://localhost:8080/nodes/t/r/a/d/e/s/h/i/f/t` to create the best nodes

`curl http://localhost:8080/nodes/` to list existing children of the root node

`curl http://localhost:8080/nodes/f` to list existing children of the `/f` node

`curl http://localhost:8080/nodes/f/o/o` to list existing children of the `/f/o/o` node

`curl -XPUT http://localhost:8080/nodes/f/o/o/q?newParent=/` to assign `/f/o/o/q` the root node as its new parent.
