#!/bin/bash

set -x

GOPATH=$(go env GOPATH)
PACKAGE_NAME=github.com/appscode/stash

pushd "$GOPATH/src/$PACKAGE_NAME"

# Generate defaults
defaulter-gen \
    --v 1 --logtostderr \
    --go-header-file "hack/gengo/boilerplate.go.txt" \
    --input-dirs "$PACKAGE_NAME/apis/stash" \
    --input-dirs "$PACKAGE_NAME/apis/stash/v1alpha1" \
    --extra-peer-dirs "$PACKAGE_NAME/apis/stash" \
    --extra-peer-dirs "$PACKAGE_NAME/apis/stash/v1alpha1" \
    --output-file-base "zz_generated.defaults"

# Generate deep copies
deepcopy-gen \
    --v 1 --logtostderr \
    --go-header-file "hack/gengo/boilerplate.go.txt" \
    --input-dirs "$PACKAGE_NAME/apis/stash" \
    --input-dirs "$PACKAGE_NAME/apis/stash/v1alpha1" \
    --output-file-base zz_generated.deepcopy

# Generate conversions
conversion-gen \
    --v 1 --logtostderr \
    --go-header-file "hack/gengo/boilerplate.go.txt" \
    --input-dirs "$PACKAGE_NAME/apis/stash" \
    --input-dirs "$PACKAGE_NAME/apis/stash/v1alpha1" \
    --output-file-base zz_generated.conversion

# Generate the internal clientset (client/clientset_generated/internalclientset)
client-gen \
   --go-header-file "hack/gengo/boilerplate.go.txt" \
   --input-base "$PACKAGE_NAME/apis/" \
   --input "stash/" \
   --clientset-path "$PACKAGE_NAME/client/" \
   --clientset-name internalclientset

# Generate the versioned clientset (client/clientset_generated/clientset)
client-gen \
   --go-header-file "hack/gengo/boilerplate.go.txt" \
   --input-base "$PACKAGE_NAME/apis/" \
   --input "stash/v1alpha1" \
   --clientset-path "$PACKAGE_NAME/" \
   --clientset-name "client"

# generate lister
lister-gen \
   --go-header-file "hack/gengo/boilerplate.go.txt" \
   --input-dirs="$PACKAGE_NAME/apis/stash" \
   --input-dirs="$PACKAGE_NAME/apis/stash/v1alpha1" \
   --output-package "$PACKAGE_NAME/listers"

# generate informer
informer-gen \
   --go-header-file "hack/gengo/boilerplate.go.txt" \
   --input-dirs "$PACKAGE_NAME/apis/stash" \
   --input-dirs "$PACKAGE_NAME/apis/stash/v1alpha1" \
   --internal-clientset-package "$PACKAGE_NAME/client/internalclientset" \
   --versioned-clientset-package "$PACKAGE_NAME/client" \
   --listers-package "$PACKAGE_NAME/listers" \
   --output-package "$PACKAGE_NAME/informers"

#go-to-protobuf \
#  --proto-import="${KUBE_ROOT}/vendor" \
#  --proto-import="${KUBE_ROOT}/third_party/protobuf"

popd
