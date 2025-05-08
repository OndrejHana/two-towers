import * as THREE from "three";
import WebGL from "three/addons/capabilities/WebGL.js";
import { setBaseTileColor } from "./lib";

export function init(width, height) {
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
 * @param {THREE.Scene} scene
 * @returns {THREE.Mesh[][]}
 */
export function addWorld(world, scene, tileSize = 1) {
  const grid = world.grid;
  const geometry = new THREE.BoxGeometry(tileSize, 0.1, tileSize);
  return grid.map((row, v) =>
    row.map((tile, u) => {
      console.log(u, v, tile);
      const material = new THREE.MeshBasicMaterial();
      const mesh = new THREE.Mesh(geometry, material);
      mesh.position.set(v * tileSize, 0, u * tileSize);
      setBaseTileColor(mesh, tile);
      scene.add(mesh);
      return mesh;
    }),
  );
}

export function startRendering(container, renderer, loop) {
  if (WebGL.isWebGL2Available()) {
    container.appendChild(renderer.domElement);
    renderer.setAnimationLoop(loop);
  } else {
    container.appendChild(WebGL.getWebGL2ErrorMessage());
  }
}
