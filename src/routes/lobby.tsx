import { createFileRoute } from "@tanstack/react-router"

export const Route = createFileRoute('/lobby')({
    component: Page,
})

function Page() {
    return (
        <div className="w-screen h-screen overflow-none flex items-center justify-center">
            <div className="flex flex-col p-2 gap-2 justify-center">
                <div className="flex gap-2">
                    <button className="text-center p-2 rounded hover:bg-neutral-100">Leave</button>
                    <button className="text-center p-2 rounded hover:bg-neutral-100">Ready</button>
                </div>
            </div>
        </div>
    )
}
