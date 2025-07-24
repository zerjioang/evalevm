## How to use this repository

This repository is used to evaluate many different EVM security tools at once saving engineers and researcher many hours. First you need to build all the docker images on you device. To do so, run:

```go
go build && ./evalevm analyzer build --tools ./tools
```

> Make sure to point to ./tools directory of this same repository as the code will be used to create the images.
## Analyzers out of scope

The following analyzers and research tools are not considered because they require either Solidity, Vyper source code or ABI files definitions. Our initial selection of tools focus only on those that work with the bytecode only.

* SmartCheck: this tool needs the contract source code to work.

## Vulnerability types and source code examples

You can use the following examples as input for tools evaluation.

### Hardcoded address detection

```solidity
pragma solidity 0.4.24;

    contract C {
        function f(uint a, uint b) pure returns (address) {
            address public multisig = 0xf64B584972FE6055a770477670208d737Fff282f;
            return multisig;
        }
    }
```

## Exact ETH equality

```solidity
pragma solidity 0.4.24;

    contract C {
        function valid pure returns (bool) {
            return address(this).balance == 42 ether
        }
    }
```

## Division before multiplication

```solidity
pragma solidity 0.4.25;

contract MyContract {

    uint constant BONUS = 500;
    uint constant DELIMITER = 10000;

    function calculateBonus(uint amount) returns (uint) {
        return amount/DELIMITER*BONUS;
    }
}
```

## Time equality

```solidity
pragma solidity 0.4.25;

contract Game {

    function oddOrEven(bool yourGuess) external payable {
        if (yourGuess == now % 2 > 0) {
            uint fee = msg.value / 10;
            msg.sender.transfer(msg.value * 2 - fee);
        }
    }

    function () external payable {}
}
```

## Current Block Hash usage

In EVM, the current block.hash is always zero.

<code>blockhash</code> function returns a non-zero value only for 256 last blocks. Besides, it always returns 0 for the current block, i.e. <code>blockhash(block.number)</code> always equals to 0.

```solidity
pragma solidity 0.8.16;

contract C {
    function currentBlockHash() public view returns(bytes32) {
        return blockhash(block.number); // 0
    }
}
```

## Contract that can lock Ether

In the following example, contracts programmed to receive ether does not call <code>transfer</code>, <code>send</code>, or <code>call.value</code> function

```solidity
pragma solidity 0.4.25;

contract BadMarketPlace {
    function deposit() payable {
        require(msg.value > 0);
    }
}
```


## Infinite loop

```solidity
pragma solidity 0.4.24;

contract GreaterOrEqualToZero {

    function infiniteLoop(uint border) returns(uint ans) {

        for (uint i = border; i >= 0; i--) {
            ans += i;
        }
    }
}
```

In this case, <code>i >= 0</code> condition will always evaluate to true. The next value of <code>i</code> variable after <code>0</code> will be <code>2**256-1</code>. Thus, the loop will be infinite.

```solidity
contract Malicious {
    function foo(uint) public {
        while (true) {} // infinite loop
    }
}
```

## Transfer in a loop

```solidity
pragma solidity 0.4.25;

contract MyContract {

    address[] public users;
    uint internal id;

    function transferBatch(address[] users) public {
        uint amount = address(this).balance / users.length;
        for (uint i = 0; i < users.length; i++) {
            users[i].transfer(amount);
        }
    }

    function () public payable {
        users[id] = msg.sender;
        id++;
    }
}
```

## Unchecked call

```solidity
contract SolidityUncheckedSend {
    function unseatKing(address addr, uint value) {
        addr.send(value);
    }
}
```

## Troubleshooting

How to check if there is an existing version of perf installed in the container

```bash
ls /usr/lib/linux-tools/*/perf
```

## References

* https://github.com/hzysvilla/Academic_Smart_Contract_Papers
* https://github.com/smartbugs/smartbugs
