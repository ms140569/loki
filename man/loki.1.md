% loki(1) Version 1.0 | loki Password Manager

NAME
====

**loki** — Loki Password Manager

SYNOPSIS
========

**loki** \[**[flags]** **subcommand** _record[.loki]]_

DESCRIPTION
===========

loki is a tool to manage your various credentials in AES256 encrypted file located in your homedirectory.
Default location is $HOME/.loki


copy | cp - Copy a Record or a subtree.   Example: loki [flags] copy <file|dir> <file|dir>

move | mv - Moves a Record or a subtree.   Example: loki [flags] move <file|dir> <file|dir>

insert | add - Inserts new data into file.   Example: loki [flags] insert filename

import - Imports a KeepassX CSV file.   Example: loki [flags] import keepass-filename

search | grep | find - Searches for given string in all fields and recordnames.   Example: loki [flags] search <querystring>

version | ver - Shows version information.   Example: loki [flags] version

complete - Generates bash programmable completion statement.   Example: loki [flags] complete

init - Initialize a new password store.   Example: loki [flags] init

login | pw | pass - Authenticate against password store.   Example: loki [flags] login

edit - Edit one Record.   Example: loki [flags] edit filename

remove | rm | del - Delete a Record.   Example: loki [flags] remove filename

shutdown | stop - Stops the Agent.   Example: loki [flags] shutdown

help - Shows general help information.   Example: loki [flags] help

ls | list - Lists the password store in a treelike fashion.   Example: loki [flags] ls

show - Shows the contents of file.   Example: loki [flags] show filename

change - Changes the masterpassword in all files.   Example: loki [flags] change


Options
-------

-b

:   Blindmode. Do not show password.

-c

:   Copy password to clipboard

-d

:   Debug mode. Equivalent to -l debug.

-e

:   Use external editor given in the EDITOR environment variable.

-g

:   Automatically run git commit after each modifiying command.

-l <{Off, Trace, Debug, Info, Warning, Error, Fatal, All}>

:   Loglevel the program is running with. (default "Info")



FILES
=====

*$LOKI_BASE/.config*

:   Configfile to preset configuration without using commandline switches.

*$LOKI_BASE/.master*

:   Changes to the masterpassword (generation) are tracked here.

*$LOKI_BASE/<subidr>/<filename>.loki*

:   Single Loki record named "filename" stored in files with the suffix .loki and optionally organized in
    subdirectories.

ENVIRONMENT
===========

**LOKI_BASE**

:   Used to override Loki's default store location $HOME/.loki

**EDITOR**

:   When editing Loki records using an external editor (-e) this variable traditionally
    points to the editor to use.

**LOKI_LOGLEVEL**

:   Used to specify the Loglevel of the program. Equivallent to the -l <logleve> flag


BUGS
====

See GitHub Issues: <https://github.com/[owner]/[repo]/issues>

AUTHOR
======

Matthias Schmidt (matthias.schmidt@gmail.com)

SEE ALSO
========

**pass(1)**
