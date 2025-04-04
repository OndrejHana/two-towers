import { render } from "./render";

let conn = null;

window.addEventListener("load", async function () {
  const res = await fetch("/game/new");
  const payload = await res.json();

  conn = new WebSocket("ws://" + document.location.host + "/game/ws");

  const button = document.createElement("button");
  button.innerText = "click me";
  document.querySelector("#app").appendChild(button);

  conn.onclose = function (evt) {
    console.log("connection closed");
  };
  conn.onmessage = function (evt) {
    var messages = evt.data.split("\n");
    console.log(messages);
  };

  button.addEventListener("mousedown", (e) => {
    console.log("clicked");
    conn.send("sup");
  });

  console.log(payload);
  render(payload, conn);
});
