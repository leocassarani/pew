# pew

Pew records the memory usage of an arbitrary process.

## Usage

Use `pew <command>` to run `command` and record its memory usage. Pew
will keep running until the process exits or it receives a SIGINT signal
(Ctrl-C).

Pew will write the memory profile of your process to a timestamped CSV
file inside the `.pew.` subdirectory of the current working directory
from where it's launched.
