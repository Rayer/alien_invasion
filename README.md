# Alien Invasion concept game

Mad aliens are about to invade the earth and you are tasked with simulating the invasion.

You are given a map containing the names of cities in the non-existent world of X. The map is in a file, with one city per line. The city name is first, followed by 1-4 directions (north, south, east, or west). Each one represents a road to another city that lies in that direction.

## Overview

![](https://github.com/Rayer/alien_invasion/blob/master/Overview.jpg)

## Roles
1. Alien will enter a random city
2. Alien will try to enter an adjacent city
3. When 2 aliens enters same city, they will fight, and result the city being destroyed.
4. If city is destroyed, all path lead to, and leads from this city, will be removed, preventing other aliens from entering or exiting.
5. After any alien goes 10000 steps, game will conclude, and dump the map file with remaining cities.

## Commandline Example

This release ships with a sample map file, and a sample map file with error :
- test_resources/sample_map.txt
- test_resources/sample_map_error.txt

And for `sample_map.txt`, you can take a look upon `sample_map.png` in same directory to know it's actual layout for each cities :
![](https://github.com/Rayer/alien_invasion/blob/master/test_resources/sample_map.png)

We can start the game by running the following command:

``` 
cd cmd
go build -o alien_invasion
./alien_invasion ../test_resources/sample_map.txt 5
```

and see the output.

## Development

Branch `develop` is the current development branch, and will be merged to `master` when ready.
Most description in code itself.

## Issue Tracking
TBD

## History
### 1.0.0
Initial Release
