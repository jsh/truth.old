# Tests

These tests answer three basic questions about the suite of behaviors that define true(1).


## Questions

1. Does it terminate?
1. Does it return success?
1. Is the output correct?

## Behaviors

- true should return success (exit code 0)
- true --help should return the help message
- true --version should output the (correct) version number
- true should ignore other command-line arguments
- true should handle I18N
- true should ignore multiple args

## Prerequisites

To test for I18N, you must install a language pack.
Run setup-i18-testing.
