task "install-deps" {
    description = "Install all of package dependencies"
    pipeline = [
        "go get -t {{.files}}",
    ]
}

task "tests" {
    description = "Run the test suite"
    command = "go test {{.files}}"
}

variables {
    files = "$(go list -v ./... | grep -iEv \"tests|examples\")"
}

