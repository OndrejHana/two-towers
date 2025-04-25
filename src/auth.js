import { Clerk } from "@clerk/clerk-js";

const clerkPubKey = import.meta.env.VITE_CLERK_PUBLISHABLE_KEY;
let clerk = null;

export async function initAuth() {
  const a = new Clerk(clerkPubKey);
  await a.load({});
  clerk = a;
  return a;
}

/**
 * @returns {Clerk | null}
 */
export function getAuth() {
  return clerk;
}

/**
 * @param {HTMLDivElement} parent
 * @param {Clerk} clerk
 */
export function renderLoginPage(parent, clerk) {
  const html = `<div class="w-full h-full flex items-center justify-center"><div id="login-container"></div></div>`;
  parent.innerHTML = html;
  const container = document.getElementById("login-container");
  console.log(container, clerk);
  clerk.mountSignIn(container);
}
