/**
 * @typedef {object} Player
 * @property {string} color
 */

/**
 * @typedef {object} Point
 * @property {number} x - The x-coordinate of the tile.
 * @property {number} y - The y-coordinate of the tile.
 */

/**
 * @typedef {object} Tile
 * @property {number} structure - The structure type (e.g., 0 for none, 1 for something).
 * @property {number | null} structureId - The ID of the structure, or null if none.
 * @property {number | null} unitId - The ID of the unit on the tile, or null if none.
 */

/**
 * @typedef {object} Player
 * @property {string} color - The player's color (hex code).
 */

/**
 * @typedef {object} Tower
 * @property {number} playerId - The ID of the player who owns the tower.
 * @property {number | null} targetRoadId - The ID of the road the tower is targeting, or null if none.
 * @property {Point} point
 */

/**
 * @typedef {object} Road
 * @property {Point[]} points - The tiles that make up the road.
 * @property {Tower} from - The starting point of the road.
 * @property {Tower} to - The ending point of the road.
 */

/**
 * @typedef {object} RoadEnd
 * @property {number} playerId - The ID of the player associated with this end (if any).
 * @property {number | null} targetRoadId - The ID of the road this end is targeting, or null if none.
 * @property {number} x - The x-coordinate of the end.
 * @property {number} y - The y-coordinate of the end.
 */

/**
 * @typedef {object} Unit
 * @property {number} playerId - The ID of the player who owns the unit.
 * @property {Point} point
 */

/**
 * @typedef {object} Payload
 * @property {Tile[][]} world - The game world grid.
 * @property {Player[]} players - The list of players.
 * @property {Tower[]} towers - The list of towers.
 * @property {Road[]} roads - The list of roads.
 * @property {Unit[]} units - The list of units.
 */
