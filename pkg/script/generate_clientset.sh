#!/bin/bash

client-gen --clientset-name versioned --input-base 'github.com/yasker/example-crd/apis' --input "message/v1" --clientset-path github.com/yasker/example-crd/pkg/client/clientset
