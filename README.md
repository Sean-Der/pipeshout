# pipeshout

Do you spend most of your life tailing logfiles? Are you piping input into multiple programs, just trying to track down a few log lines?
Dread everytime you have to ask someone to join your tmux/screen session to review logs, and then fumble with shared input?

pipeshout is a small web app, that consumes data via a *Unix Domain Socket* and then allows multiple users to view those lines and apply regexes to them.
Here are a few of its features live

* **Streaming Input** - Pipeshout uses websockets, so feels just like using `tail`
![alt text](http://i.imgur.com/dkcHcyL.gif "Streaming")

* **Regexes** - Apply regexes on prefix OR line, allowing you to quickly examine data
![alt text](http://i.imgur.com/Er4hfNi.gif "Regexes")

* **Consumes Lines** - You aren't just tied to files! You can send it tcpdump, scheduled jobs or `sl`

* **Persistence** - Pipeshout stores input even when no clients aren't connected, allowing you to open it right when an error occurs

## Requirements
* **Go**

## Running

### Installing
* `go get github.com/Sean-Der/pipeshout`
* `$GOPATH/src/github.com/Sean-Der/pipeshout`
* `git submodule init && git submodule update`
* `go run *.go`

Pipeshout is now started! By default it runs on port `8080`, so to access it on your local machine visit `http://localhost:8080`

### Adding lines
As a first test try `echo 'foo bar' | nc -U pipeshout.pipe` in the pipeshout directory. You should then get a single line with the prefix `foo` and the line `bar`

pipeshout expects each line to be made up of `$prefix $line`, a prefix is used to denote a single source of input like a logfile or program.

For example if you wanted to watch your maillog you could use

`tail -f /var/log/maillog | sed 's/^/maillog / | nc -U pipeshout.pipe`


## Contributing/TODO
All contributions are welcome! Even if you only get something half done open a PR! I will finish, and merge if it is a good feature.

Want to sharpen your Go/react.js chops? Here are some cool things I would love to see, will eventually work on myself!

* Code quality, adding JS/CSS/HTML linting and make sure the code is idiomatic
* Desktop notifications on regex matches
* Error handling
* Freeze log handling, so if you see an error you can quickly freeze to examine it

## License
The MIT License (MIT)

Copyright (c) 2014 Sean DuBois

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
