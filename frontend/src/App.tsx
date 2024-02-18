import type { Component } from 'solid-js';
import { createSignal, createResource, For } from 'solid-js';

const apiURL = 'http://vwsstorage.localdomain:3000/api/v1';
const apiGetDevices = apiURL + '/device';
const apiGetDevice = apiURL + '/device/{id}';
const apiGetDeviceEvents = apiURL + '/device/{id}/event?start={start}&end={end}';

interface IDevice {
  Id: string,
  name: string,
  label: string,
  type: string,
  room: string
}

const Device = {
  Id: '',
  name: '',
  label: '',
  type: '',
  room: ''
}

const fetchDevices = async () =>
  (await fetch(`http://vwsstorage.localdomain:3000/api/v1/device`)).json();

async function fetchAPI(apiUri: string) {
  console.log('fetchAPI: ', apiUri)
  //return (await fetch(apiUri)).json();
  try {
    const response = (await fetch(apiUri)).json();
    return response
  } catch (error) {
    console.error('Error fetching data:', error);
    throw error
  }
}

// function fetchAPI(apiUri: string) {
//   console.log('fetchAPI: ', apiUri)
//   try {
//     const response = async () => (await fetch(apiUri)).json();
//     return response
//   } catch (error) {
//     console.error('Error fetching data:', error);
//     throw error
//   }
// }
// const fetchDevices = async () =>
//   (await fetch(`http://vwsstorage.localdomain:3000/api/v1/device`)).json();

const App: Component = () => {
  const [area, setArea] = createSignal('device_details')
  const [readDevice, setReadDevice] = createSignal<boolean>(false);
  const toggle = () => setReadDevice(!readDevice())
  const [devices] = createResource(readDevice, async() => await fetchAPI(apiGetDevices));

  const [readDeviceDetail, setReadDeviceDetail] = createSignal<IDevice>(Device);
  //const [deviceDetail] = createResource(readDeviceDetail, async() => await fetchAPI(apiGetDevice.replace('{id}', () => readDeviceDetail().Id)));
  const [deviceDetail] = createResource(readDeviceDetail, async() => await fetchAPI(apiGetDevice.replace('{id}', readDeviceDetail().Id)));

  //const [devices] = createResource(readDevice, fetchDevices);

  function loadDeviceDetails(device: IDevice) {
    console.log('loadDeviceDetails: ', device)
    //setReadDeviceDetail({ ...readDeviceDetail(), Id: device.Id, name: device.name, label: device.label, type: device.type, room: Math.random().toString()})
    setReadDeviceDetail(device)
    console.log('readDeviceDetail: ', readDeviceDetail())
    console.log('test replace: ', apiGetDevice.replace('{id}', readDeviceDetail().Id))

  }

  toggle()

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
            <p>Select a device:</p>
            {/* <select id="mySelect" size="4" class="px-2">
              <option value="option1">Option 1</option>
              <option value="option2">Option 2</option>
              <option value="option3">Option 3</option>
              <option value="option4">Option 4</option>
              <option value="option4">Option 5</option>
            </select> */}

            <select id="deviceList" size="40" class="px-2 w-full">
              {/* <select id="deviceList" size="40" class="px-2 w-72"> */}
              <For each={devices()}>
                {(device) => (
                  <option value={device.id} onClick={() => loadDeviceDetails(device)}>{device.label}</option>
                )}
              </For>
            </select>

            <br></br>
            <button onClick={toggle} class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 border border-blue-700 rounded">Load Devices</button>

            <br></br>
            {/* <span>{devices.loading && "Loading..."}</span>
            <pre>{JSON.stringify(devices(), null, 2)}</pre> */}
          </div>
        </aside >
        <main class="h-full w-full overflow-auto flex flex-row flex-wrap p-2">
          <nav class="flex border-b">
            <button class={(area() == 'device_details' ? 'bg-white inline-block border-l border-t border-r rounded-t py-2 px-4 text-blue-700 font-semibold' : 'bg-white inline-block py-2 px-4 text-blue-500 hover:text-blue-800 font-semibold')} onclick={() => setArea('device_details')}>
              Device Details
            </button>
            <button class={(area() == 'device_events' ? 'bg-white inline-block border-l border-t border-r rounded-t py-2 px-4 text-blue-700 font-semibold' : 'bg-white inline-block py-2 px-4 text-blue-500 hover:text-blue-800 font-semibold')} onclick={() => setArea('device_events')}>
              Device Events
            </button>
          </nav>
          {/* <div class='w-full md:w-full bg-blue-100 h-[calc(100vh-100px)]'>
            <p class={'text-lg ' + (area() == 'device_details' ? '' : 'hidden')}>Device Details</p>
            <p class={'text-lg ' + (area() == 'device_events' ? '' : 'hidden')}>Device Events</p>
            <p class={'text-lg ' + (area() == 'device_details' ? '' : 'hidden')}>Test readDeviceDetail() trigger: {readDeviceDetail().Id}</p>
            test333333333333333333333333333333
            <p class={'text-lg ' + (area() == 'device_details' ? '' : 'hidden')}>{readDeviceDetail() == Device ? '' : JSON.stringify(deviceDetail(), null, 2)}</p>
          </div> */}
          <div id='device_details' class={'w-full md:w-full bg-blue-100 h-[calc(100vh-100px)] ' + (area() == 'device_details' ? '' : 'hidden')}>
            <p class='text-lg'>Device Details</p>
            <p class='text-lg'>Test readDeviceDetail() trigger: {readDeviceDetail().Id}</p>
            test333333333333333333333333333333
            test44444
            <p class='text-lg'>{readDeviceDetail() == Device ? '' : JSON.stringify(deviceDetail(), null, 2)}</p>
          </div>
          <div id='device_events' class={'w-full md:w-full bg-blue-100 h-[calc(100vh-100px)] ' + (area() == 'device_events' ? '' : 'hidden')}>
            <p class={'text-lg '}>Device Events</p>
            device_events
          </div>
        </main>
      </div >
      <footer class="h-[40px] text-white bg-black text-2xl px-2">footer</footer>
    </>
  );
};

export default App;
