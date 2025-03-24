import { createFileRoute, Link } from "@tanstack/react-router";

export const Route = createFileRoute("/")({
    component: Index,
    beforeLoad: async (_) => {
        const res = await fetch("/auth");
        const auth = await res.json();
        console.log(auth);
    },
});

function Index() {
    return (
        <div className="overflow-none flex h-screen w-screen items-center justify-center">
            <div className="flex w-full max-w-xl flex-col justify-center gap-2 p-2">
                <Link to="/lobby/new">New game</Link>
                <input type="text" className="bg-neutral-100" />
                <button>Join</button>
            </div>
        </div>
    );
}
