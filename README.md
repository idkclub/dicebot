dicebot
-------

![Preview of roll command](/preview@2x.png)

An appengine-powered dice rolling Slack service. Aside from [roll.go](/roll.go), the rest of the code is designed to be an extensible and reusable base for Slack bots.

Currently installs:
- /roll - Roll X Y-sided dice

Currently powers:
- [Dicebot](https://dice-b.appspot.com) - /roll

To launch your own copy:
- Clone this repo.
- Run [./secret.sh](/secret.sh) to create secrets.go.
