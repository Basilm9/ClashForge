package clashgame

import (
	"fmt"
	"image/color"
)

func NewGame() *Game {
    // Create a new grid system first
    grid := NewGridSystem()
    
    // Load the tilemap
    err := grid.LoadTileMap("clashgame/csv/tilemap.csv")
    if err != nil {
        fmt.Println("Error loading tilemap:", err)
    }
    
    game := &Game{
        Players: [2]Player{
            NewPlayer(color.RGBA{255, 0, 0, 255}, true, 0, grid),
            NewPlayer(color.RGBA{0, 0, 255, 255}, false, 1, grid),
        },
        Grid: grid,
        BuildingMap: make(map[int]*Building),
        NextBuildingID: 1,
        ShowCSVPath: true, // Set to true to show CSV path
        CSVPath: "clashgame/csv/tilemap.csv", // Set the default CSV path
    }

    SetGameInstance(game)
    
    // Assign IDs to all buildings and add them to the map
    // Player 0 buildings
    game.Players[0].KingBuilding.Building.ID = game.NextBuildingID
    game.BuildingMap[game.NextBuildingID] = &game.Players[0].KingBuilding.Building
    game.NextBuildingID++
    
    for i := range game.Players[0].Buildings {
        game.Players[0].Buildings[i].ID = game.NextBuildingID
        game.BuildingMap[game.NextBuildingID] = &game.Players[0].Buildings[i]
        game.NextBuildingID++
    }
    
    // Player 1 buildings
    game.Players[1].KingBuilding.Building.ID = game.NextBuildingID
    game.BuildingMap[game.NextBuildingID] = &game.Players[1].KingBuilding.Building
    game.NextBuildingID++
    
    for i := range game.Players[1].Buildings {
        game.Players[1].Buildings[i].ID = game.NextBuildingID
        game.BuildingMap[game.NextBuildingID] = &game.Players[1].Buildings[i]
        game.NextBuildingID++
    }
    
    // Initialize the projectile system
    InitializeProjectileSystem()
    
    return game
}