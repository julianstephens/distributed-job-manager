package httputil

import (
	"bufio"
	"errors"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/julianstephens/distributed-task-scheduler/pkg/logger"
)

type Item struct {
	Key   string
	Value any
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

// GetId parses the 'id' path param from Gin context
func GetId(ctx *gin.Context) (string, error) {
	id := ctx.Param("id")
	if id == "" {
		return id, errors.New("no resource id provided")
	}
	return id, nil
}

func ReadFile(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	stats, statsErr := file.Stat()
	if statsErr != nil {
		return nil, statsErr
	}

	var size int64 = stats.Size()
	bytes := make([]byte, size)

	bufr := bufio.NewReader(file)
	_, err = bufr.Read(bytes)

	return bytes, err
}
