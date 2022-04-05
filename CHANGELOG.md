# DwarfSweeper Changelog

## 0.4.20220401

### Features

* Added Splitscreen Local multiplayer
    * Each players' controls can be mapped in the input menu
* Added Quests
    * Pause menu has a new sub menu that is a list of all active and completed quests
    * Current quests include:
        * Several in a series for marking/flagging bombs
        * A quest for discovering each zone
* New Zone: Mossy Caverns
* New Zone: Crystal Hoard
* Big Bombs now have an explosion, including sound
* Added High Yield bombs (5x5 explosion)
* A dwarf can make 2 attacks/digs in the air
* The key to exit a level is now "interact"
* Added background parallax (missing for Crystal)
* Improved Puzzle controls
* Decreased zoom
* Implemented leading camera

### Bugfixes

* Dig blocks no longer cascade collapse blocks. Only collapse blocks cascade
* Removed extra tiles above entrance and exit of Minesweeper Levels
* Improved tile backgrounds
* Key symbol missing from Exit pop up is now there

## 0.2.20220222

### Features

* Finished first draft of the Gnome Mole Boss
* Basic 6 cave Descent complete
* Implemented new Typeface
* Created Cave Generator (for testing cave layouts)
* Changed base cave generation to use cellular automata
* Added Big Bombs that can be disarmed
* Added a Minesweeper Puzzle that must be completed to disarm the Big Bomb
* Text can now add symbols (sprites)
* Added a loading screen
* Added short delay for collectibles before the Dwarf picks them up (Gems, Apples)

### Bugs

* Fixed numbers' backgrounds
* Entities could be stuck in "ricochet" mode
* Interact button would Interact with all Interactables
* Some hitboxes were incorrectly calculating position
