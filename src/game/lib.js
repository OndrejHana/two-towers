import { COLORS, DIRECTIONS, TOWER, ROAD } from "./consts";

export function screenToNDC(x, y, width, height) {
  return new Vector2((x / width) * 2 - 1, -(y / height) * 2 + 1);
}

/**
 * @param {Point} point
 * @param {Point} world
 * @returns {Point[]}
 */
export function getAdjacentRoadCoords(point, world) {
  return DIRECTIONS.map((dir) => ({
    x: point.x + dir[0],
    y: point.y + dir[1],
  }))
    .filter(
      (coords) =>
        coords.x >= 0 &&
        coords.x < world.length &&
        coords.y >= 0 &&
        coords.y < world.length,
    )
    .filter((coords) => world[coords.x][coords.y].structure === ROAD);
}

/**
 * @param {Mesh} mesh
 * @param {Tile} tile
 */
export function setBaseTileColor(mesh, tile) {
  let color = COLORS.BASE_TILE;
  switch (tile.structure) {
    case TOWER:
      color = COLORS.TOWER_TILE;
      break;
    case ROAD:
      color = COLORS.ROAD_TILE;
      break;
  }

  mesh.material.color.set(color);
}
