import * as THREE from "three";
import WebGL from "three/addons/capabilities/WebGL.js";

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
    const tower = payload.towers[tile.towerId];
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
    }),
  );

  scene.add(...terrain);

  //scene.add(new THREE.GridHelper(10, 10));

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

  container.addEventListener("mousedown", (e) => {
    raycaster.setFromCamera(
      new THREE.Vector2(
        (e.clientX / window.innerWidth) * 2 - 1,
        -(e.clientY / window.innerHeight) * 2 + 1,
      ),
      camera,
    );
    const intersections = raycaster.intersectObjects(terrain);
    if (intersections.length > 0) {
      const tile = intersections[0];
      console.log(
        tile.object.position,
        payload.world[tile.object.position.z][tile.object.position.x],
      );
    }
  });
}
