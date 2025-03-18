import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/lobby")({
  component: Page,
});

function Page() {
  return (
    <div className="overflow-none flex h-screen w-screen items-center justify-center">
      <div className="flex flex-col justify-center gap-2 p-2">
        <div className="flex gap-2">
          <button className="rounded p-2 text-center hover:bg-neutral-100">
            Leave
          </button>
          <button className="rounded p-2 text-center hover:bg-neutral-100">
            Ready
          </button>
        </div>
      </div>
    </div>
  );
}
