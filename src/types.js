/**
 * @typedef {number} StructureType
 * @enum {number}
 * @property {number} None - No structure (0)
 * @property {number} Tower - Tower structure (1)
 * @property {number} Road - Road structure (2)
 */

/**
 * @typedef {Object} Point
 * @property {number} x - X coordinate
 * @property {number} y - Y coordinate
 */

/**
 * @typedef {Object} Tile
 * @property {StructureType} structureType - Type of structure on the tile
 * @property {number|null} structureId - ID of the structure, if any
 */

/**
 * @typedef {Object} Tower
 * @property {number} id - Unique identifier of the tower
 * @property {string} playerId - ID of the player who owns the tower
 * @property {Point} point - Position of the tower
 */

/**
 * @typedef {Object} Unit
 * @property {number} id - Unique identifier of the unit
 * @property {string} playerId - ID of the player who owns the unit
 * @property {Point} point - Current position of the unit
 * @property {number} roadId - ID of the road the unit is following
 * @property {number} targetTowerId - ID of the tower the unit is targeting
 */

/**
 * @typedef {Object} Road
 * @property {number} id - Unique identifier of the road
 * @property {Point[]} points - Array of points defining the road's path
 * @property {Tower} from - Starting tower of the road
 * @property {Tower} to - Destination tower of the road
 */

/**
 * @typedef {Object} World
 * @property {Tile[][]} grid - 2D grid of tiles
 * @property {Road[]} roads - Array of all roads in the world
 */

/**
 * @typedef {Object} State
 * @property {Tower[]} towers - Array of all towers
 * @property {Unit[]} units - Array of all units
 */
