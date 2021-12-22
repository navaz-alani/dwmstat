# dwmstat

A simple event-driven modular status-bar for dwm.
It is written in Go & uses Unix sockets for signaling. 

The status bar is conceptualized as a set of modules, which are identified by
their names (e.g. `battery`, `date`).
Each module, when executed, outputs a string which is used in the status-bar.
A module has one main property: an interval duration, which specifies how often
the module should be re-executed.
If this duration is negative, then the module is never updated.

The more flexible method of forcing modules to update is through "signals",
which should not be mistaken with OS signals (e.g. `SIG_INT`, `SIG_SEGV`).
When `dwmstat` is run, it opens up a Unix socket on `/tmp/dwmstat_sig`, which
can be interacted with using a program like `socat`.
For example, the command (which can be easily shorthanded to a script)
```bash
printf 'volume' | socat UNIX-CONNECT:/tmp/dwmstat_sig -`
```
would send the text `volume` to the socket and `dwmstat` will interpret this as
"update the 'volume' module".
The socket functionality is quite simple right now, but due to its flexibility,
there is opportunity to do quite a bit more, e.g. activate/deactivate named
modules in the status-bar, update module properties, etc.
This is why this approach was preferred to the OS signal approach.
