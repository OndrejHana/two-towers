import { render } from "./render";

console.log("hello world");
const payload = {
  world: [
    [
      { structure: 0, towerId: null, unitId: null },
      { structure: 0, towerId: null, unitId: null },
      { structure: 0, towerId: null, unitId: null },
      { structure: 0, towerId: null, unitId: null },
      { structure: 0, towerId: null, unitId: null },
      { structure: 0, towerId: null, unitId: null },
      { structure: 0, towerId: null, unitId: null },
      { structure: 0, towerId: null, unitId: null },
      { structure: 0, towerId: null, unitId: null },
      { structure: 0, towerId: null, unitId: null },
    ],
    [
      { structure: 0, towerId: null, unitId: null },
      { structure: 1, towerId: 0, unitId: null },
      { structure: 2, towerId: null, unitId: null },
      { structure: 2, towerId: null, unitId: null },
      { structure: 2, towerId: null, unitId: 0 },
      { structure: 2, towerId: null, unitId: null },
      { structure: 2, towerId: null, unitId: null },
      { structure: 2, towerId: null, unitId: null },
      { structure: 0, towerId: null, unitId: null },
      { structure: 0, towerId: null, unitId: null },
    ],
    [
      { structure: 0, towerId: null, unitId: null },
      { structure: 0, towerId: null, unitId: null },
      { structure: 0, towerId: null, unitId: null },
      { structure: 0, towerId: null, unitId: null },
      { structure: 0, towerId: null, unitId: null },
      { structure: 0, towerId: null, unitId: null },
      { structure: 0, towerId: null, unitId: null },
      { structure: 0, towerId: null, unitId: null },
      { structure: 2, towerId: null, unitId: null },
      { structure: 0, towerId: null, unitId: null },
    ],
    [
      { structure: 0, towerId: null, unitId: null },
      { structure: 0, towerId: null, unitId: null },
      { structure: 0, towerId: null, unitId: null },
      { structure: 0, towerId: null, unitId: null },
      { structure: 0, towerId: null, unitId: null },
      { structure: 0, towerId: null, unitId: null },
      { structure: 0, towerId: null, unitId: null },
      { structure: 0, towerId: null, unitId: null },
      { structure: 2, towerId: null, unitId: null },
      { structure: 0, towerId: null, unitId: null },
    ],
    [
      { structure: 0, towerId: null, unitId: null },
      { structure: 0, towerId: null, unitId: null },
      { structure: 0, towerId: null, unitId: null },
      { structure: 0, towerId: null, unitId: null },
      { structure: 0, towerId: null, unitId: null },
      { structure: 0, towerId: null, unitId: null },
      { structure: 0, towerId: null, unitId: null },
      { structure: 0, towerId: null, unitId: null },
      { structure: 2, towerId: null, unitId: null },
      { structure: 0, towerId: null, unitId: null },
    ],
    [
      { structure: 0, towerId: null, unitId: null },
      { structure: 0, towerId: null, unitId: null },
      { structure: 0, towerId: null, unitId: null },
      { structure: 0, towerId: null, unitId: null },
      { structure: 0, towerId: null, unitId: null },
      { structure: 0, towerId: null, unitId: null },
      { structure: 0, towerId: null, unitId: null },
      { structure: 0, towerId: null, unitId: null },
      { structure: 2, towerId: null, unitId: null },
      { structure: 0, towerId: null, unitId: null },
    ],
    [
      { structure: 0, towerId: null, unitId: null },
      { structure: 0, towerId: null, unitId: null },
      { structure: 0, towerId: null, unitId: null },
      { structure: 1, towerId: 2, unitId: null },
      { structure: 1, towerId: 3, unitId: null },
      { structure: 0, towerId: null, unitId: null },
      { structure: 0, towerId: null, unitId: null },
      { structure: 0, towerId: null, unitId: null },
      { structure: 2, towerId: null, unitId: null },
      { structure: 0, towerId: null, unitId: null },
    ],
    [
      { structure: 0, towerId: null, unitId: null },
      { structure: 0, towerId: null, unitId: null },
      { structure: 0, towerId: null, unitId: null },
      { structure: 0, towerId: null, unitId: null },
      { structure: 0, towerId: null, unitId: null },
      { structure: 0, towerId: null, unitId: null },
      { structure: 0, towerId: null, unitId: null },
      { structure: 0, towerId: null, unitId: null },
      { structure: 2, towerId: null, unitId: null },
      { structure: 0, towerId: null, unitId: null },
    ],
    [
      { structure: 0, towerId: null, unitId: null },
      { structure: 0, towerId: null, unitId: null },
      { structure: 0, towerId: null, unitId: null },
      { structure: 0, towerId: null, unitId: null },
      { structure: 0, towerId: null, unitId: null },
      { structure: 0, towerId: null, unitId: null },
      { structure: 0, towerId: null, unitId: null },
      { structure: 0, towerId: null, unitId: null },
      { structure: 1, towerId: 1, unitId: null },
      { structure: 0, towerId: null, unitId: null },
    ],
    [
      { structure: 0, towerId: null, unitId: null },
      { structure: 0, towerId: null, unitId: null },
      { structure: 0, towerId: null, unitId: null },
      { structure: 0, towerId: null, unitId: null },
      { structure: 0, towerId: null, unitId: null },
      { structure: 0, towerId: null, unitId: null },
      { structure: 0, towerId: null, unitId: null },
      { structure: 0, towerId: null, unitId: null },
      { structure: 0, towerId: null, unitId: null },
      { structure: 0, towerId: null, unitId: null },
    ],
  ],
  players: [{ color: "#0000ff" }, { color: "#00ff00" }, { color: "#ff0000" }],
  towers: [
    { playerId: 0 },
    { playerId: null },
    { playerId: 2 },
    { playerId: 1 },
  ],
  units: [{ playerId: 1 }],
};

window.addEventListener("load", async function () {
  render(payload);
});
