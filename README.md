# Dungeon Crawler Game in Go

A simple roguelike dungeon crawler game written in Go.

## Features

- Procedurally generated dungeons with rooms and corridors
- Turn-based gameplay
- Player character with health, attack, defense stats
- Enemies with basic AI
- Items and inventory system
- Experience and leveling system
- Multiple dungeon levels

## How to Play

1. Compile and run the game:
   ```
   go build
   ./dungeon-game-golang
   ```

2. Controls:
   - Movement: w/a/s/d or up/down/left/right
   - Open inventory: i
   - Use stairs: > (when standing on them)
   - Rest to recover health: r
   - Help: h
   - Quit: q

## Game Elements

- **@**: Player character
- **#**: Wall (impassable)
- **.**: Floor (walkable)
- **+**: Door (can be opened)
- **$**: Treasure (collect for gold)
- **^**: Trap (causes damage)
- **>**: Stairs to next level
- **g/o/T/s/r**: Enemies (goblin, orc, troll, skeleton, rat)

## Combat

Move into enemies to attack them. Combat is turn-based - you attack first, then the enemy counterattacks if it survives.

## Development

This game is a simple demonstration of game development concepts in Go, including:
- Procedural content generation
- Game state management
- Turn-based mechanics
- Simple AI for enemies

Feel free to extend and modify the game as you like!