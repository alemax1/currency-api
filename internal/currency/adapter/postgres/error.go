package postgres

import "fmt"

const duplicateErrorCode = "23505"

func newExecContextErr(err error) error {
	return fmt.Errorf("exec context: %w", err)
}

func newUpdatedRowsErr(err error) error {
	return fmt.Errorf("get updated rows: %w", err)
}

func newQueryErr(err error) error {
	return fmt.Errorf("query: %w", err)
}

func newScanErr(err error) error {
	return fmt.Errorf("scan: %w", err)
}

func newRowsErr(err error) error {
	return fmt.Errorf("rows: %w", err)
}
