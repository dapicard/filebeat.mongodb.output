# filebeat.mongodb.output
A Filebeat embedding a MongoDB output.

You can use the output in every beat you want. This repository offers a Filebeat "main" that embeds it.
You can compile and use it by following the Golang setup detailed in the CONTRIBUTE instructions of beats :
https://www.elastic.co/guide/en/beats/devguide/current/beats-contributing.html#setting-up-dev-environment

And compiling this repository using :
* go build

To build the filebeat, and :

* ./filebeat.mongodb.output
To launch it.
