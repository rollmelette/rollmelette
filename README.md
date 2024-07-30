# Rollmelette

![CI](https://github.com/gligneul/rollmelette/actions/workflows/ci.yaml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/gligneul/rollmelette)](https://goreportcard.com/report/github.com/gligneul/rollmelette)
[![GoDoc](https://img.shields.io/static/v1?label=godoc&message=reference&color=blue)](https://pkg.go.dev/github.com/rollmelette/rollmelette)
[![codecov](https://codecov.io/gh/gligneul/rollmelette/graph/badge.svg?token=467YCJV8PQ)](https://codecov.io/gh/gligneul/rollmelette)

![Logo](./logo.png)

Rollmelette is a high-level framework for Cartesi Rollups in Go.
The main objective of Rollmelette is to facilitate the development of the back-end of Cartesi applications.
It offers an abstraction on top of the Rollups API, manages assets from portals, and has unit-testing functions.

## Table of Contents

* [Getting Started](#getting-started)
* [The Application Interface](#the-application-interface)
* [Sending Outputs](#sending-outputs)
* [Managing assets](#managing-assets)
* [Unit Testing](#unit-testing)
* [Examples](#examples)

## Getting Started

> [!IMPORTANT]
> This README assumes you are familiar with the [Cartesi Rollups](https://docs.cartesi.io/cartesi-rollups/) and Go.

The recommended way of using Rollmelette is starting from the [project template][template].
Creating a new repository from the template and starting your own project is possible on the template GitHub page.
(Click the "Use this template" button.)
The template contains the skeleton of a Rollmelette application and the build files to run this application with [cartesi][cartesi].

### Pre-requisites

Before developing applications with Rollmelette, ensure you have the following dependencies installed.

- [Golang](https://go.dev/doc/install)
- [Docker](https://docs.docker.com/desktop/)
- [NoNodo](https://github.com/Calindra/nonodo)
- [Cartesi-CLI](https://docs.cartesi.io/cartesi-rollups/1.3/quickstart/)

### Sending advance-state inputs

After creating your new repository and cloning it to your machine, run `make dev` to run the application with NoNodo.
NoNodo allows you to run the application in the host machine for quick prototyping.
You should see the following output in your terminal.

```
Http Rollups for development started at http://localhost:5004
GraphQL running at http://localhost:8080/graphql
Inspect running at http://localhost:8080/inspect/
Press Ctrl+C to stop the node
```

Then, you can send an advance input with the binary payload `0xDEADBEEF` to the application with the following cartesi command.

```sh
cartesi send generic \
    --dapp=0x70ac08179605AF2D9e75782b8DEcDD3c22aA4D0C \
    --chain-id=31337 \
    --rpc-url=http://127.0.0.1:8545 \
    --mnemonic-passphrase='test test test test test test test test test test test junk' \
    --input=0xDEADBEEF
```

You should see the output below in the terminal running NoNodo.
This output means the Rollmelette application running inside NoNodo received the input.

```
[18:26:23.832] INF nonodo: added advance input index=0 sender=0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266 payload=0xdeadbeef
[18:26:23.906] INF nonodo: processing advance index=0
[18:26:23.908] INF command: log command=app buffer=stderr line="DBG received advance payload=0xdeadbeef inputIndex=0 msgSender=0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266 blockNumber=1 blockTimestamp=1705094783"
[18:26:23.908] INF nonodo: finished advance
```

### Sending inspect-state inputs

You can use `curl` to send an inspect input to the running application.
For instance, the command `curl http://localhost:8080/inspect/hello` sends the input with the payload `hello` to the application.
You should see the output below in the terminal running NoNodo.
Notice the application receives the input as a binary payload; hence, the string `hello` becomes the payload `0x68656c6c6f`.

```
[18:36:12.746] INF nonodo: added inspect input index=0 payload=0x68656c6c6f
[18:36:12.747] INF nonodo: processing inspect index=0
[18:36:12.752] INF command: log command=app buffer=stderr line="DBG received inspect payload=0x68656c6c6f"
[18:36:12.752] INF nonodo: finished inspect
```

### Building the application

The Rollmelette template contains a Dockerfile to build the RISC-V snapshot for the Cartesi Machine.
To build this snapshot, run `make build`.
Eventually, you will see the output below in the terminal, meaning cartesi built the application successfully, and it is ready to run.

```
         .
        / \
      /    \
\---/---\  /----\
 \       X       \
  \----/  \---/---\
       \    / CARTESI
        \ /   MACHINE
         '

[INFO  rollup_http_server] starting http dispatcher service...
[INFO  rollup_http_server::http_service] starting http dispatcher http service!
[INFO  actix_server::builder] starting 1 workers
[INFO  actix_server::server] Actix runtime found; starting in Actix runtime
[INFO  rollup_http_server::dapp_process] starting dapp: /var/opt/cartesi-app/app

Manual yield rx-accepted (0x100000000 data)
Cycles: 162546926
162546926: 8db686eb1b7a38b23dc33c7692440dd49eed77a902a74434b44d09639dbca17e
Storing machine: please wait
```

### Running in production mode

Once you built the Cartesi Machine snapshot, you can run the application with cartesi by running `make run`.
After some time, you should see the output below in your terminal.
You can now send advance-state and inspect-state inputs to the Rollmelette application inside the Cartesi machine.

```bash
Attaching to prompt-1, validator-1
validator-1  | 2024-07-30 20-40-15 info remote-cartesi-machine pid:119 ppid:72 Initializing server on localhost:0
prompt-1     | Anvil running at http://localhost:8545
prompt-1     | GraphQL running at http://localhost:8080/graphql
prompt-1     | Inspect running at http://localhost:8080/inspect/
prompt-1     | Explorer running at http://localhost:8080/explorer/
prompt-1     | Press Ctrl+C to stop the node
```

## The Application Interface

In the project template, the `application.go` file contains the definition of the Rollmelette application.
This file contains two parts: the application structure and the main function.
The application struct in the template file implements the [`Application`][roll.application] interface from Rollmelette.
The application interface requires a method for handling the advance-state input and another for handling the inspect-state input.
The main function sets the configuration options for Rollmelette, starts the application by calling [`Run`][roll.run], and logs the error before exiting.

### Advance

```go
func (a *MyApplication) Advance(
	env rollmelette.Env,
	metadata rollmelette.Metadata,
	deposit rollmelette.Deposit,
	payload []byte,
) error {
	// Handle advance input
	return nil
}
```

The snippet above shows the definition of the `Advance` method for the template application.
The application should use the [`Env`][roll.env] object to send outputs and manage assets.
[`Metadata`][roll.metadata] is a plain struct that contains the advance-input metadata, such as the message sender and the input index.
[`Deposit`][roll.deposit] is an interface that contains information about deposits.
(Handling deposits will be explained in the [Managing assets](#managing-assets) section.)
The `payload` parameter is the binary payload of the advance-state input.

### Inspect

```go
func (a *MyApplication) Inspect(env rollmelette.EnvInspector, payload []byte) error {
	// Handle inspect input
	return nil
}
```

The snippet above shows the definition of the `Inspect` method.
The [`EnvInspector`][roll.envinspector] is a subset of the `Env` interface;
it allows the application to send reports and inspect the assets managed by Rollmelette.
Like `Advance`, the `payload` parameter is the binary payload of the input.

### Main

```go
func main() {
	ctx := context.Background()
	opts := rollmelette.NewRunOpts()
	app := new(MyApplication)
	err := rollmelette.Run(ctx, opts, app)
	if err != nil {
		slog.Error("application error", "error", err)
	}
}
```

The snippet above contains the definition of the main function of the template application.
It will likely not be required to change this function.

## Sending Outputs

The Rollmelette application should use the `Env` and `EnvInspector` interfaces to send outputs.

### Reports

The report is a mechanism to expose the application log or a piece of diagnostic information.
The `Report` method receives a binary payload from both the `Inspect` and `Advance` methods.
The example below uses the `Report` method to log the received payload as a string.

```go
func (a *MyApplication) Inspect(env rollmelette.EnvInspector, payload []byte) error {
	msg := fmt.Sprintf("Received the payload %v", string(payload))
	env.Report([]byte(msg))
	return nil
}
```

### Notices

The notice is a mechanism to expose information about the application to the external world.
The `Notice` method receives a binary payload and can only be sent from the `Advance` method.
Unlike reports, the Rollups Node computes proofs for the notices to verify them on the base layer.
The example below encodes a JSON string and passes it to the `Notice` method.
(It is common to send notices as structured data.)

```go
func (a *MyApplication) Advance(
	env rollmelette.Env,
	metadata rollmelette.Metadata,
	deposit rollmelette.Deposit,
	payload []byte,
) error {
	noticeData := struct {
		Payload string `json:"payload"`
	}{
		Payload: hexutil.Encode(payload),
	}
	noticeJson, err := json.Marshal(noticeData)
	if err != nil {
		return err
	}
	env.Notice(noticeJson)
	return nil
}
```

### Vouchers

The voucher is a mechanism to call a contract in the base layer.
The `Voucher` method in the `Env` interface receives the destination contract and the payload, which should conform to the [Solidity-ABI][solabi] specification of the given contract.
Like notices, vouchers have proofs and can only be sent from the `Advance` method.

The example below sends input to another Cartesi application using the input box contract.
The code uses the function [`abi.JSON`][abi.json] from the [go-ethereum][go-ethereum] library to load the Solidity ABI of the input box contract.
Then, it uses the `Pack` method to encode the input according to the specification.
The code uses the function [`NewAddressBook`][roll.book] to create an address book and obtain the address of the input box.
Finally, the code calls the method `Voucher` from the `Env` interface to send the input to the input box.

```go
const INPUT_BOX_ABI = `[
  {
    "name": "addInput",
    "type": "function",
    "stateMutability": "nonpayable",
    "inputs": [
      {"type": "address", "internalType": "address", "components": null},
      {"type": "bytes", "internalType": "bytes", "components": null}
    ]
  }
]`

func (a *MyApplication) Advance(
	env rollmelette.Env,
	metadata rollmelette.Metadata,
	deposit rollmelette.Deposit,
	payload []byte,
) error {
	inputBox, err := abi.JSON(strings.NewReader(INPUT_BOX_ABI))
	if err != nil {
		return err
	}
	appAddress := common.HexToAddress("0xfafafafafafafafafafafafafafafafafafafafa")
	voucherPayload, err := inputBox.Pack("addInput", appAddress, payload)
	if err != nil {
		return err
	}
	addresses := rollmelette.NewAddressBook()
	env.Voucher(addresses.InputBox, voucherPayload)
	return nil
}
```

## Managing assets

Rollmelette provides built-in mechanisms to manage Ethereum assets.
When Rollmelette receives an input from a Cartesi portal, it parses it and registers the asset in an internal wallet.
Then, Rollmelette calls the `Advance` method, passing the corresponding deposit information to the `deposit` parameter.
Rollmelette passes the execution-layer data from the portal payload as the payload parameter.
The deposit parameter will be nil if the input does not come from a portal.

### Handling Deposits

The code below examplifies how a Rollmelette application can handle deposits.
It makes a type switch to know which kind of deposit the application received.
Then, the application should handle the deposit accordingly.

```go
func (a *MyApplication) Advance(
	env rollmelette.Env,
	metadata rollmelette.Metadata,
	deposit rollmelette.Deposit,
	payload []byte,
) error {
	if deposit == nil {
		// The input is not from a portal
	} else {
		switch deposit := deposit.(type) {
		case *rollmelette.EtherDeposit:
			// The input is from the Ether portal
		case *rollmelette.ERC20Deposit:
			// The input is from the ERC20 portal
		default:
			return fmt.Errorf("unsupported deposit: %T", deposit)
		}
	}
	return nil
}
```

### Ether

```go
type EtherDeposit struct {
	Sender common.Address
	Value *big.Int
}
```

The snippet above contains the definition of the Ether deposit.
The deposit contains the account that sent the Ether to the portal and the value sent in Wei.
Rollmelette stores these values in a wallet that maps accounts to the value deposited.
When an account sends another deposit to the application, Rollmelette adds up the value to the existing entry in the wallet.

Rollmelette offers functions in the `Env` interface to manipulate this wallet.
The functions are described in the table below.

| **Function** | **Description** |
|-|-|
| `env.EtherAddresses` | returns the list of addresses that have Ether. |
| `env.EtherBalanceOf` | returns the balance of the given address. |
| `env.EtherTransfer` | transfers the given amount of funds from source to destination. |
| `env.EtherWithdraw` | withdraws the asset from the wallet, generates the voucher to withdraw it from the application contract, and returns the voucher index. |

#### Application Address Relay

The application contract address is necessary to withdraw Ether.
To obtain this address, the application running inside the Cartesi rollups must receive an input from the application address relay.
This input is handled automatically by Rollmelette; the `Advance` method will not be called in this case.

### ERC20

```go
type ERC20Deposit struct {
	Token common.Address
	Sender common.Address
	Amount *big.Int
}
```

The snippet above contains the definition of the ERC20 deposit.
The deposit contains:

- The ERC20 token contract address.
- The account that sent the token to the portal.
- The amount of tokens sent.

Rollmelette stores the amount of tokens in a wallet that maps tokens to accounts to the amount deposited.
When an account sends another deposit to the application, Rollmelette adds up the amount of tokens to the existing entry in the wallet.

Rollmelette offers functions in the `Env` interface to manipulate this wallet.
The functions are described in the table below.

| **Function** | **Description** |
|-|-|
| `ERC20Tokens` | returns the list of tokens that have a non-zero balance in the application. |
| `ERC20Addresses` | returns the list of addresses that have the given token. |
| `ERC20BalanceOf` | returns the balance of the given address for the given token. |
| `ERC20Transfer` | transfers the given amount of tokens from source to destination. |
| `ERC20Withdraw` | withdraws the token from the wallet, generates the voucher to withdraw it from the ERC20 contract, and returns the voucher index. |

## Unit Testing

The Rollmelette template contains a unit test file called `application_test.go`.
You may execute the command `make test` to run the test.
The snippet below contains the expected output for this command.

```
go test -v ./...
=== RUN   TestMyApplicationSuite
=== RUN   TestMyApplicationSuite/TestAdvance
DBG received advance payload=0xdeadbeef inputIndex=0 msgSender=0xfafafafafafafafafafafafafafafafafafafafa blockNumber=0 blockTimestamp=1705176206
=== RUN   TestMyApplicationSuite/TestInspect
DBG received inspect payload=0xdeadbeef
--- PASS: TestMyApplicationSuite (0.00s)
    --- PASS: TestMyApplicationSuite/TestAdvance (0.00s)
    --- PASS: TestMyApplicationSuite/TestInspect (0.00s)
PASS
ok  	dapp	0.009s
```

The test file uses the [`Tester`][roll.tester] structure to send inputs to the application.
You may use the [`NewTester`][roll.newtester] function to create a tester struct, which receives the application struct that shall be tested.
To send advance-state inputs, the test code call the methods [`Advance`][roll.tester.advance], [`RelayAppAddress`][roll.tester.relayappaddress] [`DepositEther`][roll.tester.depositether], and [`DepositERC20`][roll.tester.depositerc20].
To send inspect-state inputs, the test code may call the [`Inspect`][roll.tester.inspect] method.
These methods call the application directly, collect the outputs, and return them for assertions.

## Examples

The Rollmelette repository contains some example applications under the `examples` directory.
The table below describes those examples.

| **Name** | **Description** |
|-|-|
| `address` | receives the app address in an advance input and returns it in an inspect input.|
| `echo` | emits a voucher, a notice, and a report for each advance input; and a report for each inspect input. |
| `error` | always returns an error; rejecting the input. |
| `honeypot` | is honeypot application that stores Ether. |
| `json` | simple game that uses JSON as inputs and outputs. |
| `panic` | always panic; Rollmelette captures the panic and reject the input. |

---

[NoNodo]: https://github.com/Calindra/nonodo
[Cartesi-CLI]: https://docs.cartesi.io/
[template]: https://github.com/gligneul/rollmelette-template
[solabi]: https://docs.soliditylang.org/en/latest/abi-spec.html
[go-ethereum]: https://geth.ethereum.org/docs/developers/dapp-developer/native

[abi.json]: https://pkg.go.dev/github.com/ethereum/go-ethereum/accounts/abi#JSON

[roll.application]: https://pkg.go.dev/github.com/rollmelette/rollmelette#Application
[roll.book]: https://pkg.go.dev/github.com/rollmelette/rollmelette#NewAddressBook
[roll.deposit]: https://pkg.go.dev/github.com/rollmelette/rollmelette#Deposit
[roll.env]: https://pkg.go.dev/github.com/rollmelette/rollmelette#Env
[roll.envinspector]: https://pkg.go.dev/github.com/rollmelette/rollmelette#EnvInspector
[roll.metadata]: https://pkg.go.dev/github.com/rollmelette/rollmelette#Metadata
[roll.newtester]: https://pkg.go.dev/github.com/rollmelette/rollmelette#NewTester
[roll.run]: https://pkg.go.dev/github.com/rollmelette/rollmelette#Run
[roll.tester.advance]: https://pkg.go.dev/github.com/rollmelette/rollmelette#Tester.Advance
[roll.tester.depositerc20]: https://pkg.go.dev/github.com/rollmelette/rollmelette#Tester.DepositERC20
[roll.tester.depositether]: https://pkg.go.dev/github.com/rollmelette/rollmelette#Tester.DepositEther
[roll.tester]: https://pkg.go.dev/github.com/rollmelette/rollmelette#Tester
[roll.tester.inspect]: https://pkg.go.dev/github.com/rollmelette/rollmelette#Tester.Inspect
[roll.tester.relayappaddress]: https://pkg.go.dev/github.com/rollmelette/rollmelette#Tester.RelayAppAddress
