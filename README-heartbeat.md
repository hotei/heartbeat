<center>
# heartbeat
</center>

## What Is It?

_heartbeat_ is a go program that demonstrates a simple client server system.
The server collects heartbeat signals from multiple clients and displays them
in a list.


### Installation

If you have a working go installation on a Unix-like OS:

> ```go get github.com/hotei/heartbeat```

Will copy github.com/hotei/heartbeat to the first entry of your $GOPATH

or if go is not installed yet :

> ```cd DestinationDirectory```

> ```git clone https://github.com/hotei/heartbeat.git```

### Features

You can provide an "expected" client list such as:
``` go
	newClient("JOHN", "10.1.2.124")
	newClient("MARY", "10.1.1.113")
	newClient("FRED", "10.1.2.112")
```

Uses UDP so heartbeat delivery is not guaranteed - see also bugs.

Demonstrates how client and server can be same program.

If not in the list heartbeat respondees will still show up but with
numbers instead of names.

### Limitations

* Uses UDP so heartbeat delivery is not guaranteed - see also features.

### Usage

Typical usage is :

On clients:
	heartbeat -port=2345

On server:
	heartbeat -server -port=2345

### BUGS

none known

### To-Do

* Nice:
	* ServerIPstr should be gotten from command line or automatically,not hardcoded
	* Get client IP from flag/env/syscall?
	* Sort reports by either name or ipnumber (switch without restart)
* Nice but no immediate need:
	* Add some randomness in client repeater timing to avoid floods?  Minor benefit
	unless we have tons of clients in short interval loop
	* What to do if server consistantly not writable?
	* Use _hotei/ansiterm_ to make updates pretty

### Change Log

* 2015-05-25 validated working with go 1.4.2
* 2013-03-28 sorted list of system names working
* 2010-02-21 started, working same day
 
### Resources

* [go language reference] [1] 
* [go standard library package docs] [2]
* [Source for program] [3]

[1]: http://golang.org/ref/spec/ "go reference spec"
[2]: http://golang.org/pkg/ "go package docs"
[3]: http://github.com/hotei/heartbeat "github.com/hotei/heartbeat"

Comments can be sent to <hotei1352@gmail.com> or to user "hotei" at github.com.

License
-------
The 'heartbeat' go package/program is distributed under the Simplified BSD License:

> Copyright (c) 2010-2015 David Rook. All rights reserved.
> 
> Redistribution and use in source and binary forms, with or without modification, are
> permitted provided that the following conditions are met:
> 
>    1. Redistributions of source code must retain the above copyright notice, this list of
>       conditions and the following disclaimer.
> 
>    2. Redistributions in binary form must reproduce the above copyright notice, this list
>       of conditions and the following disclaimer in the documentation and/or other materials
>       provided with the distribution.
> 
> THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDER ``AS IS'' AND ANY EXPRESS OR IMPLIED
> WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND
> FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL <COPYRIGHT HOLDER> OR
> CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
> CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
> SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON
> ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING
> NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF
> ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

Documentation (c) 2015 David Rook 

// EOF README-heartbeat.md  (markdown tested with blackfriday)

