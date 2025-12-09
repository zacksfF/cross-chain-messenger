// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "./interfaces/IMessenger.sol";

contract SourceMessenger is IMessenger {
    uint256 public nonce;
    mapping(bytes32 => bool) public messageExists;
    
    error InvalidDestinationChain();
    error EmptyPayload();
    
    function sendMessage(uint256 _destChainId, bytes calldata _payload) external returns (bytes32) {
        if (_destChainId == block.chainid) revert InvalidDestinationChain();
        if (_payload.length == 0) revert EmptyPayload();
        
        bytes32 messageHash = keccak256(
            abi.encodePacked(nonce, block.chainid, _destChainId, msg.sender, _payload, block.timestamp)
        );
        
        messageExists[messageHash] = true;
        
        emit MessageSent(nonce, _destChainId, msg.sender, _payload, block.timestamp);
        
        nonce++;
        return messageHash;
    }
    
    function getMessageHash(
        uint256 _nonce,
        uint256 _sourceChainId,
        uint256 _destChainId,
        address _sender,
        bytes calldata _payload,
        uint256 _timestamp
    ) public pure returns (bytes32) {
        return keccak256(abi.encodePacked(_nonce, _sourceChainId, _destChainId, _sender, _payload, _timestamp));
    }
    
    function verifyMessage(bytes32 _messageHash) external view returns (bool) {
        return messageExists[_messageHash];
    }
}