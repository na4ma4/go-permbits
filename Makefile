-include .makefiles/Makefile
-include .makefiles/pkg/go/v1/Makefile
-include .makefiles/ext/na4ma4/lib/golangci-lint/v1/Makefile
-include .makefiles/ext/na4ma4/lib/goreleaser/v1/Makefile

.makefiles/ext/na4ma4/%: .makefiles/Makefile
	@curl -sfL https://raw.githubusercontent.com/na4ma4/makefiles-ext/main/v1/install | bash /dev/stdin "$@"

.makefiles/%:
	@curl -sfL https://makefiles.dev/v1 | bash /dev/stdin "$@"

######################
# Testing
######################

GINKGO := artifacts/ginkgo/bin/ginkgo
$(GINKGO):
	@mkdir -p "$(MF_PROJECT_ROOT)/$(@D)"
	GOBIN="$(MF_PROJECT_ROOT)/$(@D)" go get github.com/onsi/ginkgo/ginkgo

test:: $(GINKGO)
	-@mkdir -p "artifacts/test"
	$(GINKGO) -outputdir "artifacts/test/" -r --randomizeAllSpecs --randomizeSuites --failOnPending --cover --trace --race --compilers=2 --nodes=2
