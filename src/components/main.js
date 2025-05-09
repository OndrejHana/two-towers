import { getAuth } from "../auth";
import { navigate } from "../routes";

export async function renderMain() {
  const clerk = getAuth();
  if (!clerk.session) {
    navigate("/login");
    return;
  }

  const c = new AbortController();
  const app = document.getElementById("app");

  app.innerHTML = `
<div class="w-full bg-neutral-100 h-full flex md:pt-32 justify-center relative">
    <div class="p-4 bg-neutral-50 shadow-sm rounded-lg h-fit max-w-64 w-full space-y-8">
        <h1 class="text-2xl font-bold text-center">Two towers</h1>
        <div class="space-y-4 text-sm font-semibold">
            <div class="flex flex-col space-y-2">
                <input type="text" placeholder="Your username" class="p-1 border border-neutral-200 outline-none block grow rounded font-mono font-thin shadow-sm" id="username-input"/>
                <button type="button" class="inline-flex items-center rounded bg-green-600 px-3 py-2 text-sm font-semibold text-white shadow-xs hover:bg-green-500 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-green-600 w-full" id="new-lobby-button">Create Lobby</button>
            </div>
            <div class="flex flex-col space-y-2">
                <input type="text" placeholder="Enter lobby code" class="p-1 border border-neutral-200 outline-none block grow rounded font-mono font-thin shadow-sm" id="lobby-code-input"/>
                <button type="button" class="inline-flex items-center rounded bg-green-600 px-3 py-2 text-sm font-semibold text-white shadow-xs hover:bg-green-500 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-green-600" id="join-lobby-button">Join Lobby</button>
            </div>
        </div>
    </div>
    <div class="absolute top-0 right-0 p-4"><div id="user-button-container"></div></div>
</div>
`;

  const ubc = document.getElementById("user-button-container");
  clerk.mountUserButton(ubc);

  const createLobbyButton = document.getElementById("new-lobby-button");
  const joinLobbyButton = document.getElementById("join-lobby-button");
  const lobbyCodeInput = document.getElementById("lobby-code-input");
  const usernameInput = document.getElementById("username-input");

  createLobbyButton.addEventListener(
    "mousedown",
    async function (_) {
      const nameInput = usernameInput.value;
      if (!nameInput) {
        alert("Enter username");
        return;
      }

      try {
        const token = await clerk.session.getToken();
        const promise = fetch("/lobby/new", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify({
            username: nameInput.trim(),
          }),
        });

        createLobbyButton.disabled = true;
        const response = await promise;
        createLobbyButton.disabled = false;

        if (!response.ok) {
          throw new Error("Failed to create lobby");
        }

        const data = await response.json();
        navigate(`/lobby/${data.code}`);
      } catch (error) {
        console.error("Error creating lobby:", error);
        alert("Failed to create lobby. Please try again.");
      }
    },
    {
      signal: c.signal,
    },
  );

  joinLobbyButton.addEventListener(
    "mousedown",
    async function (_) {
      const codeInput = lobbyCodeInput.value;
      const nameInput = usernameInput.value;

      if (!codeInput) {
        alert("Please enter a lobby code");
        return;
      }
      const code = codeInput.trim();

      if (!nameInput) {
        alert("Enter username");
        return;
      }
      const name = nameInput.trim();

      try {
        const token = await clerk.session.getToken();
        const response = await fetch("/lobby/join", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify({
            code: code,
            username: name,
          }),
        });

        if (!response.ok) {
          const error = await response.text();
          throw new Error(error);
        }

        navigate(`/lobby/${code}`);
      } catch (error) {
        console.error("Error joining lobby:", error);
        alert(error.message || "Failed to join lobby. Please try again.");
      }
    },
    {
      signal: c.signal,
    },
  );

  return c;
}
