package utils

import (
	"os"

	"github.com/julianstephens/distributed-job-manager/pkg/logger"
)

func ContainsAll[T comparable](mainSlice, subset []T) bool {
	if len(subset) == 0 {
		return true
	}

	mainMap := make(map[T]bool)
	for _, item := range mainSlice {
		mainMap[item] = true
	}

	for _, item := range subset {
		if !mainMap[item] {
			return false
		}
	}
	return true
}

// If mimics the ternary operator s.t. cond ? vtrue : vfalse
func If[T any](cond bool, vtrue T, vfalse T) T {
	if cond {
		return vtrue
	}
	return vfalse
}

// Difference implements slice subtraction s.t. a - b
func Difference(a []string, b []string) []string {
	mb := make(map[string]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var diff []string
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}

// DeleteElement removes an item from a slice at the given index
func DeleteElement[T any](slice []T, index int) []T {
	return append(slice[:index], slice[index+1:]...)
}

// Ensure checks if the given path exists and creates it if not
func Ensure(path string, isDir bool) error {
	var f *os.File
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if isDir {
			if err = os.MkdirAll(path, os.ModePerm); err != nil {
				logger.Errorf("unable to create key pair dir: %v", err)
				return err
			}
		} else {
			f, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
			if err != nil {
				logger.Errorf("unable to open key pair file: %v", err)
				return err
			}
			f.Close()
		}
	}

	return nil
}

// StringPtr returns a pointer to the given string.
func StringPtr(s string) *string {
	return &s
}
