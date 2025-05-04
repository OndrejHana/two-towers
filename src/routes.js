import { render404 } from "./components/404";
import { renderGame } from "./components/game";
import { renderLobby } from "./components/lobby";
import { renderLogin } from "./components/login";
import { renderMain } from "./components/main";

const routes = {
  "/": { pattern: /^\/$/, render: renderMain },
  "/login": { pattern: /^\/login$/, render: renderLogin },
  "/lobby": { pattern: /^\/lobby\/[a-zA-Z0-9]{4}$/, render: renderLobby },
  "/game": {
    pattern:
      /^\/game\/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$/,
    render: renderGame,
  },
  "/404": { pattern: /^\/404$/, render: render404 },
};

/**
 * @global
 * @type {?AbortSignal}
 */
let signal = null;

export async function onNavigate() {
  const path = location.pathname;
  let route = null;

  for (const [key, value] of Object.entries(routes)) {
    if (value.pattern.test(path)) {
      route = { path: key, render: value.render };
      break;
    }
  }

  if (!route) {
    route = { path: "/404", render: render404 };
  }

  console.log("on navigate: ", path, route);

  leave();
  signal = await route.render();
}

export function navigate(path) {
  history.pushState({ path }, "", path);
  onNavigate();
}

export function leave() {
  if (signal) {
    signal.abort();
    signal = null;
  }
  const parentDiv = document.getElementById("app");
  parentDiv.innerHTML = "";
}
