import { Clerk } from "@clerk/clerk-js";

const clerkPubKey = import.meta.env.VITE_CLERK_PUBLISHABLE_KEY;

export async function initAuth() {
  console.log("start");
  const clerk = new Clerk(clerkPubKey);
  await clerk.load({});
  return clerk;
}

/**
 * @param {HTMLDivElement} parent
 * @param {Clerk} clerk
 */
export function renderLoginPage(parent, clerk) {
  console.log("rendering login");
  const html = `<div class="w-full h-full flex items-center justify-center"><div id="login-container"></div></div>`;
  parent.innerHTML = html;
  const container = document.getElementById("login-container");
  console.log(container, clerk);
  clerk.mountSignIn(container);
}
