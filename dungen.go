package main

import (
	"fmt"
	"math/rand"
	"time"
)

// TileType represents different types of dungeon tiles
type TileType rune

const (
	Floor     TileType = '.'  // Empty floor space
	Wall      TileType = '#'  // Impassable wall
	Door      TileType = '+'  // Door (can be opened)
	Treasure  TileType = '$'  // Treasure (can be collected)
	Trap      TileType = '^'  // Trap (causes damage)
	StairsDown TileType = '>' // Stairs to next level
)

// Room represents a rectangular room in the dungeon
type Room struct {
	X, Y          int // Top-left corner
	Width, Height int
}

// Enemy represents a monster in the dungeon
type Enemy struct {
	X, Y    int
	Health  int
	Symbol  rune
	Name    string
	Damage  int
	Hostile bool
}

// Dungeon represents the game map as a 2D grid of runes (characters)
type Dungeon struct {
	Width, Height int       // Dimensions of the dungeon
	Grid          [][]rune  // 2D grid representing the dungeon layout
	Rooms         []Room    // List of rooms in the dungeon
	Enemies       []*Enemy  // List of enemies in the dungeon
	Items         []Item    // List of items in the dungeon
	Level         int       // Current dungeon level
}

// NewDungeon creates a new dungeon of width w and height h
func NewDungeon(w, h int) *Dungeon {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())
	
	// Create a new dungeon instance
	d := &Dungeon{
		Width:  w,
		Height: h,
		Level:  1,
	}
	
	// Initialize the grid with walls
	d.Grid = make([][]rune, h)
	for y := range d.Grid {
		d.Grid[y] = make([]rune, w)
		for x := range d.Grid[y] {
			d.Grid[y][x] = rune(Wall) // Initialize all cells as walls
		}
	}
	
	// Generate rooms and corridors
	d.generateRooms(4, 8) // Generate between 4-8 rooms
	d.connectRooms()      // Connect rooms with corridors
	d.addFeatures()       // Add doors, traps, treasures
	d.spawnEnemies(3, 6)  // Spawn 3-6 enemies
	
	return d
}

// generateRooms creates random rooms in the dungeon
func (d *Dungeon) generateRooms(minRooms, maxRooms int) {
	// Determine number of rooms to generate
	numRooms := minRooms + rand.Intn(maxRooms-minRooms+1)
	
	// Room size constraints
	minSize := 4
	maxSize := 10
	
	// Try to place rooms
	for i := 0; i < numRooms; i++ {
		// Random room dimensions
		width := minSize + rand.Intn(maxSize-minSize+1)
		height := minSize + rand.Intn(maxSize-minSize+1)
		
		// Random position (leaving border)
		x := 1 + rand.Intn(d.Width-width-2)
		y := 1 + rand.Intn(d.Height-height-2)
		
		// Create new room
		newRoom := Room{X: x, Y: y, Width: width, Height: height}
		
		// Check if room overlaps with existing rooms
		overlap := false
		for _, room := range d.Rooms {
			if roomsOverlap(newRoom, room) {
				overlap = true
				break
			}
		}
		
		// If no overlap, add the room
		if !overlap {
			d.carveRoom(newRoom)
			d.Rooms = append(d.Rooms, newRoom)
		}
	}
	
	// If no rooms were created, ensure at least one room exists
	if len(d.Rooms) == 0 {
		room := Room{X: d.Width/4, Y: d.Height/4, Width: d.Width/2, Height: d.Height/2}
		d.carveRoom(room)
		d.Rooms = append(d.Rooms, room)
	}
}

// roomsOverlap checks if two rooms overlap (including a 1-tile buffer)
func roomsOverlap(r1, r2 Room) bool {
	return r1.X-1 <= r2.X+r2.Width && r1.X+r1.Width+1 >= r2.X &&
		r1.Y-1 <= r2.Y+r2.Height && r1.Y+r1.Height+1 >= r2.Y
}

// carveRoom creates a room by setting floor tiles
func (d *Dungeon) carveRoom(room Room) {
	for y := room.Y; y < room.Y+room.Height; y++ {
		for x := room.X; x < room.X+room.Width; x++ {
			d.Grid[y][x] = rune(Floor)
		}
	}
}

// connectRooms creates corridors between rooms
func (d *Dungeon) connectRooms() {
	// Skip if there's only one or no rooms
	if len(d.Rooms) <= 1 {
		return
	}
	
	// Connect each room to the next one
	for i := 0; i < len(d.Rooms)-1; i++ {
		// Get center points of current and next room
		x1 := d.Rooms[i].X + d.Rooms[i].Width/2
		y1 := d.Rooms[i].Y + d.Rooms[i].Height/2
		x2 := d.Rooms[i+1].X + d.Rooms[i+1].Width/2
		y2 := d.Rooms[i+1].Y + d.Rooms[i+1].Height/2
		
		// Randomly decide whether to go horizontal first or vertical first
		if rand.Intn(2) == 0 {
			// Horizontal then vertical
			d.createHorizontalCorridor(x1, x2, y1)
			d.createVerticalCorridor(y1, y2, x2)
		} else {
			// Vertical then horizontal
			d.createVerticalCorridor(y1, y2, x1)
			d.createHorizontalCorridor(x1, x2, y2)
		}
	}
}

