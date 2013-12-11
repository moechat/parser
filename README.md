MoeChat Parser
==============

An XSS-secured parser for BBCode and chat-style markdown used in MoeChat,
written in GoLang.

The parser uses Go's html/template package to inject all user input.
Remember that this code has *not* been checked by any professionals, and
comes with **NO WARRANTY**.

This is designed for use inside HTML - *DO NOT* use it inside ```<script>```
tags or ```<style>``` tags!

This code is licensed under the FreeBSD license, described in the LICENSE file.
