hackyslack2 [![Build Status](https://circleci.com/gh/arkie/hackyslack2.svg?style=shield)](https://circleci.com/gh/arkie/hackyslack2)
-----------

A framework for appengine-powered Slack apps.

Currently used in:
- [Dicebot](https://dice-b.appspot.com) - Simple app to add a /roll command.

To use in your own apps:
- View the parent [dicebot](/) app which in turn uses the [dicebot/roll](/roll) dice parsing library.
- Import the repository.
  
  ```import "github.com/arkie/hackyslack2/slack"```

- Configure with the application id and secret from [Slack](https://api.slack.com/applications) ([dicebot/index.go](/index.go#L17))

  ```slack.Configure(clientId, clientSecret)```

- Register a command to run when the application is called. ([dicebot/roll.go](/roll.go#L14))

  ```slack.Register("roll", command)```
  
- Deploy the new app with ```goapp deploy```
  - Or use continuous deployment similar to dicebot's [CircleCI](/.circleci/config.yml) configuration, inspired by [yosukesuzuki/hugo-gae](https://github.com/yosukesuzuki/hugo-gae).
