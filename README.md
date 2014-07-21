# pew

Pew records the memory usage of an arbitrary process.

## Usage

Use `pew <command>` to run `command` and record its memory usage. Pew
will keep running until the process exits or it receives a SIGINT signal
(Ctrl-C).

Pew will write the memory profile of your process to a timestamped CSV
file inside the `.pew.` subdirectory of the current working directory
from where it's launched.

## README-driven development

The following features have not yet been implemented, but pew is
developed using a README-driven development approach. The following
section is not representative of the current development state, but it's
the roadmap that all development will take place against.

### Multiple git revisions

If it's run from within a git repository, pew picks up the SHA of the
current revision and adds it to the memory profile metadata, so it can
be more easily identified and multiple revisions can be compared against
each other.

This is because, while benchmarking or tracking down memory leaks, it is
often useful to make one small change at a time, and create a separate
git commit for every change, so the history of one's debugging attempts
can be played back later.

This behaviour can be disabled by using the `--no-git` command-line
flag.

### Attach to an existing process

Use the `attach` command to attach pew to an existing process by passing
in the latter's PID.

    $ pew attach 1360

Pew will monitor the process and quit when it senses that the process
has terminated. Using Ctrl-C on pew to stop the data collection early
will have no effect on the process it's monitoring.

### Web server

Use the `server` command to start a local HTTP server that exposes a
minimal, real-time web interface to the data that another instance of
pew is collecting.

    $ pew server
    [pew] HTTP server is listening on 0.0.0.0:7777

Open the address from the command-line output in a web browser to see a
list of the individual samples being collected. Selecting a sample
that's currently in progress will display graphs that are updated in
real time, as new data is collected.

## License

MIT
