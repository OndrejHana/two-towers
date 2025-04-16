import { Clerk } from "@clerk/clerk-js";

const clerkPubKey = import.meta.env.VITE_CLERK_PUBLISHABLE_KEY;

export async function initAuth() {
  console.log("start");
  const clerk = new Clerk(clerkPubKey);
  await clerk.load({});
  console.log("sup", clerk);
  return clerk;
}
