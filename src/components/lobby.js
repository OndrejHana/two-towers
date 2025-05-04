import { navigate } from "../routes";
import { getAuth } from "../auth";

export async function renderLobby() {
  const clerk = getAuth();
  if (!clerk.session) {
    navigate("/login");
    return;
  }

  const c = new AbortController();
  const parent = document.getElementById("app");

  const lobbyCode = location.pathname.split("/").pop();

  parent.innerHTML = `
<div class="w-full bg-neutral-100 h-full flex md:pt-32 justify-center relative">
    <div class="p-4 bg-neutral-50 shadow-sm rounded-lg h-fit max-w-xl w-full space-y-8">
        <h1 class="text-2xl font-bold text-center">Lobby</h1>
        <div class="bg-white shadow-sm rounded p-2 flex flex-col gap-2" id="players-list">
            <p class="text-center text-neutral-500">Loading players...</p>
        </div>
        <div class="flex justify-between">
            <button type="button" class="p-2 rounded hover:bg-neutral-200  hover:cursor-pointer text-left" id="leave-lobby-button">Leave</button>
            <div class="flex gap-2" >
                <div id="code-copy" class="flex items-center justify-center gap-2 bg-neutral-200 rounded font-mono px-2 hover:bg-neutral-200/60 hover:cursor-copy overflow-hidden">
                    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-copy-icon lucide-copy"><rect width="14" height="14" x="8" y="8" rx="2" ry="2"/><path d="M4 16c-1.1 0-2-.9-2-2V4c0-1.1.9-2 2-2h10c1.1 0 2 .9 2 2"/></svg>
                    <p>
                        ${lobbyCode}
                    </p>
                </div>
                <button type="button" class="inline-flex items-center rounded bg-green-600 px-3 py-2 text-sm font-semibold text-white shadow-xs hover:bg-green-500 hover:cursor-pointer focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-green-600" id="ready-button">Ready</button>
            </div>
        </div>
    </div>
</div>
`;

  const leaveButton = document.getElementById("leave-lobby-button");
  const readyButton = document.getElementById("ready-button");
  const playersList = document.getElementById("players-list");
  const codeCopy = document.getElementById("code-copy");

  function fillPlayerList(players) {
    playersList.innerHTML = players
      .map(
        (player) => `
        <div class="flex items-center justify-between p-2 border-b border-neutral-300 last:border-b-0">
          <span class="font-medium">${player.username}</span>
          <span class="text-sm ${player.ready ? "text-green-600" : "text-neutral-500"}">
            ${player.ready ? "Ready" : "Not Ready"}
          </span>
        </div>
      `,
      )
      .join("");
  }

  const token = await clerk.session.getToken();
  fetch(`/lobby/${lobbyCode}`, {
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
    },
    signal: c.signal,
  })
    .then((res) =>
      res
        .json()
        .then((data) => {
          fillPlayerList(data.players);
        })
        .catch((err) => {
          console.log(err);
          alert("Invalid response");
        }),
    )
    .catch((err) => {
      console.log(err);
      alert("Coult not load lobby");
    });

  codeCopy.addEventListener(
    "mousedown",
    async function (_) {
      await navigator.clipboard.writeText(lobbyCode);
    },
    { signal: c.signal },
  );

  async function checkStatus() {
    while (!c.signal.aborted) {
      const token = await clerk.session.getToken();
      const res = await fetch(`/lobby/${lobbyCode}/status`, {
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        signal: c.signal,
      });
      if (res.ok) {
        const body = await res.json();
        if (body.event.EventType === 4 && !!body.event.GameId) {
          navigate(`/game/${body.event.GameId}`);
        }
        fillPlayerList(body.players);
      }
    }
  }
  checkStatus();

  leaveButton.addEventListener(
    "mousedown",
    async function (_) {
      const token = await clerk.session.getToken();
      const res = await fetch(`/lobby/${lobbyCode}/leave`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        signal: c.signal,
      });
      console.log(res);
      if (res.ok) {
        navigate("/");
      }
    },
    {
      signal: c.signal,
    },
  );

  readyButton.addEventListener(
    "mousedown",
    async function (e) {
      const token = await clerk.session.getToken();
      await fetch(`/lobby/${lobbyCode}/toggle`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        signal: c.signal,
      });
    },
    { signal: c.signal },
  );

  return c;
}
