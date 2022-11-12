.PHONY: genall
genall:
	hack/update-codegen.sh
	hack/update-crdgen.sh

.PHONY: vendor
vendor:
	go mod tidy && go mod vendor