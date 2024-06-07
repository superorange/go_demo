package picture

import (
	"fmt"
	"os"

	"github.com/davidbyttow/govips/v2/vips"
)

func checkError(err error) {
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}

func Run() {
	vips.Startup(nil)
	defer vips.Shutdown()

	image1, err := vips.NewImageFromFile("/Users/tian/projects/private/git_clone/super/c_demo/input.jpeg")
	checkError(err)

	err = image1.Resize(0.5, vips.KernelAuto)
	if err != nil {
		return
	}

	checkError(err)
	image1bytes, _, err := image1.ExportNative()
	err = os.WriteFile("output.jpg", image1bytes, 0644)
	checkError(err)

}
