import { getAuth } from "../auth";
import { navigate } from "../routes";

export async function renderLogin() {
  const clerk = getAuth();
  if (clerk.session) {
    navigate("/");
    return;
  }

  const app = document.getElementById("app");
  app.innerHTML = `<div class="w-full h-full flex items-center justify-center"><div id="login-container"></div></div>`;
  const container = document.getElementById("login-container");
  clerk.mountSignIn(container);
}
