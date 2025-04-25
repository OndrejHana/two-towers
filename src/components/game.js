import { Clerk } from "@clerk/clerk-js";
import { getAuth } from "../auth";
import { init } from "../game/lib";
import { navigate } from "../routes";

async function fetchPayload(clerk, signal) {
  const res = await fetch("game/new", {
    headers: {
      Authorization: `Bearer ${await clerk.session.getToken()}`,
    },
    signal,
  });
  return await res.json();
}

/**
 * @param {Clerk} clerk
 * @param {AbortSignal} signal
 */
async function initWs(clerk, signal) {
  const url = `ws://${document.location.host}/game/ws?Authorization=${await clerk.session.getToken()}`;
  const ws = new WebSocket(url);
  signal.addEventListener("abort", () => ws.close());
  ws.onopen = () => console.log(ws.readyState);
  return ws;
}

/**
 * @param {HTMLDivElement} parent
 * @param {WebSocket} ws
 */
export async function renderGame() {
  const clerk = getAuth();
  if (!clerk) {
    navigate("/login");
    return;
  }

  const c = new AbortController();

  const payload = await fetchPayload(clerk, c.signal);
  const ws = await initWs(clerk, c.signal);

  console.log(payload, ws, ws.readyState);
  init(payload, ws, document.getElementById("app"));

  return c;
}
