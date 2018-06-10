# *mgmt*: This is: mgmt!

[![Build Status](https://secure.travis-ci.org/purpleidea/mgmt.png)](http://travis-ci.org/purpleidea/mgmt)
[![Documentation](https://img.shields.io/docs/markdown.png)](DOCUMENTATION.md)
[![IRC](https://img.shields.io/irc/%23mgmtconfig.png)](https://webchat.freenode.net/?channels=#mgmtconfig)

## Documentation:
Please see: [DOCUMENTATION.md](DOCUMENTATION.md) or [PDF](https://pdfdoc-purpleidea.rhcloud.com/pdf/https://github.com/purpleidea/mgmt/blob/master/DOCUMENTATION.md).

## Questions:
Come join us in [#mgmtconfig](https://webchat.freenode.net/?channels=#mgmtconfig) on Freenode!

## Examples:
Please look in the [examples/](examples/) folder!

## Notes:
* This is currently a research project into next generation config management technologies!
* This is my first complex project in golang, please notify me of any issues.
* I have some well thought out designs for the future of this project, which I'll try and write up clearly and publish as soon as possible.
* Please don't expect stable interfaces, code, or any data safety.
* This design is the result of ideas I've had from hacking on advanced config management projects.
* I first started hacking on this in ~2013, even though I had very little time for it.
* I couldn't think of a good name for the project, so it's now being called `mgmt` until someone contributes a better one!
* I've published a number of articles about this tool:
  * [https://ttboj.wordpress.com/2016/01/18/next-generation-configuration-mgmt/](https://ttboj.wordpress.com/2016/01/18/next-generation-configuration-mgmt/)
* There are some screencasts available:
  * TODO

## Dependencies:
* golang (required, available in most distros)
* golang libraries (required, available with `go get`)
  ```
  go get github.com/coreos/etcd/client
  go get gopkg.in/yaml.v2
  go get gopkg.in/fsnotify.v1
  go get github.com/codegangsta/cli
  go get github.com/coreos/go-systemd/dbus
  go get github.com/coreos/go-systemd/util
  ```
* pandoc (optional, for building a pdf of the documentation)
* graphviz (optional, for building a visual representation of the graph)

## Patches:
We'd love to have your patch! Please send it by email, or as a pull request.

##

Happy hacking!
