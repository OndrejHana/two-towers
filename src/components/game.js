import { getAuth } from "../auth";

export async function renderGame() {
  const gameId = location.pathname.split("/").pop();
  console.log("rendering game", gameId);

  const c = new AbortController();

  const ws = new WebSocket(`ws://localhost:8000/game/${gameId}/ws`);
  ws.addEventListener(
    "open",
    async function (_) {
      console.log("hit onopen");
      const token = await getAuth().session.getToken();
      ws.send(JSON.stringify({ userToken: token }));
      console.log("sent client token");
    },
    { signal: c.signal },
  );
}
