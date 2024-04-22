#!/bin/bash

pgrep example | xargs kill
sleep 3s
pgrep haobase | xargs kill
pgrep haomatch | xargs kill
pgrep haoquote | xargs kill
pgrep haoadm | xargs kill
pgrep haotrader | xargs kill
