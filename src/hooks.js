import "./types";
import { Raycaster, Mesh, Camera } from "three";
import { getAdjacentRoadCoords } from "./lib";
import { screenToNDC } from "./lib";

/**
 * @callback onSelectionChange
 * @param {Tower | null} old
 * @param {Tower | null} next
 */

/**
 *  Registers a hook, that selects a tower when clicked on.
 *
 * @param {Mesh[]} towerMeshes
 * @param {onSelectionChange} onSelectionChange
 * @param {Tower[]} towers
 * @param {HTMLDivElement} container
 * @param {Camera} camera
 */
export function registerTowerSelection(
  towerMeshes,
  towers,
  container,
  camera,
  onSelectionChange,
  raycaster = new Raycaster(),
  controller = new AbortController(),
) {
  let selected = null;
  const { width, height } = container.getBoundingClientRect();
  container.addEventListener(
    "mousedown",
    (e) => {
      raycaster.setFromCamera(
        screenToNDC(e.clientX, e.clientY, width, height),
        camera,
      );
      const intersections = raycaster.intersectObjects(towerMeshes);
      if (intersections.length > 0) {
        const intersection = intersections[0];
        const x = intersection.object.position.x;
        const y = intersection.object.position.z;
        const old = selected;
        const next = towers.filter(
          (tower) => tower.point.x === x && tower.point.y === y,
        )[0];
        if (old !== next) {
          selected = next;
          onSelectionChange(old, next);
        }
      }
    },
    {
      signal: controller.signal,
    },
  );

  document.addEventListener(
    "keydown",
    function (event) {
      if (event.key === "Escape") {
        const old = selected;
        selected = null;
        onSelectionChange(old, selected);
      }
    },
    { signal: controller.signal },
  );

  return {
    getSelected: () => selected,
    abort: controller.abort,
  };
}

/**
 * @callback getSelected
 * @returns {Tower | null}
 */

/**
 *  Registers a hook, that selects a tower when clicked on.
 *
 * @param {Tower | null} selected
 * @param {Mesh[]} highlightedRoads
 * @param {HTMLDivElement} container
 * @param {Camera} camera
 * @param {Mesh[][]} grid
 * @param {WebSocket} conn
 * @param {Tile[][]} world
 * @param {getSelected} getSelected
 */
export function registerTowerRoadSelection(
  getSelected,
  world,
  grid,
  conn,
  container,
  camera,
  raycaster = new Raycaster(),
  controller = new AbortController(),
) {
  const { width, height } = container.getBoundingClientRect();
  container.addEventListener(
    "mousedown",
    function (e) {
      const selected = getSelected();

      if (selected === null) {
        console.log("nothing selected");
        return;
      }

      const highlightedRoads = getAdjacentRoadCoords(selected.point, world).map(
        (coords) => grid[coords.x][coords.y],
      );

      raycaster.setFromCamera(
        screenToNDC(e.clientX, e.clientY, width, height),
        camera,
      );

      const intersections = raycaster.intersectObjects(highlightedRoads);

      if (intersections.length === 0) {
        console.log("no intersections");
        return;
      }

      const clicked = intersections[0].object.position;
      const clickedGridCoords = { x: clicked.x, y: clicked.z };

      const selectedTowerId =
        world[selected.point.x][selected.point.y].structureId;
      const clickedRoadId =
        world[clickedGridCoords.x][clickedGridCoords.y].structureId;

      if (clickedRoadId !== selected.targetRoadId) {
        conn.send(
          JSON.stringify({
            towerId: selectedTowerId,
            roadId: clickedRoadId,
          }),
        );
      }
    },
    {
      signal: controller.signal,
    },
  );
}
