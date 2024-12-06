# Golang Archipelago Client

This is a golang archipelago randomizer client. The goals are to provide all basic client functions. As part of this library, there is an openapi definition for all types.

## Licensing

The default license of the repository is MIT licensed.
The goal is to provide these files/this library in as permissive a method as possible.
If the MIT license does not meet your needs, please get in contact.
All we ask (tho will not require) is that attribution be kept to recognize the work done here.

## construction

The client library is broken up into several parts.
THe main parts however are

- the websocket handler, which is located in the Network folder,
- the datatypes, which are defined via openapi in the api folder,
- the caching client located in the root folder.
