import { renderRoot } from "./components/root";
import { renderLogin } from "./components/login";
import { renderMain } from "./components/main";
import { renderLobby } from "./components/lobby";
import { renderGame } from "./components/game";
import { render404 } from "./components/404";

// Define route patterns
const routePatterns = {
  "/": { pattern: /^\/$/, render: renderRoot },
  "/login": { pattern: /^\/login$/, render: renderLogin },
  "/main": { pattern: /^\/main$/, render: renderMain },
  "/lobby": { pattern: /^\/lobby\/[a-zA-Z0-9]{4}$/, render: renderLobby },
  "/game": { pattern: /^\/game$/, render: renderGame },
  "/404": { pattern: /^\/404$/, render: render404 },
};

let signal = null;

export async function onNavigate() {
  const path = location.pathname;
  let route = null;

  // Find matching route pattern
  for (const [key, value] of Object.entries(routePatterns)) {
    if (value.pattern.test(path)) {
      route = { path: key, render: value.render };
      break;
    }
  }

  // If no route matches, use 404
  if (!route) {
    route = { path: "/404", render: render404 };
  }

  leavePage();
  signal = await route.render();
}

export function navigate(path = location.pathname) {
  // Add new history entry with state
  history.pushState({ path }, '', path);
  onNavigate();
}

export function leavePage() {
  if (signal) {
    signal.abort();
    signal = null;
  }
  const parentDiv = document.getElementById("app");
  parentDiv.innerHTML = "";
}
