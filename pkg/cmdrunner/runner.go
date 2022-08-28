package cmdrunner

import "fmt"

// Runner provides the neccessary methods to
// collect and compare data from a target device.
type Runner interface {
	Run() error
	Compare() error
}

// Stepper steps through a list of runners, grabbing and comparing
// the received data while validating success.
func Stepper(cmds []Runner) error {
	for _, cmd := range cmds {
		fmt.Printf("[%v] running...\n", cmd)
		if err := cmd.Run(); err != nil {
			return err
		}
		fmt.Printf("[%v] validating...\n", cmd)
		if err := cmd.Compare(); err != nil {
			return err
		}
		fmt.Printf("[%v] validation passed\n", cmd)
	}
	return nil
}