// createHorizontalCorridor creates a horizontal corridor
func (d *Dungeon) createHorizontalCorridor(x1, x2, y int) {
	// Ensure x1 is the smaller value
	if x1 > x2 {
		x1, x2 = x2, x1
	}
	
	// Create corridor
	for x := x1; x <= x2; x++ {
		if y >= 0 && y < d.Height && x >= 0 && x < d.Width {
			d.Grid[y][x] = rune(Floor)
		}
	}
}

// createVerticalCorridor creates a vertical corridor
func (d *Dungeon) createVerticalCorridor(y1, y2, x int) {
	// Ensure y1 is the smaller value
	if y1 > y2 {
		y1, y2 = y2, y1
	}
	
	// Create corridor
	for y := y1; y <= y2; y++ {
		if y >= 0 && y < d.Height && x >= 0 && x < d.Width {
			d.Grid[y][x] = rune(Floor)
		}
	}
}

// addFeatures adds doors, traps, and treasures to the dungeon
func (d *Dungeon) addFeatures() {
	// Add doors between corridors and rooms
	d.addDoors()
	
	// Add treasures in rooms
	d.addTreasures()
	
	// Add traps in corridors
	d.addTraps()
	
	// Add stairs to next level in the last room
	if len(d.Rooms) > 0 {
		lastRoom := d.Rooms[len(d.Rooms)-1]
		stairsX := lastRoom.X + lastRoom.Width/2
		stairsY := lastRoom.Y + lastRoom.Height/2
		d.Grid[stairsY][stairsX] = rune(StairsDown)
	}
}

// addDoors adds doors at appropriate locations
func (d *Dungeon) addDoors() {
	// For simplicity, we'll just add some random doors
	// A more sophisticated algorithm would place doors at corridor-room junctions
	for y := 1; y < d.Height-1; y++ {
		for x := 1; x < d.Width-1; x++ {
			// Check if this is a potential door location (floor with walls on opposite sides)
			if d.Grid[y][x] == rune(Floor) {
				if (d.Grid[y-1][x] == rune(Wall) && d.Grid[y+1][x] == rune(Wall)) ||
					(d.Grid[y][x-1] == rune(Wall) && d.Grid[y][x+1] == rune(Wall)) {
					// 10% chance to place a door
					if rand.Intn(100) < 10 {
						d.Grid[y][x] = rune(Door)
					}
				}
			}
		}
	}
}

// addTreasures adds treasure items to rooms
func (d *Dungeon) addTreasures() {
	// Add treasures to some rooms
	for _, room := range d.Rooms {
		// 40% chance for a room to have treasure
		if rand.Intn(100) < 40 {
			// Place treasure at random position in room
			treasureX := room.X + rand.Intn(room.Width)
			treasureY := room.Y + rand.Intn(room.Height)
			d.Grid[treasureY][treasureX] = rune(Treasure)
			
			// Add to items list
			d.Items = append(d.Items, Item{
				X:      treasureX,
				Y:      treasureY,
				Type:   ItemTreasure,
				Name:   "Gold",
				Value:  10 + rand.Intn(90), // 10-99 gold
				Symbol: '$',
			})
		}
	}
}

// addTraps adds dangerous traps to the dungeon
func (d *Dungeon) addTraps() {
	// Add some traps in corridors and rooms
	numTraps := 2 + rand.Intn(4) // 2-5 traps
	
	for i := 0; i < numTraps; i++ {
		// Try to place a trap
		for attempts := 0; attempts < 50; attempts++ {
			x := 1 + rand.Intn(d.Width-2)
			y := 1 + rand.Intn(d.Height-2)
			
			// Only place traps on floor tiles
			if d.Grid[y][x] == rune(Floor) {
				d.Grid[y][x] = rune(Trap)
				break
			}
		}
	}
}

// spawnEnemies creates enemies in the dungeon
func (d *Dungeon) spawnEnemies(min, max int) {
	numEnemies := min + rand.Intn(max-min+1)
	
	enemyTypes := []struct {
		name   string
		symbol rune
		health int
		damage int
	}{
		{"Goblin", 'g', 3, 1},
		{"Orc", 'o', 5, 2},
		{"Troll", 'T', 8, 3},
		{"Rat", 'r', 1, 1},
		{"Skeleton", 's', 4, 2},
	}
	
	// Spawn enemies in rooms (not the first room, which is the player's starting point)
	for i := 0; i < numEnemies; i++ {
		if len(d.Rooms) <= 1 {
			break
		}
		
		// Choose a random room (not the first one)
		roomIndex := 1 + rand.Intn(len(d.Rooms)-1)
		room := d.Rooms[roomIndex]
		
		// Choose a random position in the room
		x := room.X + rand.Intn(room.Width)
		y := room.Y + rand.Intn(room.Height)
		
		// Choose a random enemy type
		enemyType := enemyTypes[rand.Intn(len(enemyTypes))]
		
		// Create the enemy
		enemy := &Enemy{
			X:       x,
			Y:       y,
			Health:  enemyType.health,
			Symbol:  enemyType.symbol,
			Name:    enemyType.name,
			Damage:  enemyType.damage,
			Hostile: true,
		}
		
		// Add to enemies list
		d.Enemies = append(d.Enemies, enemy)
	}
}

