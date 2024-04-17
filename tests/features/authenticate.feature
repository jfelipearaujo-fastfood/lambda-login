Feature: authenticate
    In order to be authenticated
    As an user
    I want to be able to authenticate with my CPF and password

    Scenario: authenticate with valid credentials
        Given the user CPF is "548.644.620-97"
        And the user password is "12345678"
        When the user request to be authenticated
        Then the user should be authenticated successfully