import * as THREE from "three";
import WebGL from "three/addons/capabilities/WebGL.js";

import "./types";
import { COLORS } from "./consts";
import { setBaseTileColor, getAdjacentRoadCoords } from "./lib";
import { registerTowerRoadSelection, registerTowerSelection } from "./hooks";

function createScene(width, height) {
  const scene = new THREE.Scene();

  const frustumSize = 20;
  const aspect = width / height;
  const camera = new THREE.OrthographicCamera(
    (-frustumSize * aspect) / 2, // left
    (frustumSize * aspect) / 2, // right
    frustumSize / 2, // top
    -frustumSize / 2, // bottom
    0.1, // near
    1000, // far
  );

  camera.position.set(10, 10, 10);
  camera.lookAt(scene.position);

  const renderer = new THREE.WebGLRenderer({ antialias: true });
  renderer.setPixelRatio(window.devicePixelRatio * 1.5);
  renderer.setSize(width, height);

  return {
    scene,
    renderer,
    camera,
  };
}

/**
 * @param {THREE.WebGLRenderer} renderer
 * @param {HTMLDivElement} rootDiv
 * @param {Function} loop
 */
function renderCanvas(renderer, rootDiv, loop) {
  const container = document.createElement("div");
  rootDiv.appendChild(container);

  if (WebGL.isWebGL2Available()) {
    container.appendChild(renderer.domElement);
    renderer.setAnimationLoop(loop);
  } else {
    const warning = WebGL.getWebGL2ErrorMessage();
    container.appendChild(warning);
  }
}

/**
 * @param {Tile[][]} world
 * @param {THREE.Scene} scene
 * @param {number} tileSize
 * @param {THREE.ColorRepresentation} baseColor
 * @returns {THREE.Mesh[][]}
 */
function initGrid(world, scene, tileSize = 1) {
  const geometry = new THREE.BoxGeometry(tileSize, 0.1, tileSize);

  return world.map((row, x) => {
    const terrainRow = row.map((tile, z) => {
      const material = new THREE.MeshBasicMaterial();
      const mesh = new THREE.Mesh(geometry, material);
      mesh.position.set(x * tileSize, 0, z * tileSize);
      setBaseTileColor(mesh, tile);
      scene.add(mesh);
      return mesh;
    });
    return terrainRow;
  });
}

/**
 * @param {Tower[]} towers
 * @param {THREE.Scene} scene
 * @param {Player[]} players
 * @param {number} tileSize
 * @returns {Mesh[]}
 */
function initTowers(towers, scene, players, tileSize = 1) {
  const geometry = new THREE.CylinderGeometry(
    tileSize / 4,
    tileSize / 4,
    tileSize,
    8,
  );
  return towers.map((tower) => {
    const color =
      tower.playerId !== null
        ? players[tower.playerId].color
        : COLORS.NEUTRAL_TOWER;
    const material = new THREE.MeshBasicMaterial({ color });
    const mesh = new THREE.Mesh(geometry, material);
    mesh.position.set(tower.point.x, 0.5, tower.point.y);
    scene.add(mesh);
    return mesh;
  });
}

/**
 * @param {Unit[]} units
 * @param {Player[]} players
 * @param {THREE.Scene} scene
 */
function initUnits(units, players, scene, tileSize = 1) {
  const geometry = new THREE.CapsuleGeometry(tileSize / 4, tileSize / 4);
  units.map((unit) => {
    const mesh = new THREE.Mesh(
      geometry,
      new THREE.MeshBasicMaterial({ color: players[unit.playerId].color }),
    );
    mesh.position.set(unit.point.x, 0.5, unit.point.y);
    scene.add(mesh);
    return mesh;
  });
}

/**
 * @param {WebSocket} conn
 * @param {Payload} payload
 * @param {HTMLDivElement} appDiv
 */
function init(payload, conn, appDiv) {
  const { scene, renderer, camera } = createScene(
    window.innerWidth,
    window.innerHeight,
  );
  const grid = initGrid(payload.world, scene);
  const towers = initTowers(payload.towers, scene, payload.players);
  const units = initUnits(payload.units, payload.players, scene);

  let highlightedRoads = [];

  const { getSelected } = registerTowerSelection(
    towers,
    payload.towers,
    appDiv,
    camera,
    (old, next) => {
      if (old !== null) {
        getAdjacentRoadCoords(old.point, payload.world).forEach((coords) =>
          setBaseTileColor(
            grid[coords.x][coords.y],
            payload.world[coords.x][coords.y],
          ),
        );
      }
      if (next !== null) {
        highlightedRoads = getAdjacentRoadCoords(next.point, payload.world).map(
          (coords) => {
            const tile = grid[coords.x][coords.y];
            tile.material.color.set(
              next.targetRoadId ===
                payload.world[coords.x][coords.y].structureId
                ? COLORS.SELECTED_ACTIVE
                : COLORS.SELECTED_INACTIVE,
            );
            return tile;
          },
        );
      } else {
        highlightedRoads = [];
      }
    },
  );

  registerTowerRoadSelection(
    getSelected,
    payload.world,
    grid,
    conn,
    appDiv,
    camera,
  );

  renderCanvas(renderer, appDiv, () => {
    renderer.render(scene, camera);
  });
}

window.addEventListener("load", async function () {
  const res = await fetch("/game/new");
  const payload = await res.json();

  const conn = new WebSocket("ws://" + document.location.host + "/game/ws");

  const button = document.createElement("button");
  button.innerText = "click me";
  button.className = "ws-testbutton";
  const appDiv = document.querySelector("#app");
  appDiv.appendChild(button);

  conn.onclose = function (_) {
    console.log("connection closed");
  };
  conn.onmessage = function (evt) {
    var messages = evt.data.split("\n");
    console.log(messages);
  };

  button.addEventListener("mousedown", (_) => {
    console.log("clicked");
    conn.send("sup");
  });

  console.log(payload);
  init(payload, conn, appDiv);
});
