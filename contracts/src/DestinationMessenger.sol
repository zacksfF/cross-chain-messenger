// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "./interfaces/IMessenger.sol";

contract DestinationMessenger is IMessenger {
    address public relayer;
    mapping(bytes32 => bool) public processedMessages;
    mapping(address => uint256) public receivedCount;
    
    error OnlyRelayer();
    error AlreadyProcessed();
    error InvalidSourceChain();
    
    modifier onlyRelayer() {
        if (msg.sender != relayer) revert OnlyRelayer();
        _;
    }
    
    constructor(address _relayer) {
        relayer = _relayer;
    }
    
    function receiveMessage(
        uint256 _nonce,
        uint256 _sourceChainId,
        address _sender,
        bytes calldata _payload,
        uint256 _timestamp
    ) external onlyRelayer returns (bytes32) {
        if (_sourceChainId == block.chainid) revert InvalidSourceChain();
        
        bytes32 messageHash = keccak256(
            abi.encodePacked(_nonce, _sourceChainId, block.chainid, _sender, _payload, _timestamp)
        );
        
        if (processedMessages[messageHash]) revert AlreadyProcessed();
        
        processedMessages[messageHash] = true;
        receivedCount[_sender]++;
        
        emit MessageReceived(messageHash, _sourceChainId, _sender, _payload, _nonce);
        
        return messageHash;
    }
    
    function updateRelayer(address _newRelayer) external onlyRelayer {
        relayer = _newRelayer;
    }
    
    function isProcessed(bytes32 _messageHash) external view returns (bool) {
        return processedMessages[_messageHash];
    }
}