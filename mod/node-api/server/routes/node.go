// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
// Use of this software is governed by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package routes

import (
	"github.com/berachain/beacon-kit/mod/node-api/server/handlers"
	"github.com/labstack/echo/v4"
)

func aasignNodeRoutes[ValidatorT any](
	e *echo.Echo,
	h handlers.RouteHandlers[ValidatorT],
) {
	e.GET("/eth/v1/node/identity",
		h.NotImplemented)
	e.GET("/eth/v1/node/peers",
		h.NotImplemented)
	e.GET("/eth/v1/node/peers/:peer_id",
		h.NotImplemented)
	e.GET("/eth/v1/node/peers/peer_count",
		h.NotImplemented)
	e.GET("/eth/v1/node/version",
		h.NotImplemented)
	e.GET("/eth/v1/node/syncing",
		h.NotImplemented)
	e.GET("/eth/v1/node/health",
		h.NotImplemented)
}