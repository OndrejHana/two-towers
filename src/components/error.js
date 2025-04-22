/**
 * @param {HTMLDivElement} parent
 * @param {Error | null} error
 */
export function renderError(parent, error) {
  const html = `<div class="w-full h-full p-4"><div class="rounded p-2 bg-red-100 text-red-500 text-sm space-y-2">${error !== null ? `<h1 class="font-bold text-xl">${error.name}</h1><p>${error.message}</p>` : `<p>No error in state :)</p>`}</div></div>`;
  parent.innerHTML = html;
}
