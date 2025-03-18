import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/lobby/$id")({
  component: RouteComponent,
});

function RouteComponent() {
  return <div>Hello "/lobby/$id"!</div>;
}
