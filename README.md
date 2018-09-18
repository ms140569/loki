# Loki Password Manager

**Introduction**

The Loki Password Manager manages passwords/credentials in an comfortable but secure manner. It is highly influence by the Password Store [pass](https://www.passwordstore.org/ "pass") written by Jason A. Donenfeld. The credential data is stored in separate files (with suffix .loki), structured in a directory-tree which contain the following fields:

* Title (String)
* Account (String)
* Password (String)
* Url (String)
* Notes (Multiline String)

_Example:_

```
Title       : Amazon
Account     : example@mail.com
Password    : secret
Url         : https://www.amazon.de

-------------------------------------------
This is my private Amazon shopping account.
```
The individual files are encrypted with the AES-256 algorithm. The store is usually located in the users Homedirectory ($HOME/.loki). Besides the datafiles are two special files:

* .config - Human-editable configuration file (analog to the flags)
* .master - This file keeps track of active generation ( version of master password)

To save the user from authenticate against the store multiple times, the program creates once sucessfully authenticated a daemon process (loki-agentd) which buffers the key in memory. This behavior is similar to the ssh-agent. Subsequent invocations of the loki command fetch the authentification key via unix domain socket from the agent.

**Installation**

The software supports MacOS and Linux (Windows Pull-Requests welcome). Under the Linux a debian package is created, under MacOS the files are copied to there final destination (as long as there are not found on Homebrew). The installation based on the cloned repository is:

* Linux:
```
make install
dpkg -i loki*.deb
```
* MacOS:
```
make install
```

**Usage**

The software is invoked with this basic synopsis:

**loki** \[**[flags]** [**subcommand**] [_record[.loki]]_]

whereas the subcommand could be one of:

* ls | list - Lists the password store in a treelike fashion.
* init - Initialize a new password store.
* login | pw | pass - Authenticate against password store.
* version | ver - Shows version information.
* remove | rm | del - Delete a Record.
* change - Changes the masterpassword in all files.
* help - Shows general help information.
* show - Shows the contents of file.
* insert | add - Inserts new data into file.
* copy | cp - Copy a Record or a subtree.
* move | mv - Moves a Record or a subtree.
* shutdown | stop - Stops the Agent.
* import - Imports a KeepassX CSV file.
* search | grep | find - Searches for given string in all fields and recordnames.
* edit - Edit one Record.

If no command is given, the _list_ subcommand is executed.

The valid _flags_ are:
```
  -b	Blindmode. Do not show password.
  -c	Copy password to clipboard
  -d	Debug mode. Equivalent to -l debug.
  -e	Use external editor given in the EDITOR environment variable.
  -g	Automatically run git commit after each modifiying command.
  -l string
    	Loglevel the program is running with. (default "INFO")
```


**Examples**
```
$ loki
Loki Password Manager, ver 1.0.0, data: /home/mattschmidt/.loki

.
├── frumpy
├── gonzo
├── nas
├── private
│   ├── gmail
│   └── webde
└── public
    ├── accounts
    │   ├── gmail
    │   └── webde
    ├── amazon
    └── google
```

**Manpage**

[Here ist the link to the manpage](man/loki.1.md)

**Datafiles**

Loki's datafiles containing the secret data are structured like this:
```
Datafile format (Big endian)

Magic    : 4c 4f 4b 49     :  4 : "LOKI" Magic Header
Version  : 00 00 00 01     :  4 : v1 - Protocol/Format version
Counter  : 00 00 00 17     :  4 : Version number of Masterpassword
Size     : 00 00 00 00     :  4 : Size of encrypted payload
MD5 Hash : 16 Bytes        : 16 : md5sum of encrypted payload

Data     : .........       : Variable-sized, AES-256-encrypted payload
```

The _Data_ section, variable sized payload is AES-256 encrypted, the un-encrypted payload is formated using Google's [Protocol buffers](https://developers.google.com/protocol-buffers/ "Protocol buffers"):
```
syntax = "proto3";
package storage;

message Record {
    string magic = 1;
    string md5 = 2;
    string title = 3;
    string account = 4;
    string password = 5;
    string url = 6;
    string notes = 7;
}
```
**Git Integration**

Loki's credential data is **not** stored in one single file, but in a collection of separate files arranged in a directory tree. This is to support teams accessing a tree of credential data for different projects, departments etc. and work in parallel on different parts of the tree.

Working on different parts of the directory tree leads to changes on individual files and directories which are expected to be tracked by a version control system. This could be done by any version control system one prefers.

If working on separate files and directories of the data-tree no merge conflicts should occur, but there is one case, where special care is needed:

Changing the Masterpassword (with the _change_ command) changes every file in the tree, should be done only in **one single operation** and **could not** merge with other, normal operations!

If the password store was created using the -g flag, the _.config_ file in the password store will remember this and keep the _Gitmode_ turned on for the store:

```
[basic]
Gitmode = true
```
Having this mode enabled triggers separate git commit operations after each successful tree-modifying operation (insert, import, init, edit, remove, copy, move, change).

**Cryptography**

The password to lock and unlock the password store is processed as UTF-8 String with the Argon2 Key derivation algorithm which produces the 32-byte fixed-sized input to the AES-256 encryption algorithm.

The libaries used are:

* [argon2](https://godoc.org/golang.org/x/crypto/argon2) - External
* crypto/md5 - Standard Go
* crypto/aes - Standard Go
* crypto/cipher - Standard Go
* crypto/rand - Standard Go

**Development**

To build loki you need: 

* Go Version 1.11 (for the Modules support)
* GNU make
* [pandoc](http://pandoc.org) to generate the manpage from Markdown
* protoc from Google's [Protocol buffers](https://developers.google.com/protocol-buffers/ "Protocol buffers")
* [The Go-managed dependencies could be found here](go.mod)
* Under Linux the debian-package generator (_dpkg-deb_)

**Support**

To troubleshoot problems you should used the _Debug_ Loglevel which could be activated using the -d flag (Shorthand for -l debug)

