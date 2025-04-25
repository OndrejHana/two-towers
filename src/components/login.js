import { getAuth } from "../auth";

export function renderLogin() {
  const clerk = getAuth();
  if (!clerk) {
    console.log("noclerk");
    return;
  }
  const app = document.getElementById("app");
  app.innerHTML = `<div class="w-full h-full flex items-center justify-center"><div id="login-container"></div></div>`;
  console.log(clerk);
  const container = document.getElementById("login-container");
  clerk.mountSignIn(container);
}
