
Here is the corrected version of your text:

</br> <div align="center"> <a href="https://github.com/berachain/beacon-kit"> <picture> <source media="(prefers-color-scheme: dark)" srcset="https://res.cloudinary.com/duv0g402y/image/upload/v1718034312/BeaconKitBanner.png"> <img alt="beacon-kit-banner" src="https://res.cloudinary.com/duv0g402y/image/upload/v1718034312/BeaconKitBanner.png" width="auto" height="auto"> </picture> </a> </div> <h2 > A modular framework for building EVM consensus clients ⛵️✨ </h2>
The project is still heavily under construction. See the disclaimer below.

<div>


</div>
What is BeaconKit?
BeaconKit is a modular framework for building EVM-based consensus clients.
The framework offers the most user-friendly way to build and operate an EVM blockchain while ensuring a functionally identical execution environment to that of the Ethereum Mainnet.

Supported Execution Clients
Through utilizing the Ethereum Engine API, BeaconKit supports all six major Ethereum execution clients:

Geth: The official Go implementation of the Ethereum protocol.
Erigon: A more performant, feature-rich client forked from go-ethereum.
Nethermind: A .NET-based client with full support for Ethereum protocols.
Besu: An enterprise-grade client, Apache 2.0 licensed, and written in Java.
Reth: A Rust-based client focusing on performance and reliability.
Ethereumjs: A JavaScript-based client managed by the Ethereum Foundation.
Running a Local Development Network
Prerequisites:

Docker
Golang 1.23.0+
Foundry
Begin by opening two terminals side by side:

Terminal 1:

bash
Copy code
# Start the sample BeaconKit Consensus Client:
make start
Terminal 2:

bash
Copy code
# Start an Ethereum Execution Client:
make start-reth # or start-geth start-besu start-erigon start-nethermind start-ethereumjs
The account with
private-key=0xfffdbb37105441e14b0ee6330d855d8504ff39e705c3afa8f859ac9865f99306,
corresponding to address=0x20f33ce90a13a4b5e7697e3544c3083b8f8a51d4, is preloaded with the native EVM token.

Multinode Local Devnet
Please refer to the Kurtosis README for more information on how to run a multinode local devnet.

Status
This project is a work in progress and subject to frequent changes as we are still wiring up the final system. Audits on BeaconKit are currently ongoing. We don't recommend using BeaconKit in a production environment yet.

