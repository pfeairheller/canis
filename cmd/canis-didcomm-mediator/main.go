/*
Copyright Scoir Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	arieslog "github.com/hyperledger/aries-framework-go/pkg/common/log"

	"github.com/scoir/canis/pkg/didcomm/mediator/cmd"
)

func main() {
	arieslog.SetLevel("aries-framework/out-of-band/service", arieslog.CRITICAL)
	arieslog.SetLevel("aries-framework/ws", arieslog.CRITICAL)
	arieslog.SetLevel("aries-framework/did-exchange/service", arieslog.DEBUG)
	arieslog.SetLevel("aries-framework/route/service", arieslog.DEBUG)

	cmd.Execute()
}
