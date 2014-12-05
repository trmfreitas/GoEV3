// Provides APIs for interacting with EV3's motors.
package Motor

import (
	"fmt"
	"github.com/mattrajca/GoEV3/utilities"
	"log"
	"io/ioutil"
	"strings"
)

// Constants for output ports.
type OutPort string

const (
	OutPortA OutPort = "outA"
	OutPortB         = "outB"
	OutPortC         = "outC"
	OutPortD         = "outD"
)


func findFilename(port OutPort) string {
	motors, _ := ioutil.ReadDir("/sys/class/tacho-motor")
	for _, item := range motors {
		if strings.HasPrefix(item.Name(), "motor") {
			motorPath := fmt.Sprintf("/sys/class/tacho-motor/%s", item.Name())
			portr := utilities.ReadStringValue(motorPath, "port_name")
			if OutPort(portr) == port {
				return fmt.Sprintf("/sys/class/tacho-motor/%s", item.Name())
			}
		}
	}

	log.Fatalf("Could not find %v motor\n", port)

	return ""
}

// Runs the motor at the given port.
// `speed` ranges from -100 to 100 and indicates the target speed of the motor, with negative values indicating reverse motion. Depending on the environment, the actual speed of the motor may be lower than the target speed.
func Run(port OutPort, speed int16) {
	if speed > 100 || speed < -100 {
		log.Fatal("The speed must be in range [-100, 100]")
	}

	utilities.WriteIntValue(findFilename(port), "run", 1)
	utilities.WriteIntValue(findFilename(port), "duty_cycle_sp", int64(speed))
}

// Stops the motor at the given port.
func Stop(port OutPort) {
	utilities.WriteIntValue(findFilename(port), "run", 0)
}

// Reads the operating speed of the motor at the given port.
func CurrentSpeed(port OutPort) int16 {
	return utilities.ReadInt16Value(findFilename(port), "duty_cycle_sp")
}

// Reads the operating power of the motor at the given port.
func CurrentPower(port OutPort) int16 {
	return utilities.ReadInt16Value(findFilename(port), "power")
}

// Enables regulation mode, causing the motor at the given port to compensate for any resistance and maintain its target speed.
func EnableRegulationMode(port OutPort) {
	utilities.WriteStringValue(findFilename(port), "regulation_mode", "on")
}

// Disables regulation mode. Regulation mode is off by default.
func DisableRegulationMode(port OutPort) {
	utilities.WriteStringValue(findFilename(port), "regulation_mode", "off")
}

// Enables brake mode, causing the motor at the given port to brake to stops.
func EnableBrakeMode(port OutPort) {
	utilities.WriteStringValue(findFilename(port), "stop_mode", "break")
}

// Disables brake mode, causing the motor at the given port to coast to stops. Brake mode is off by default.
func DisableBrakeMode(port OutPort) {
	utilities.WriteStringValue(findFilename(port), "stop_mode", "coast")
}
