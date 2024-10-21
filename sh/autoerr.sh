#!/bin/bash

sh/autotest.sh | grep -E 'Inc|FAIL'
