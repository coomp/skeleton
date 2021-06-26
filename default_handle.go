package skeleton

import "fmt"

func NewDefaultHandle() HandlerFunc {
	return func(c Context) error {
		fmt.Printf("NewDefaultHandle %v\n", c.GetString())
		return nil
	}
}
