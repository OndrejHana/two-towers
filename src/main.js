import { navigate, onNavigate } from "./routes";
import { initAuth } from "./auth";

async function main() {
  console.log("start");
  window.addEventListener("popstate", (_) => onNavigate());
  const clerk = await initAuth();
  if (clerk.session) {
    navigate(location.pathname);
  } else {
    navigate("/login");
  }
}

window.addEventListener("load", main);
