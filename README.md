# 2025gopherconf.co.za
Networking talk for GopherCon South Africa

Examples are indexed sequentially, with each program of a given index operating in unison.

There is also a helpers.go file which contains support capabilities used generally across
the indexed programs.

To run any individual program:

		go run ''name''.go helpers.go [command-line parameters]

Client programs take one or more filenames as parameters which they pass to Server programs
using a variety of network APIs.