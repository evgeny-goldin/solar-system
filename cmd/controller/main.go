package main

import (
	"knative.dev/pkg/injection/sharedmain"
	"solar-system/pkg/reconciler/solar"
)

func main() {
	sharedmain.Main("solar_system_controller", solar.NewController)
}
