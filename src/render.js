import * as THREE from "three";
import WebGL from "three/addons/capabilities/WebGL.js";

const NONE = 0;
const TOWER = 1;
const ROAD = 2;

const DIRECTIONS = [
  [1, 0],
  [-1, 0],
  [0, 1],
  [0, -1],
  [1, 1],
  [1, -1],
  [-1, 1],
  [-1, -1],
];
/**
 * @param {WebSocket} conn
 */
export function render(payload, conn) {
  const app = document.querySelector("#app");
  const container = document.createElement("div");
  app.appendChild(container);
  const scene = new THREE.Scene();

  const frustumSize = 20;
  const aspect = window.innerWidth / window.innerHeight;
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
  renderer.setSize(window.innerWidth, window.innerHeight);

  function createTower(x, y, tile) {
    const tower = payload.towers[tile.structureId];
    const color =
      tower.playerId !== null
        ? payload.players[tower.playerId].color
        : "#64748b";

    const geometry = new THREE.CylinderGeometry(0.25, 0.25, 1, 8);
    const material = new THREE.MeshBasicMaterial({ color });
    const cylinder = new THREE.Mesh(geometry, material);
    cylinder.position.set(x, 0.5, y);

    scene.add(cylinder);
  }

  /**
   * @param {THREE.Mesh} mesh
   */
  function setColor(mesh, color) {
    console.log("setting color", color);
    mesh.material.color.set(color);
  }

  function createUnit(x, y, tile) {
    const unit = payload.units[tile.unitId];
    const color = payload.players[unit.playerId].color;

    const capsule = new THREE.Mesh(
      new THREE.CapsuleGeometry(0.25, 0.25, 4, 8),
      new THREE.MeshBasicMaterial({ color }),
    );
    capsule.position.y = 0.5;
    capsule.position.x = x;
    capsule.position.z = y;

    scene.add(capsule);
  }

  const terrain = [];
  payload.world.forEach((row, y) =>
    row.forEach((tile, x) => {
      let color = 0xd4d4d4;

      switch (tile.structure) {
        case 1:
          color = 0xf5f5f5;
          createTower(x, y, tile);
          break;
        case 2:
          color = 0xa3a3a3;
          if (tile.unitId !== null) {
            createUnit(x, y, tile);
          }
          break;
      }

      const mesh = new THREE.Mesh(
        new THREE.BoxGeometry(1, 0.2, 1),
        new THREE.MeshBasicMaterial({ color }),
      );
      mesh.position.x = x;
      mesh.position.z = y;
      terrain.push(mesh);
      tile.mesh = mesh;
    }),
  );

  scene.add(...terrain);

  function animate() {
    renderer.render(scene, camera);
  }

  if (WebGL.isWebGL2Available()) {
    container.appendChild(renderer.domElement);
    renderer.setAnimationLoop(animate);
  } else {
    const warning = WebGL.getWebGL2ErrorMessage();
    container.appendChild(warning);
  }

  const raycaster = new THREE.Raycaster();

  let selected = null;
  container.addEventListener("mousedown", (e) => {
    raycaster.setFromCamera(
      new THREE.Vector2(
        (e.clientX / window.innerWidth) * 2 - 1,
        -(e.clientY / window.innerHeight) * 2 + 1,
      ),
      camera,
    );
    const intersections = raycaster.intersectObjects(terrain);
    console.log(intersections);
    if (intersections.length > 0) {
      const intersection = intersections[0];
      const x = intersection.object.position.z;
      const y = intersection.object.position.x;
      const tile = payload.world[x][y];
      console.log([x, y], tile);
      if (tile.structure === TOWER) {
        selected = payload.towers[tile.structureId];
        console.log("selected", selected);
        const adjacentRoads = DIRECTIONS.map((dir) => [x + dir[0], y + dir[1]])
          .filter(
            (dir) =>
              dir[0] >= 0 &&
              dir[0] < payload.world.length &&
              dir[1] >= 0 &&
              dir[1] < payload.world.length,
          )
          .filter(
            (coords) => payload.world[coords[0]][coords[1]].structure === ROAD,
          )
          .map((dir) => ({
            tile: payload.world[dir[0]][dir[1]],
            road: payload.roads[tile.structureId],
            coords: dir,
          }));
        console.log("adjacentRoads", adjacentRoads);
        //if (selected.targetRoadId !== null) {
        adjacentRoads.forEach(({ tile }) => {
          let color = 0xfca5a5;
          if (tile.structureId === selected.targetRoadId) {
            color = 0xd9f99d;
          }
          console.log(tile, color);
          setColor(tile.mesh, color);
        });
        //}
      }
    }
  });
}
