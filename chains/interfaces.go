// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package chains

import (
	"github.com/rjman-ljm/platdot-utils/msg"
)

type Router interface {
	Send(message msg.Message) error
}
