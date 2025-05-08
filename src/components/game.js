import { getAuth } from "../auth";
import { addWorld, init, startRendering } from "../gameRender";

function startWs(gameId, token, signal) {
  const ws = new WebSocket(`ws://localhost:8000/game/${gameId}/ws`);
  ws.addEventListener(
    "open",
    async function (_) {
      console.log("opened");
      ws.send(JSON.stringify({ userToken: token }));
    },
    { signal },
  );

  return ws;
}

export async function renderGame() {
  const gameId = location.pathname.split("/").pop();
  const c = new AbortController();
  const ws = startWs(gameId, await getAuth().session.getToken(), c.signal);
  const parent = document.getElementById("app");
  const { scene, renderer, camera } = init(
    parent.clientWidth,
    parent.clientHeight,
  );

  ws.addEventListener(
    "message",
    function (e) {
      const message = JSON.parse(e.data);
      if (message.messageType !== 1) {
        return;
      }
      const world = message.world;
      addWorld(world, scene);
      startRendering(parent, renderer, () => renderer.render(scene, camera));
    },
    { signal: c.signal },
  );
}
