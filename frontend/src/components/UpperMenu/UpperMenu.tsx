import { Component, createSignal, createResource, For } from "solid-js";

interface DevicesProps {
    deviceSelected: () => string
    area: () => string
    setArea: (area: string) => void
}

const UpperMenu: Component<DevicesProps> = (props) => {

    return (
        <nav class="flex border-b">
        <button class={(props.area() == 'device_details' ? 'bg-white inline-block border-l border-t border-r rounded-t py-2 px-4 text-blue-700 font-semibold' : 'bg-white inline-block py-2 px-4 text-blue-500 hover:text-blue-800 font-semibold')} onclick={() => props.setArea('device_details')}>
          Device Details
        </button>
        <button class={(props.area() == 'device_events' ? 'bg-white inline-block border-l border-t border-r rounded-t py-2 px-4 text-blue-700 font-semibold' : 'bg-white inline-block py-2 px-4 text-blue-500 hover:text-blue-800 font-semibold')}
          onclick={() => props.setArea('device_events')}
          disabled={props.deviceSelected() == 'false'}
          //onclick={deviceSelected() == 'true' ? () => setArea('device_events') : () => setArea('device_details')}
          >
          Device Events
        </button>

      </nav>
    )
}

export default UpperMenu;