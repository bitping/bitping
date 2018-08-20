## Bitping

Bitping is the root of all the things.

## Getting started

Most of the work we'll do within `bitping` is through the `Makefile`. Checkout the `Makefile` for details about how these things work.

Basically, the way that this works is that `bitping` listens on all the available blockchains for block events. This can be events where blocks are mined or contracts are called.

Blockchain status:

- [x] ethereum
- [ ] bitpoing
- [ ] eos

Feel free to add another blockchain here. We'll add the as we see fit and the need. To add a blockchain to the `watch` command is straight-forward. Each blockchain needs to implement the `blockchain` interface and implement the following methods:

- `NewClient()`
- `Run()`

## What it is

## License

Ari Lerner
