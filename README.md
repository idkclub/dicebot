hackyslack2
-----------

A framework for appengine-powered Slack apps.

Currently powers:
- [Dicebot](https://dice-b.appspot.com) - Simple app to add a /roll command.

To use in your own apps:
- View the included [dicebot](/dicebot) sample app.
- Import the repository.
  
  ```import "github.com/arkie/hackyslack2"```

- Configure with the application id and secret from [Slack](https://api.slack.com/applications) ([dicebot/index.go](/dicebot/index.go#L14))

  ```hackyslack.Configure(clientId, clientSecret)```

- Register a command to run when the application is called. ([dicebot/roll.go](/dicebot/roll.go#L14))

  ```hackyslack.Register("roll", command)```
  
- Deploy the new app with ```goapp deploy```
