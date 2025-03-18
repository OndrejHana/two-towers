import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/")({
  component: Index,
});

function Index() {
  return (
    <div className="overflow-none flex h-screen w-screen items-center justify-center">
      <div className="flex w-full max-w-xl flex-col justify-center gap-2 p-2">
        <button className="bg-neutral-100">New game</button>
        <input type="text" className="bg-neutral-100" />
        <button className="bg-neutral-100">Join</button>
      </div>
    </div>
  );
}
