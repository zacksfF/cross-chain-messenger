// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "forge-std/Test.sol";
import "../src/SourceMessenger.sol";

contract SourceMessengerTest is Test {
    SourceMessenger public messenger;
    address public user = address(0x123);
    
    function setUp() public {
        messenger = new SourceMessenger();
    }
    
    function testSendMessage() public {
        vm.prank(user);
        bytes memory payload = "Hello Mumbai";
        
        bytes32 hash = messenger.sendMessage(80001, payload);
        
        assertEq(messenger.nonce(), 1);
        assertTrue(messenger.messageExists(hash));
    }
    
    function testCannotSendToSameChain() public {
        vm.prank(user);
        vm.expectRevert(SourceMessenger.InvalidDestinationChain.selector);
        messenger.sendMessage(block.chainid, "test");
    }
    
    function testCannotSendEmptyPayload() public {
        vm.prank(user);
        vm.expectRevert(SourceMessenger.EmptyPayload.selector);
        messenger.sendMessage(80001, "");
    }
}