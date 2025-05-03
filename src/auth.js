import { Clerk } from "@clerk/clerk-js";

const clerkPubKey = import.meta.env.VITE_CLERK_PUBLISHABLE_KEY;

/**
 * @global
 * @type {?Clerk}
 */
let clerk = null;

export async function initAuth() {
  clerk = new Clerk(clerkPubKey);
  await clerk.load({});
  return clerk;
}

/**
 * @returns {?Clerk}
 */
export function getAuth() {
  return clerk;
}
