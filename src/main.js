import "./types";
import { initAuth, renderLoginPage } from "./auth";
import { State, MAIN, LOGIN } from "./state";

window.addEventListener("load", async function () {
  const app = document.querySelector("#app");
  const state = new State(app);
  const clerk = await initAuth();
  state.context.clerk = clerk;
  if (!clerk.user) {
    console.log(state.renderWith(LOGIN));
  } else {
    console.log(state.renderWith(MAIN));
  }

  //if (clerk.user) {
  //  state.state = MAIN;
  //} else {
  //  state.state = LOGIN;
  //}

  //
  //const appDiv = document.querySelector("#app");
  //const button = document.createElement("button");
  //  clerk.mountUserButton(button);
  //} else {
  //  renderLoginPage(appDiv, clerk);
  //}
  //appDiv.appendChild(button);
  //console.log("printing user", clerk.user);
  //
  //const appDiv = document.querySelector("#app");
  //const button = document.createElement("button");
  //if (clerk.user) {
  //  appDiv.innerHTML = JSON.stringify(clerk.user);
  //  clerk.mountUserButton(button);
  //} else {
  //  clerk.mountSignIn(button);
  //  button.innerText = "sign in";
  //}
  //appDiv.appendChild(button);
  //
  //const res = await fetch("/game/new", {
  //  headers: {
  //    Authorization: `Bearer ${await clerk.session.getToken()}`,
  //  },
  //});
  //const payload = await res.json();
  //const conn = new WebSocket("ws://" + document.location.host + "/game/ws");
  //
  //conn.onclose = function (_) {
  //  console.log("connection closed");
  //};
  //conn.onmessage = function (evt) {
  //  var messages = evt.data.split("\n");
  //  console.log(messages);
  //};
  //init(payload, conn, appDiv);
});
