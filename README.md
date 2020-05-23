DiceBot
-------

[DiceBot](https://slack.com/apps/A0K7BALJ1-dicebot) is a [Slack](https://slack.com/) app.

 - [roll](/roll) -
   The dice parsing logic.
 - [slack](/slack) -
   The Slack specific authentication and parsing.

Installing
----------

Click `Add to Slack` on [this page](https://dice-b.appspot.com/) to install into your own Slack.

Usage
-----

The DiceBot app does a pretty forgiving parsing of the text passed to the `/roll` command,
to figure out what sort of dice to roll and how to use them.

First, it scans the string for all the "roll expressions" it can find.
Each roll expression consists of:

1. An optional operation, either:
	* `+` for addition
	* `-` for subtraction
	* `*` or `×` for multiplication
	* `/` for division
	* `^` for maximum
	* `v` for minimum

	If there is no operation, addition is assumed.
	(So, `1d20 1d20` will roll two d20s and add the results.)

	There is no "order of operations" or grouping of roll expressions;
	each expression's operation applies its value to the result of all the preceding expressions.
	For example, `1d4 * 1d6 ^ 1d8 + 1d10` will roll 1d4,
	then roll 1d6 and multiply the previous result by that,
	then roll 1d8 and take either that value or the previous result, whichever is higher,
	then roll 1d10 and add it to the previous result.
	If the rolls were `1 5 3 8`, the result would be `(((1 * 5) ^ 3) + 8)`, totalling 13.
2. A "dice expression", defined below.
3. Optionally, a "for" value, like `2d6 for damage`,
	which will print out in the results to help explain what was rolled.
	(The "for" value will consume all the text up to the next comma or semicolon.)

Any text other than the roll expressions is ignored,
and won't have any effect on the result.

A "dice expression" instructs the bot how many and what kind of dice to roll,
and what additional special actions to take.
There are a few possible forms:

1. Just a number, by itself, between 1 and 99999.
2. Fudge dice, like `5f`, which rolls 5 fudge dice (between 1 and 99999).
3. A full [standard notation dice expression](https://en.wikipedia.org/wiki/Dice_notation#Standard_notation), with the following pieces (in order):
	* Optionally, the number of dice to roll, between 1 and 999. If omitted, defaults to `1`.
	* The type of dice to roll: either `df` for fudge dice, `d%` for percentile dice, or `d6` (for any number between 1 and 1000) for dice with that many sides.
	* Optionally, a `!` to indicate that the dice "explode", rolling an additional die of the same size when the max value is rolled.
	* Optionally, a `<5` or `>5` (for any number between 2 and the die size - 1) to indicate that you want to count up how many dice rolled below or above the value.
	* Optionally, a `k5` or `k-5` (for any number between -999 and 999) to indicate that only the best K (for positive values) or worst K (for negative values) rolls should be counted in the result.

If no dice expression is provided, the DiceBot app will default to `2d6` or two six-sided dice similar to `//roll` from AOL Instant Messenger.

Examples
--------

D&D rolls can be written just like they are in text: `1d20 + 4`, `8d6`, etc. To roll with advantage or disadvantage, use the `k` modifier, like `2d20k1 + 2` for advantage, or `2d20k-1 + 2` for disadvantage.  The "standard" stat-rolling method can be written as `4d6k3 for Str`, etc., rolling 4 dice and dropping the lowest.

World of Darkness rolls should be written like `5d10!>7`, which will roll 5 exploding d10s, and count how many exceed 7 (are 8 or higher).

"For" values can make a result easier to read, like `1d8+1 for sword, +1d6 for fire` which will print out something like "@dm rolled *5* + *1* for *sword* + *1* for *fire* = *7*". (Note the comma after the first "for" value - it'll parse wrong if you leave that off.)
