import { navigate } from "../routes";
import { getAuth } from "../auth";

export function renderLobby() {
  const clerk = getAuth();
  if (!clerk.session) {
    navigate("/login");
    return;
  }

  const c = new AbortController();
  const parent = document.getElementById("app");

  // Extract lobby code from URL
  const lobbyCode = location.pathname.split("/").pop();

  parent.innerHTML = `
<div class="w-full bg-neutral-100 h-full flex md:pt-32 justify-center relative">
    <div class="p-4 bg-neutral-50 shadow-sm rounded-lg h-fit max-w-xl w-full space-y-8">
        <h1 class="text-2xl font-bold text-center">Lobby ${lobbyCode}</h1>
        <div class="bg-white shadow-sm rounded p-2 flex flex-col gap-2" id="players-list">
            <p class="text-center text-neutral-500">Loading players...</p>
        </div>
        <div class="flex justify-between">
            <button type="button" class="p-2 rounded hover:bg-white hover:shadow-sm hover:cursor-pointer text-left" id="leave-lobby-button">Leave</button>
            <button type="button" class="inline-flex items-center rounded bg-green-600 px-3 py-2 text-sm font-semibold text-white shadow-xs hover:bg-green-500 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-green-600" id="ready-button">Ready</button>
        </div>
    </div>
</div>
`;

  const leaveButton = document.getElementById("leave-lobby-button");
  const readyButton = document.getElementById("ready-button");
  const playersList = document.getElementById("players-list");

  // Fetch and display players
  const updatePlayers = async () => {
    try {
      const token = await clerk.session.getToken();
      const response = await fetch(`/lobby/${lobbyCode}/players`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      if (!response.ok) {
        throw new Error("Failed to fetch players");
      }

      const data = await response.json();
      playersList.innerHTML = data.players.map(player => `
        <div class="flex items-center justify-between p-2 border-b last:border-b-0">
          <span class="font-medium">${player.username}</span>
          <span class="text-sm ${player.ready ? "text-green-600" : "text-neutral-500"}">
            ${player.ready ? "Ready" : "Not Ready"}
          </span>
        </div>
      `).join("");
    } catch (error) {
      console.error("Error fetching players:", error);
      playersList.innerHTML = `<p class="text-center text-red-500">Error loading players</p>`;
    }
  };

  // Initial fetch
  updatePlayers();

  // Set up polling for player updates
  const pollInterval = setInterval(updatePlayers, 2000);

  leaveButton.addEventListener("mousedown", async () => {
    try {
      const token = await clerk.session.getToken();
      await fetch(`/lobby/${lobbyCode}/leave`, {
        method: "POST",
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });
      navigate("/main");
    } catch (error) {
      console.error("Error leaving lobby:", error);
      alert("Failed to leave lobby. Please try again.");
    }
  }, { signal: c.signal });

  readyButton.addEventListener("mousedown", async () => {
    try {
      const token = await clerk.session.getToken();
      const response = await fetch(`/lobby/${lobbyCode}/ready`, {
        method: "POST",
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      if (!response.ok) {
        throw new Error("Failed to set ready status");
      }

      const data = await response.json();
      console.log(data);
      if (data.gameStarted) {
        navigate(`/game/${lobbyCode}`);
      }
    } catch (error) {
      console.error("Error setting ready status:", error);
      alert("Failed to set ready status. Please try again.");
    }
  }, { signal: c.signal });

  // Cleanup
  c.signal.addEventListener("abort", () => {
    clearInterval(pollInterval);
  });

  return c;
}
