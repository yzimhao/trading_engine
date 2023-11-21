#!/bin/bash

pgrep example | xargs kill

pgrep haobase | xargs kill
pgrep haomatch | xargs kill
pgrep haoquote | xargs kill
pgrep haoadm | xargs kill
