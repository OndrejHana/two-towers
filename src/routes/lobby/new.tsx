import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/lobby/new")({
  component: RouteComponent,
});

function RouteComponent() {
  return <div>Hello "/lobby/new"!</div>;
}
