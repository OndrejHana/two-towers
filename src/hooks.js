import "./types";
import { Mesh, Camera, Raycaster } from "three";
import { screenToNDC } from "./lib";

/**
 * @callback onSelectionChange
 * @param {Tower | null} old
 * @param {Tower | null} next
 */

/**
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

  return {
    selected,
    abort: controller.abort,
  };
}
