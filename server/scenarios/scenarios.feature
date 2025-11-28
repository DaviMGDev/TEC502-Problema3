Feature: User Management
  Scenario: Login Successful
    Given the user is on the beginning of the chat application
    When the user enters valid credentials
  Scenario: Login Failed
    Given the user is on the beginning of the chat application
    When the user enters invalid credentials

Feature: Cards Management
  Scenario: Get Cards 
    Given the user is logged in
    When the user requests a card package

  Scenario: Exchange Cards
    Given the user has cards in their collection
    When the user exchanges cards with another user in the same room 

  Scenario: View Cards 
    Given the user is logged in
    When the user views their card collection

Feature: Chat System
  Scenario: Send Message
    Given the user is in a chat room
    When the user sends a message to the room 

Feature: Room Management
  Scenario: Create Room
    Given the user is logged in
    When the user creates a new chat room 

  Scenario: Join Room
    Given the user is logged in
    When the user joins an existing chat room 

  Scenario: Search Room
    Given the user is logged in
    When the user searches for a chat room by name

  Scenario: Leave Room
    Given the user is in a chat room
    When the user leaves the chat room 

Feature: Game System 
  Scenario: 


