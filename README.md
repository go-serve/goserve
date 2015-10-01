goserve
=======

[![Travis CI results][travis]](https://travis-ci.org/yookoala/goserve)

[travis]: https://api.travis-ci.org/yookoala/goserve.svg?branch=master

`goserve` is a small utility to serve static HTML files in a directory to a given port.

It is intented for developer who wants a quick access to their code through browser. Especially so when their code are not using relative path in link / image / CSS what so ever. Just compile the binary and put it in your `PATH`, then it is good to go.

The code is dead simple. I simply don't want to write it all the time. And I wish it might be of help to you, too.


Requirement
-----------
`goserve` requires only the core go libraries. No need to `go get` anything other than this.


Installation
------------
If you have set `GOPATH/bin` to your `PATH`, you may install and use this by:

```sh
go get github.com/yookoala/goserve
```

Alternatively, you may compile and copy the binary to your directory in `PATH`.

```sh
git clone https://github.com/yookoala/goserve.git
cd goserve
go build
cp goserve YOUR_DIR_IN_PATH/.
```


Usage
-----

Just type this, a server will be serving the files in the current directory to default port 8080:

```sh
goserve
```

To specify the directory, you may add 1 directory path as the argument:

```sh
goserve ./data
```

You may specify the port with environment variable `PORT`:

```sh
PORT=8123 goserve ./data
```

You may manually override the default or `PORT` with `-port` parameter:
```sh
goserve -port=8123 ./data
```


Author
------
This software is written by [Koala Yeung](https://github.com/yookoala) (koalay at gmail.com).


Licence
-------
This software is licenced under GPL v3. You may obtain a copy of the licence in the `LICENSE` file in this repository.


Bug Report
----------
You are always welcome to report issue here:
https://github.com/yookoala/goserve/issues