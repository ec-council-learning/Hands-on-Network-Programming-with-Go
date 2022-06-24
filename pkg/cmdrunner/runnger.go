package cmdrunner

import "fmt"

type Runner interface {
	Run() error
	Compare() error
}

func Stepper(cmds []Runner) error {
	for _, cmd := range cmds {
		fmt.Printf("[%v] running\n", cmd)
		if err := cmd.Run(); err != nil {
			return err
		}
		fmt.Printf("[%v] validating\n", cmd)
		if err := cmd.Compare(); err != nil {
			return err
		}
		fmt.Printf("[%v] validation passed\n", cmd)
	}
	return nil
}
