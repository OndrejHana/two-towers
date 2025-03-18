import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/')({
    component: Index,
})

function Index() {
    return (
        <div className="w-screen h-screen overflow-none flex items-center justify-center">
            <div className="flex flex-col p-2 gap-2 justify-center max-w-xl w-full">
                <button className='bg-neutral-100'>New game</button>
                <input type='text' className='bg-neutral-100' />
                <button className='bg-neutral-100'>Join</button>
            </div>
        </div>
    )
}
