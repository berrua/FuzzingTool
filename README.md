# FuzzingTool

This repository contains a CLI-based fuzzing tool developed in Go, utilizing the concurrency provided by the worker pool working mechanism.

![asciiart](asciiart.png)

Usage:
  fuzzing fuzzer [flags]

Flags:
  -h, --help              help for fuzzer
      --speed int         Specify the speed in milliseconds for fuzzing (default 500) 
      --status string     Specify the status code to filter the results (ex: 200, 403)
      --timeout int       Specify the timeout in milliseconds for HTTP requests (default 5000)
      --url string        Specify the target URL for fuzzing
      --wordlist string   Specify the name of the wordlist file
