// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

interface IMessenger {
    struct Message {
        uint256 nonce;
        uint256 sourceChainId;
        uint256 destinationChainId;
        address sender;
        bytes payload;
        uint256 timestamp;
    }
    
    event MessageSent(
        uint256 indexed nonce,
        uint256 indexed destinationChainId,
        address indexed sender,
        bytes payload,
        uint256 timestamp
    );
    
    event MessageReceived(
        bytes32 indexed messageHash,
        uint256 indexed sourceChainId,
        address indexed sender,
        bytes payload,
        uint256 nonce
    );
}