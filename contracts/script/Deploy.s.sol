// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "forge-std/Script.sol";
import "../src/SourceMessenger.sol";
import "../src/DestinationMessenger.sol";

contract DeployScript is Script {
    function run() external {
        uint256 deployerPrivateKey = vm.envUint("PRIVATE_KEY");
        address relayerAddress = vm.envAddress("RELAYER_ADDRESS");
        
        vm.startBroadcast(deployerPrivateKey);
        
        // Deploy based on chain
        uint256 chainId = block.chainid;
        
        if (chainId == 11155111) { // Sepolia
            SourceMessenger source = new SourceMessenger();
            console.log("SourceMessenger deployed on Sepolia:", address(source));
            
            DestinationMessenger dest = new DestinationMessenger(relayerAddress);
            console.log("DestinationMessenger deployed on Sepolia:", address(dest));
        } else if (chainId == 80001) { // Mumbai
            SourceMessenger source = new SourceMessenger();
            console.log("SourceMessenger deployed on Mumbai:", address(source));
            
            DestinationMessenger dest = new DestinationMessenger(relayerAddress);
            console.log("DestinationMessenger deployed on Mumbai:", address(dest));
        }
        
        vm.stopBroadcast();
    }
}