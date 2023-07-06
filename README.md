# Gossip Glomers

## A series of distributed systems challenges brought to you by Fly.io

### Introduction

We've teamed up with Kyle Kingsbury, author of Jepsen, to build this series of distributed systems challenges so you can try your hand and see how your skills stack up.

The challenges are built on top of a platform called Maelstrom, which in turn, is built on Jepsen. This platform lets you build out a "node" in your distributed system and Maelstrom will handle the routing of messages between the those nodes. This lets Maelstrom inject failures and perform verification checks based on the consistency guarantees required by each challenge.

The documentation for these challenges will be in Go, however, Maelstrom is language agnostic so you can rework these challenges in any programming language.

[https://fly.io/dist-sys/](https://fly.io/dist-sys/)

### Challenges

- [x] 1: Echo
- [x] 2: Unique ID Generation
- [x] 3a: Broadcast - Single Node Broadcast
- [x] 3b: Broadcast - Multi-Node Broadcast
- [x] 3c: Broadcast - Fault Tolerant Broadcast
