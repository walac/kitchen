This is a sample order delivery simulator. You basically have three components:

* The supplier: which supplies new orders to the system.
* The kitchen: which receives orders from the supplier and put ready orders
on the shelf.
* The courier: that picks up ready order from the shelf and deliveries it.

We have 4 kinds of shelfs: for hot, cold and frozen food, and the overflow shelf,
which we put orders when their corresponding shelf is full. The orders have a
kind of "time to live" property, here called "waste". When the order is wasted,
it is removed from the shelf.

Build and running
=================

To build it type `go build ./cmd/kitchen` and to run the tests `go test ./...`.
To run the program type `./kitchen data/orders.json`. The command also accepts an
argument to control the rate in which the kitchen receives new orders, in orders/sec.
By default it is 2 orders/sec. Example: `./kitchen -r 10 data/orders.json` will run
the simulator with an arrival rate of 10 orders/sec.

To exit the application, type `ctrl-c`.
