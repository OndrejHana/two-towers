import * as THREE from "three";
import WebGL from "three/addons/capabilities/WebGL.js";

export function render(payload) {
  const container = document.querySelector("#app");
  const scene = new THREE.Scene();

  const frustumSize = 16;
  const aspect = window.innerWidth / window.innerHeight;
  const camera = new THREE.OrthographicCamera(
    (-frustumSize * aspect) / 2, // left
    (frustumSize * aspect) / 2, // right
    frustumSize / 2, // top
    -frustumSize / 2, // bottom
    0.1, // near
    1000, // far
  );

  camera.position.set(20, 20, 20);
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

    const geometry = new THREE.CylinderGeometry(0.25, 0.25, 2, 8);
    const material = new THREE.MeshBasicMaterial({ color });
    const cylinder = new THREE.Mesh(geometry, material);
    cylinder.position.x = x;
    cylinder.position.z = y;
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
      scene.add(mesh);
    }),
  );

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
}
