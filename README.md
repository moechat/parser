MoeParser
=========

A (hopefully) XSS-secured parser for BBCode and chat-style markdown used in
MoeChat, written in pure Go.

GoLang's HTML escape is run *BEFORE* any parsing is done. Therefore, any
vulnerabilties stem from insecure code inside the parser. Remember that this
code has *not* been checked by any professionals, and comes with *NO WARRANTY*.

Tags can be customized in tags.go

This is designed for use inside HTML - *DO NOT* use it inside ```<script>```
tags or ```<style>``` tags!

This code is licensed under the FreeBSD license, described in the LICENSE file.
