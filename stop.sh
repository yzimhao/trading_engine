#!/bin/bash


pgrep haobase | xargs kill
pgrep haomatch | xargs kill
pgrep haoquote | xargs kill
pgrep haoadm | xargs kill
pgrep example | xargs kill