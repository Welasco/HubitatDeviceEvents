import { Component, createSignal, createResource, For } from "solid-js";
import { apiURL, apiGetDevices, apiGetDevice, apiGetDeviceEvents } from '../../constants'
import { IDevice } from '../../types'

interface SidebarProps {
    fetchAPI: (apiGetDevices: string) => any
    loadDevice: (device: IDevice) => void
}

const Sidebar: Component<SidebarProps> = (props) => {

    const [readDevice, setReadDevice] = createSignal<number>(1);
    const [devices] = createResource(readDevice, async () => await props.fetchAPI(apiGetDevices));
    const toggleDevice = () => setReadDevice(readDevice() == 1 ? setReadDevice(2) : setReadDevice(1))

    toggleDevice()

    return (
        <aside class="bg-slate-200 p-2 w-1/4 overflow-auto">
            <section class="h-[40px] bg-gray-300 px-2 space-x-2 flex flex-row"></section>
            <div class="bg-slate-600 h-[calc(100vh-160px)] px-2">
                <p>Select a device:</p>

                <select id="deviceList" size="49" class="px-2 w-full">
                    <For each={devices()}>
                        {(device) => (
                            <option value={device.id} onClick={() => props.loadDevice(device)}>{device.label}</option>
                        )}
                    </For>
                </select>

                <br></br>
                {/* <button onClick={toggleDevice} class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 border border-blue-700 rounded">Load Devices</button>
            <button onClick={toggleDeviceEvents} class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 border border-blue-700 rounded">Load Device Events</button> */}

                <br></br>
                {/* <span>{devices.loading && "Loading..."}</span>
            <pre>{JSON.stringify(devices(), null, 2)}</pre> */}
            </div>
        </aside >
    )
}

export default Sidebar;