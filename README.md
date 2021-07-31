# Task

Implement a tcp server for fragmented packets server functions:
    
    1. Send messages to all clients connected.
    2. Send messages to clients specified (you should tag the client)

Implement a client connect to the server.

## Message format

    | 2 bytes | x  bytes |
    | content length | content|

* Server and client use same message format for communication.

### Example

- msg1: [00 02 41 42] => content length 2, content "AB"    
- msg2: [00 03 41 42 41] => content length 3, content "ABA"    

## Test case

10 clients, get broadcast messages and client 1 send message to client 2.

# Setup

## Environment variables

See `.envDist`

## Tag

Message tag should be at the begining of content:

```
length 10, content "#test#some"
```