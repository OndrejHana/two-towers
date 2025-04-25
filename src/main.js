import { initAuth } from "./auth";
import { onNavigate } from "./routes";

async function main() {
  await initAuth();

  // Handle initial URL load
  const initialPath = window.location.pathname;
  if (initialPath !== '/') {
    // For non-root paths, we need to add the initial state to history
    history.replaceState({ path: initialPath }, '', initialPath);
  }

  // Render the initial route
  onNavigate();

  // Handle browser back/forward buttons
  window.addEventListener("popstate", (event) => {
    onNavigate();
  });
}

window.addEventListener("load", main);
