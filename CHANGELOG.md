# DwarfSweeper Changelog

## 0.5.20230127

### Features

* Added Mole Gnomes
    * Critters that when out of the ground, chase and leap at the dwarf
    * One point of health
    * Will burrow into diggable blocks when unable to get to the dwarf
    * Will then dig through blocks to find a spot to ambush the dwarf
* Bats now start flying if a block is destroyed near them
* New Zone: Glacier
* Added a Biome selector on the main menu
* Added a new object, a bomb dispenser: stand close to get a new bomb
    * Used for bosses to keep the game from softlocking
    * Has a five second timer between dropping bombs
* Added shaky cam back in

### Bugfixes

* Fixed a whole bunch of small issues with the Gnome Boss
    * The Boss would freeze randomly
    * Widened the bottom of the stairs
    * Fixed camera movement during "cutscene"
    * Made the boss be farther away from the player at the start
    * Gave the players more time to dodge the first charge
    * Fixed the exit not showing up correctly
    * The music now fades out when the boss is defeated
    * Added new emerge detection for the boss
    * Fixed sounds

## 0.4.20220729

### Features

* Improved Collision
* Massively improved Cave generation time
* Added ability to Drop Items (hold Interact button)
* Added Slugs and Bats back in (w/improvements!)
* Added Evil Bats
* Added Throwing Shovel item
* Added Pickaxe item
* Added Metal Detector item

### Bugfixes

* Blinking now works
* Fixed particles entering tiles and rapidly shooting off the top of the screen
* Score screen now correctly displays the cave behind it

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

### Bugfixes

* Fixed numbers' backgrounds
* Entities could be stuck in "ricochet" mode
* Interact button would Interact with all Interactables
* Some hitboxes were incorrectly calculating position