// IsWalkable checks whether the (x, y) position is within bounds and walkable
func (d *Dungeon) IsWalkable(x, y int) bool {
	// Check bounds
	if x < 0 || y < 0 || x >= d.Width || y >= d.Height {
		return false // Out of bounds
	}
	
	// Check tile type
	tile := TileType(d.Grid[y][x])
	switch tile {
	case Floor, Door, Treasure, Trap, StairsDown:
		return true // These tiles are walkable
	default:
		return false // Walls and other tiles are not walkable
	}
}

// GetTileAt returns the tile type at the given coordinates
func (d *Dungeon) GetTileAt(x, y int) TileType {
	if x < 0 || y < 0 || x >= d.Width || y >= d.Height {
		return Wall // Out of bounds is treated as wall
	}
	return TileType(d.Grid[y][x])
}

// GetEnemyAt returns the enemy at the given coordinates, or nil if none
func (d *Dungeon) GetEnemyAt(x, y int) *Enemy {
	for _, enemy := range d.Enemies {
		if enemy.X == x && enemy.Y == y && enemy.Health > 0 {
			return enemy
		}
	}
	return nil
}

// GetItemAt returns the item at the given coordinates, or nil if none
func (d *Dungeon) GetItemAt(x, y int) *Item {
	for i, item := range d.Items {
		if item.X == x && item.Y == y && !item.Collected {
			return &d.Items[i]
		}
	}
	return nil
}

// RemoveEnemy removes a dead enemy from the dungeon
func (d *Dungeon) RemoveEnemy(enemy *Enemy) {
	for i, e := range d.Enemies {
		if e == enemy {
			// Remove from slice
			d.Enemies = append(d.Enemies[:i], d.Enemies[i+1:]...)
			break
		}
	}
}

// MoveEnemies updates enemy positions based on simple AI
func (d *Dungeon) MoveEnemies(player *Player) {
	for _, enemy := range d.Enemies {
		// Skip dead enemies
		if enemy.Health <= 0 {
			continue
		}
		
		// Simple AI: Move randomly, but prefer moving toward player if nearby
		dx, dy := 0, 0
		
		// Calculate distance to player
		distX := player.X - enemy.X
		distY := player.Y - enemy.Y
		distance := abs(distX) + abs(distY) // Manhattan distance
		
		// If player is close (within 5 tiles), move toward them
		if distance < 5 && enemy.Hostile {
			// Move in the direction of the player
			if abs(distX) > abs(distY) {
				// Move horizontally
				if distX > 0 {
					dx = 1
				} else {
					dx = -1
				}
			} else {
				// Move vertically
				if distY > 0 {
					dy = 1
				} else {
					dy = -1
				}
			}
		} else {
			// Move randomly
			if rand.Intn(3) > 0 { // 2/3 chance to move
				directions := []struct{ dx, dy int }{
					{0, -1}, {1, 0}, {0, 1}, {-1, 0}, // Up, right, down, left
				}
				dir := directions[rand.Intn(len(directions))]
				dx, dy = dir.dx, dir.dy
			}
		}
		
		// Check if the new position is valid
		newX, newY := enemy.X+dx, enemy.Y+dy
		
		// Don't move onto the player
		if newX == player.X && newY == player.Y {
			continue
		}
		
		// Check if the new position is walkable
		if d.IsWalkable(newX, newY) && d.GetEnemyAt(newX, newY) == nil {
			enemy.X, enemy.Y = newX, newY
		}
	}
}

// abs returns the absolute value of x
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// Print renders the dungeon grid, displaying the player, enemies, and items
func (d *Dungeon) Print(p *Player) {
	// Print the dungeon level
	fmt.Printf("Dungeon Level: %d\n", d.Level)
	
	// Print the grid
	for y := 0; y < d.Height; y++ {
		for x := 0; x < d.Width; x++ {
			// Check if there's an enemy at this position
			enemy := d.GetEnemyAt(x, y)
			if enemy != nil {
				fmt.Print(string(enemy.Symbol))
				continue
			}
			
			// Check if player is at this position
			if p.X == x && p.Y == y {
				fmt.Print("@") // Player's position
				continue
			}
			
			// Otherwise print the terrain
			fmt.Print(string(d.Grid[y][x]))
		}
		fmt.Println()
	}
}
