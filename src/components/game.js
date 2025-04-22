import { init } from "../game/lib";

/**
 * @param {HTMLDivElement} parent
 * @param {Payload} payload
 * @param {WebSocket} ws
 */
export function renderGame(parent, payload, ws, clerk) {
  parent.innerHTML = "";
  const conn = new WebSocket("ws://" + document.location.host + "/game/ws");
  init(payload, conn, parent);
}
