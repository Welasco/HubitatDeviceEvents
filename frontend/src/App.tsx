import type { Component } from 'solid-js';
import { createSignal, createResource} from 'solid-js';

const fetchDevices = async () =>
  (await fetch(`http://localhost:3000/api/v1/device`)).json();

const App: Component = () => {
  const [area, setArea] = createSignal('device_details')
  const [readDevice, setReadDevice] = createSignal(false);
  const toggle = () => setReadDevice(!readDevice())
  const [devices] = createResource(readDevice, fetchDevices);


  return (
    <>
      <header class="bg-slate-950 text-white p-3 text-2xl font-bold h-[60px]">HubitatDeviceEvents</header>
      <div class="flex flex-row h-[calc(100vh-100px)]">
        <aside class="bg-slate-200 p-2 w-1/4 overflow-auto">
          {/* <nav class="h-[40px] bg-gray-300 px-2 space-x-2 flex flex-row">
            <button class='p-1 round border selection:bottom-2 selection:border-b-2' onclick={() => setArea('area1')}>
              Area 1
            </button>
            <button class='p-1 round border' onclick={() => setArea('area2')}>
              Area 2
            </button>
          </nav> */}
          <section class="h-[40px] bg-gray-300 px-2 space-x-2 flex flex-row"></section>
          <div class="bg-slate-600 h-[calc(100vh-160px)] px-2">
            test1
            <p>Please choose one or more pets:</p>
            <select id="mySelect" size="4" class="px-2">
              <option value="option1">Option 1</option>
              <option value="option2">Option 2</option>
              <option value="option3">Option 3</option>
              <option value="option4">Option 4</option>
              <option value="option4">Option 5</option>
            </select>
            <br></br>
            <button onClick={toggle}>Get Devices</button>
            <br></br>
            <span>{devices.loading && "Loading..."}</span>
            <pre>{JSON.stringify(devices(), null, 2)}</pre>
          </div>
        </aside>
        <main class="h-full w-full overflow-auto flex flex-row flex-wrap p-2">
          <nav class="h-[40px] bg-gray-300 px-2 space-x-2 flex flex-row">
            <button class='p-1 round border selection:bottom-2 selection:border-b-2' onclick={() => setArea('device_details')}>
              Device Details
            </button>
            <button class='p-1 round border' onclick={() => setArea('device_events')}>
              Device Events
            </button>
          </nav>
          <div class='w-full md:w-full bg-blue-100 h-[calc(100vh-100px)]'>
            <p class={'text-lg ' + (area() == 'device_details' ? '' : 'hidden')}>Device Details</p>
            <p class={'text-lg ' + (area() == 'device_events' ? '' : 'hidden')}>Device Events</p>
            test333333333333333333333333333333
          </div>
        </main>
      </div>
      <footer class="h-[40px] text-white bg-black text-2xl px-2">footer</footer>
    </>
  );
};

export default App;
